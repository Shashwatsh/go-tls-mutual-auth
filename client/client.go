package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// process parameters
	totalRequest := flag.Int("request", 1, "How many request perform.")
	flag.Parse()

	// load client cert
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatal(err)
	}

	// load CA cert
	caCert, err := ioutil.ReadFile("server.crt")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()

	for i := 0; i < *totalRequest; i++ {
		client := &http.Client{
			Transport: &http2.Transport{TLSClientConfig: tlsConfig},
		}

		resp, err := client.Get("https://localhost:8080/hello")
		if err != nil {
			fmt.Println(err)
		}
		contents, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s\n", string(contents))
	}
}
