package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting hello-world server with SSL...")
	http.HandleFunc("/", helloServer)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"HelloWorld"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    ":8443",
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				{
					Certificate: [][]byte{certDER},
					PrivateKey:  priv,
				},
			},
		},
	}

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!")
}
