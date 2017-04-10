package service

import (
	"gateway/model"
	"log"
	"net/http"
	"sync"
	"fmt"
	"time"
)


func GateWayPortServer(config model.Configuration, sites []model.Site, waitGroup *sync.WaitGroup, procIndex int, indexConfig model.IndexSite)  (*http.Server, error) {
	
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
	
	//log.Fatal(http.ListenAndServe(listenAddress, reverseProxy))
	//waitGroup.Done()
	go func(server *http.Server, waitGroup *sync.WaitGroup) {
		err = server.ListenAndServe()
		log.Fatal(err)
		//waitGroup.Done()
	}(server, waitGroup)
	return server, err

}
