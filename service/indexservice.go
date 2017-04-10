package service

import (
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
		fmt.Fprintf(w, "{ \"code\": %d, \"message\":\"%s\" }", 404, "Not Found")
	} else {
		json.NewEncoder(w).Encode(filteredSites)
	}
	
}

type ServerRestHandler struct {
	file  string
	token string
}

func (h ServerRestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	XToken := r.Header.Get("X-GATEWAY-TOKEN")
	if h.token == XToken {
		urlTokens := strings.Split(html.EscapeString(r.URL.Path), "/")
		if urlTokens[1] == "shutdown" {
			log.Println("Gateway shutdown in process ...")
			fmt.Fprintf(w, "{ \"code\": %d, \"message\":\"%s\" }", http.StatusOK, "GateWay poweroff in progress...")
			go func() {
				time.Sleep(2500 * time.Millisecond)
				log.Println("GateWay exit according to API request ...")
				os.Exit(0)
			}()
		} else if urlTokens[1] == "error" {
			log.Println("Gateway error shuffling in process ...")
			w.Header().Add("code", r.URL.Query().Get("code"))
			w.Header().Add("message", r.URL.Query().Get("message"))
			fmt.Fprintf(w, "{\"code\": %s, \"message\":\"%s\"}", r.URL.Query().Get("code"), r.URL.Query().Get("message"))
		} else {
			value, err := strconv.Atoi(urlTokens[1])
			log.Printf("Required Process : [%d]", value)
			if err != nil {
				fmt.Fprintf(w, "{ \"code\": %d, \"message\":\"%s\" }", http.StatusNotFound, http.StatusText(http.StatusNotFound))
			} else {
				log.Printf("Recovery whole configs service Process [%d]", value)
				configs, error0 := model.RetrieveConfig(h.file)
				if error0 != nil {
					log.Printf("Process [%d] Error On RetrieveConfig", value)
					fmt.Fprintf(w, "{ \"code\": %d, \"message\":\"%s\" }", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				} else {
					log.Println("Configurations: ")
					log.Println(configs)
					log.Printf("Validate Process Id Range [%d]", value)
					if len(configs) <= value {
						log.Printf("Process [%d] Not Found", value)
						fmt.Fprintf(w, "{ \"code\": %d, \"message\":\"%s\" }", http.StatusNotFound, http.StatusText(http.StatusNotFound))
					} else {
						log.Printf("Recovery of data for Process [%d]", value)
						sites, error1 := model.RetrieveSites(configs[value].File)
						if error1 != nil {
							log.Printf("Process [%d] Error On RetrieveSites", value)
							fmt.Fprintf(w, "{ \"code\": %d, \"message\":\"%s\" }", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
						} else {
							log.Printf("Recovered of data for Process [%d] : [%d]", value, len(sites))
							log.Printf("Trying Process [%d]", value)
							handleGatewayRequest(w, sites)
						}
					}
				}
			}
		}
	} else {
		log.Println("Process request not authorized")
		fmt.Fprintf(w, "{ \"code\": %d, \"message\":\"%s\" }", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}
}

func GateWayIndexServer(config model.IndexSite, fileName string, waitGroup *sync.WaitGroup) (*http.Server, error) {
	listenAddress := fmt.Sprintf("%s:%d", config.Address, config.Port)
	log.Println("GateWay Index Port - Listen address : " + listenAddress)
	myHandler := new(ServerRestHandler)
	myHandler.file = fileName
	myHandler.token = config.SecurityToken
	var err error
	var server *http.Server
	server = &http.Server{
		Addr:           listenAddress,
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//err := server.ListenAndServe()
	//log.Fatal(err)
	//waitGroup.Done()
	go func(server *http.Server, waitGroup *sync.WaitGroup) {
		err = server.ListenAndServe()
		log.Fatal(err)
		//waitGroup.Done()
	}(server, waitGroup)
	
	return server, err
}
