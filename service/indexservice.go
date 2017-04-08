package service

import (
	"log"
	"net/http"
	"gateway/model"
	"bytes"
	"strconv"
	"sync"
	"html"
	"strings"
	"encoding/json"
	"fmt"
	"time"
)


func FilterIndexSites(sites []model.Site) ([]model.Response) {
	vsf := make([]model.Response, 0)
	for _, v := range sites {
		res := model.Response{
			Name: v.Name,
			Type: v.Type,
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
	} else  {
		json.NewEncoder(w).Encode(filteredSites)
	}

}
type ServerRestHandler struct{
	file string
}

func (h ServerRestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlTokens := strings.Split(html.EscapeString(r.URL.Path), "/")
	value, err := strconv.Atoi(urlTokens[1])
	log.Printf("Required Process : [%d]", value)
	if err != nil {
		log.Println("Process Nil Not Found")
		fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	log.Printf("Recovery whole configs service Process [%d]", value)
	configs, error0 := model.RetrieveConfig(h.file)
	if error0 != nil {
		log.Printf("Process [%d] Error On RetrieveConfig", value)
		fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	log.Println("Configurations: ")
	log.Println(configs)
	log.Printf("Validate Process Id Range [%d]", value)
	if len(configs) <= value {
		log.Printf("Process [%d] Not Found", value)
		fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	log.Printf("Recovery of data for Process [%d]", value)
	sites, error1 := model.RetrieveSites(configs[value].File)
	if error1 != nil {
		log.Printf("Process [%d] Error On RetrieveSites", value)
		fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	log.Printf("Recovered of data for Process [%d] : [%d]", value, len(sites))
	log.Printf("Trying Process [%d]", value)
	handleGatewayRequest(w, sites)
}

func IndexServer(config model.IndexSite, fileName string, waitGroup *sync.WaitGroup) {
	buffer := bytes.NewBufferString("");
	buffer.WriteString(config.Address)
	buffer.WriteString(":")
	buffer.WriteString(strconv.FormatInt(config.Port, 10))
	listenAddress := buffer.String()
	log.Println("GateWay Index Port - Listen address : " + listenAddress)
	myHandler := new (ServerRestHandler)
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
