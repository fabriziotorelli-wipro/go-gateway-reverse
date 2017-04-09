package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/model"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"os"
)

func FilterIndexSites(sites []model.Site) []model.Response {
	vsf := make([]model.Response, 0)
	for _, v := range sites {
		res := model.Response{
			Name:    v.Name,
			Type:    v.Type,
			SiteObj: v,
		}
		vsf = append(vsf, res)
	}
	return vsf

}

func handleGatewayRequest(w http.ResponseWriter, sites []model.Site) {
	filteredSites := FilterIndexSites(sites)
	if len(filteredSites) == 0 {
		fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", 404, "Not Found")
	} else {
		json.NewEncoder(w).Encode(filteredSites)
	}

}

type ServerRestHandler struct {
	file string
}

func (h ServerRestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlTokens := strings.Split(html.EscapeString(r.URL.Path), "/")
	if urlTokens[1] == "shutdown" {
		log.Println("System exit in process ...")
		fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", http.StatusOK, "GateWay poweroff in progress...")
		go func () {
			time.Sleep(2500 * time.Millisecond)
			log.Println("GateWay exit according to API request ...")
			os.Exit(0)
		}()
	} else {
		value, err := strconv.Atoi(urlTokens[1])
		log.Printf("Required Process : [%d]", value)
		if err != nil {
			fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", http.StatusNotFound, http.StatusText(http.StatusNotFound))
		} else {
			log.Printf("Recovery whole configs service Process [%d]", value)
			configs, error0 := model.RetrieveConfig(h.file)
			if error0 != nil {
				log.Printf("Process [%d] Error On RetrieveConfig", value)
				fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			} else {
				log.Println("Configurations: ")
				log.Println(configs)
				log.Printf("Validate Process Id Range [%d]", value)
				if len(configs) <= value {
					log.Printf("Process [%d] Not Found", value)
					fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", http.StatusNotFound, http.StatusText(http.StatusNotFound))
				} else {
					log.Printf("Recovery of data for Process [%d]", value)
					sites, error1 := model.RetrieveSites(configs[value].File)
					if error1 != nil {
						log.Printf("Process [%d] Error On RetrieveSites", value)
						fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
					} else {
						log.Printf("Recovered of data for Process [%d] : [%d]", value, len(sites))
						log.Printf("Trying Process [%d]", value)
						handleGatewayRequest(w, sites)
					}
				}
			}
		}
	}
}

func IndexServer(config model.IndexSite, fileName string, waitGroup *sync.WaitGroup) {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(config.Address)
	buffer.WriteString(":")
	buffer.WriteString(strconv.FormatInt(config.Port, 10))
	listenAddress := buffer.String()
	log.Println("GateWay Index Port - Listen address : " + listenAddress)
	myHandler := new(ServerRestHandler)
	myHandler.file = fileName
	server := &http.Server{
		Addr:           listenAddress,
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
	waitGroup.Done()

}
