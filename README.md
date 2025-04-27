# IpValidator
Api endpoint to validate location of incoming IP Addresses

## local TLS Cert Gen On Mac ##
go run /Users/{you}/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.8.darwin-amd64/src/crypto/tls/generate_cert.go

## periodic updates using cron job ##
https://dev.maxmind.com/geoip/updating-databases/
[docker documentation](https://github.com/maxmind/geoipupdate/blob/main/doc/docker.md)

## unofficial golang client library ##
https://github.com/oschwald/geoip2-golang

## get ip file ##
[Install Geoipupdate](https://github.com/maxmind/geoipupdate)

## request body format ##
```
{
    "ip": "110.241.52.60",
    "countryIsoCode": [
        "US", "UA"
    ]
}
```
countryIsoCodes should be alpha-2 code from [this list](https://www.iso.org/obp/ui/#search).

## tests ##
I submitted my version of the main.go file to ChatGpt and asked it to write tests.  It wrote the main_test.go file and suggested I use an interface for the db, which I accepted. [This commit](https://github.com/SusanEnneking/IpValidator/commit/2b82feb2c5f4b158d7354b51c372df442111e131) shows the changes Chatgpt made.

to test ```go test  -coverprofile=coverage.out```