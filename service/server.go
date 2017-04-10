package service

import (
	"fmt"
	"gateway/model"
	"log"
	"net/http"
	"sync"
	"time"
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
			err = server.ListenAndServeTLS(config.X509CertFile, config.X509KeyFile)
		} else {
			err = server.ListenAndServe()
		}
		log.Fatal(err)
	}(server, config, waitGroup)
	return server, err

}
