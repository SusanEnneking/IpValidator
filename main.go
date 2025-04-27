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

// GeoIPLookup defines an interface for looking up country codes by IP
type GeoIPLookup interface {
	LookupCountry(ip net.IP) (string, error)
}

// RealGeoIPDB implements GeoIPLookup using a real GeoIP2 database
type RealGeoIPDB struct {
	db *geoip2.Reader
}

func (r *RealGeoIPDB) LookupCountry(ip net.IP) (string, error) {
	record, err := r.db.Country(ip)
	if err != nil {
		return "", err
	}
	return record.Country.IsoCode, nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbLocation := os.Getenv("DB_LOCATION")
	if dbLocation == "" {
		dbLocation = "./" // Default to current directory if not set
	}
	dbPath := dbLocation + "GeoLite2-Country.mmdb"
	db, err := geoip2.Open(dbPath)
	if err != nil {
		slog.Error("failed to open geoip database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	geoip := &RealGeoIPDB{db: db}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bodyData, err := io.ReadAll(r.Body)
		if err != nil {
			slog.ErrorContext(context.Background(), "Failed to read request body", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var validationInfo ValidationInfo
		err = json.Unmarshal(bodyData, &validationInfo)
		if err != nil {
			slog.ErrorContext(context.Background(), "Failed to unmarshal request body", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		validatedResult := validate_country(validationInfo, geoip)

		if validatedResult.IsError {
			slog.ErrorContext(context.Background(), "Validation error", slog.String("error", validatedResult.ErrorMessage))
		}

		data, err := json.Marshal(validatedResult)
		if err != nil {
			slog.ErrorContext(context.Background(), "Failed to marshal response", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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

func validate_country(validationInfo ValidationInfo, geoip GeoIPLookup) Result {
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

	isoCode, err := geoip.LookupCountry(parsedIp)
	if err != nil {
		result.IsError = true
		result.ErrorMessage = err.Error()
		return result
	}

	for _, code := range validationInfo.CountryIsoCode {
		if code == isoCode {
			result.IsValidCountry = true
			break
		}
	}

	slog.InfoContext(context.Background(), "IP Lookup Complete",
		slog.String("Ip Address", validationInfo.Ip),
		slog.String("Valid Countries", strings.Join(validationInfo.CountryIsoCode, ",")),
		slog.String("Actual country code from geo lookup", isoCode),
		slog.String("Valid Ip", strconv.FormatBool(result.IsValidCountry)),
	)

	return result
}
