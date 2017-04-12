package service

import (
	"fmt"
	"gateway/model"
	"log"
	"net/http"
	"sync"
	"time"
	"io/ioutil"
	"crypto/x509"
	"crypto/tls"
)

func GateWayPortServer(config model.Configuration, sites []model.Site, waitGroup *sync.WaitGroup, procIndex int, indexConfig model.IndexConfig) (*http.Server, error) {

	listenAddress := fmt.Sprintf("%s:%d", config.Address, config.Port)
	log.Println("GateWay Port - Listen address : " + listenAddress)
	reverseProxy := HostRewriteReverseProxy(sites, &config, procIndex, indexConfig, config.UseToken, config.SecurityToken)
	var err error
	var server *http.Server
	server = &http.Server{
		Addr:           listenAddress,
		Handler:        reverseProxy,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func(server *http.Server, config model.Configuration, waitGroup *sync.WaitGroup) {
		if config.UseTLS {
			if config.CACertFile != "" {
				caCert, _ := ioutil.ReadFile(config.CACertFile)
				caCertPool := x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)
				server.TLSConfig = &tls.Config{
					RootCAs: caCertPool,
					ClientAuth: tls.RequireAndVerifyClientCert,
					MinVersion:               tls.VersionTLS12,
					CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
					PreferServerCipherSuites: true,
					CipherSuites: []uint16{
						tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
						tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_RSA_WITH_AES_256_CBC_SHA,
					},
				}
				server.TLSNextProto= make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
			}
			err = server.ListenAndServeTLS(config.X509CertFile, config.X509KeyFile)
		} else {
			err = server.ListenAndServe()
		}
		log.Fatal(err)
	}(server, config, waitGroup)
	return server, err

}
