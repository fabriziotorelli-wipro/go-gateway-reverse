package service

import (
	"gateway/model"
	"log"
	"net/http"
	"sync"
	"fmt"
)


func RestServer(config model.Configuration, sites []model.Site, waitGroup *sync.WaitGroup, procIndex int, indexConfig model.IndexSite) {
	
	listenAddress := fmt.Sprintf("%s:%d", config.Address, config.Port)
	log.Println("GateWay Port - Listen address : " + listenAddress)
	reverseProxy := HostRewriteReverseProxy(sites, &config, procIndex, indexConfig, config.UseToken, config.SecurityToken)
	log.Fatal(http.ListenAndServe(listenAddress, reverseProxy))
	waitGroup.Done()

}
