package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type ValidationInfo struct {
	Ip             string   `json:"ip"`
	CountryIsoCode []string `json:"countryIsoCode"`
}

type Result struct {
	IsError        bool   `json:"isError"`
	ErrorMessage   string `json:"errorMessage"`
	IsValidCountry bool   `json:"isValidCountry"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bodyData, err := io.ReadAll(r.Body)
		if err != nil {
			slog.ErrorContext(
				context.Background(), "Failed to process request",
				slog.String("error", "Could not read request body"),
				slog.Any("details", map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "Internal server error",
				}),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var validationInfo ValidationInfo
		err = json.Unmarshal(bodyData, &validationInfo)
		if err != nil {
			slog.ErrorContext(
				context.Background(), "Failed to process request",
				slog.String("error", "could not unmarshal request"),
				slog.Any("details", map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "Internal server error",
				}),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		validatedResult := validate_country(validationInfo)
		if validatedResult.IsError {
			slog.ErrorContext(
				context.Background(), "Sending input validation error to client",
				slog.String("error", validatedResult.ErrorMessage),
			)
		}
		data, err := json.Marshal(validatedResult)
		if err != nil {
			slog.ErrorContext(
				context.Background(), "Failed to process request",
				slog.String("error", "could not mashal result"),
				slog.Any("details", map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "Internal server error",
				}),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("content-type", "application/json")
		w.Write(data)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	s := http.Server{
		Addr: port,
	}
	go func() {
		s.ListenAndServeTLS("./cert.pem", "./key.pem")
	}()
	fmt.Println("Server started, press <enter> to shut down")
	fmt.Scanln()
	s.Shutdown(context.Background())
	fmt.Println("Server shut down gracefully")
}

func validate_country(validationInfo ValidationInfo) Result {
	result := Result{IsValidCountry: false, IsError: false}
	if len(validationInfo.CountryIsoCode) == 0 || validationInfo.Ip == "" {
		result.ErrorMessage = "Ip address cannot be blank and at least one country code required."
		result.IsError = true
		return result
	}
	parsedIp := net.ParseIP(validationInfo.Ip)
	if parsedIp == nil {
		result.ErrorMessage = "The incoming ip address could not be parsed"
		result.IsError = true
		return result
	}
	dbLocation := os.Getenv("DB_LOCATION")
	dbLocation += "GeoLite2-Country.mmdb"
	db, err := geoip2.Open(dbLocation)
	if err != nil {
		result.IsError = true
		result.ErrorMessage = err.Error()
		return result
	}
	defer db.Close()
	record, err := db.Country(parsedIp)
	if err != nil {
		result.IsError = true
		result.ErrorMessage = err.Error()
		return result
	}
	for _, code := range validationInfo.CountryIsoCode {
		if code == record.Country.IsoCode {
			result.IsValidCountry = true
			break
		}
	}
	slog.InfoContext(context.Background(), "IP Lookup Complete",
		slog.String("Ip Address", validationInfo.Ip),
		slog.String("Valid Countries", strings.Join(validationInfo.CountryIsoCode, ",")),
		slog.String("Actual country code from geo look up", record.Country.IsoCode),
		slog.String("Actual country name from geo look up", record.Country.Names["en"]),
		slog.String("Valid Ip", strconv.FormatBool(result.IsValidCountry)),
	)
	return result
}
