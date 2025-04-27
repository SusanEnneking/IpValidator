# IpValidator
Api endpoint to validate location of incoming IP Addresses

## Local TLS Cert Gen On Mac ##
go run /Users/{you}/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.8.darwin-amd64/src/crypto/tls/generate_cert.go

## Periodic updates using cron job ##
https://dev.maxmind.com/geoip/updating-databases/
[docker documentation](https://github.com/maxmind/geoipupdate/blob/main/doc/docker.md)

## Unofficial golang client library ##
https://github.com/oschwald/geoip2-golang

## Get ip file ##
[Install Geoipupdate](https://github.com/maxmind/geoipupdate)

## Request body format ##
```
{
    "ip": "110.241.52.60",
    "countryIsoCode": [
        "US", "UA"
    ]
}
```
countryIsoCodes should be alpha-2 code from [this list](https://www.iso.org/obp/ui/#search).

## Tests ##
I submitted my version of the main.go file to ChatGpt and asked it to write tests.  It wrote the main_test.go file and suggested I use an interface for the db, which I accepted. [This commit](https://github.com/SusanEnneking/IpValidator/commit/2b82feb2c5f4b158d7354b51c372df442111e131) shows the changes Chatgpt made.

to test ```go test  -coverprofile=coverage.out```

## Run locally without Docker ##
1. Follow [instructions to download the database](https://dev.maxmind.com/geoip/updating-databases/). Create and .env file with DB_LOCATION and PORT populated and source it.  (```source my.env```).
2. Create cert.pem and key.pem files. [One way to do this](go run /Users/{you}/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.8.darwin-amd64/src/crypto/tls/generate_cert.go).

3. If you want to debug, create a launch file and debug. If not, run ```go build``` and ```go run main.go```.

4. Use Postman to post a request to https://localhost:{your port from .env file} using the request body format mentioned earlier in this README file.

## Run locally with Docker ##
1. Create an environment file like env.sample and source it.
2. Run ```docker compose up -d```

This will run two containers.  One is the canned geoipupdate container and expose the db file for the other container.  My understanding is that this geoipupdate container will automatically update the db every x hours.  I have it set to 72.

:rocket: Chatgpt helped me tremendously with the Dockerfile for the go service.