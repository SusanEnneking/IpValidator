# IpValidator
Api endpoint to validate location of incoming IP Addresses

## local TLS Cert Gen On Mac ##
go run /Users/{you}/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.8.darwin-amd64/src/crypto/tls/generate_cert.go

## periodic updates using cron job ##
https://dev.maxmind.com/geoip/updating-databases/

## unofficial golang client library ##
https://github.com/oschwald/geoip2-golang

## get ip file ##
[Install Geoipupdate](https://github.com/maxmind/geoipupdate)

## request body format ##
```
{
    ip: "valid ip string",
    validCountries: [
        {coutryIsoCode: "string"},
        {countryIsoCode: "string"}
    ]
}
```
coutryIsoCodes should be alpha-2 code from [this list](https://www.iso.org/obp/ui/#search).