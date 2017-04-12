package service

import (
	"encoding/json"
	"fmt"
	"gateway/model"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"io/ioutil"
	"crypto/x509"
	"crypto/tls"
)

func FilterIndexConfigs(sites []model.Site) []model.Response {
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
	filteredSites := FilterIndexConfigs(sites)
	if len(filteredSites) == 0 {
		message := model.HttpResponse{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		}
		json.NewEncoder(w).Encode(message)
	} else {
		json.NewEncoder(w).Encode(filteredSites)
	}

}

type ServerRestHandler struct {
	file     string
	token    string
	useToken bool
}

func (h ServerRestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	XToken := r.Header.Get("X-GATEWAY-TOKEN")
	if !h.useToken || h.token == XToken {
		urlTokens := strings.Split(html.EscapeString(r.URL.Path), "/")
		if urlTokens[1] == "shutdown" {
			log.Println("Gateway shutdown in process ...")
			message := model.HttpResponse{
				Code:    http.StatusOK,
				Message: "GateWay poweroff in progress...",
			}
			json.NewEncoder(w).Encode(message)
			go func() {
				time.Sleep(2500 * time.Millisecond)
				log.Println("GateWay exit according to API request ...")
				os.Exit(0)
			}()
		} else if urlTokens[1] == "error" {
			log.Println("Gateway error shuffling in process ...")
			w.Header().Add("code", r.URL.Query().Get("code"))
			w.Header().Add("message", r.URL.Query().Get("message"))
			val, err := strconv.Atoi(r.URL.Query().Get("code"))
			if err == nil {
				message := model.HttpResponse{
					Code:    http.StatusNotFound,
					Message: r.URL.Query().Get("message"),
				}
				json.NewEncoder(w).Encode(message)
			} else {
				message := model.HttpResponse{
					Code:    val,
					Message: r.URL.Query().Get("message"),
				}
				json.NewEncoder(w).Encode(message)
			}
		} else {
			value, err := strconv.Atoi(urlTokens[1])
			log.Printf("Required Process : [%d]", value)
			if err != nil {
				message := model.HttpResponse{
					Code:    http.StatusNotFound,
					Message: http.StatusText(http.StatusNotFound),
				}
				json.NewEncoder(w).Encode(message)
			} else {
				log.Printf("Recovery whole configs service Process [%d]", value)
				configs, error0 := model.RetrieveConfig(h.file)
				if error0 != nil {
					log.Printf("Process [%d] Error On RetrieveConfig", value)
					message := model.HttpResponse{
						Code:    http.StatusInternalServerError,
						Message: http.StatusText(http.StatusInternalServerError),
					}
					json.NewEncoder(w).Encode(message)
				} else {
					log.Println("Configurations: ")
					log.Println(configs)
					log.Printf("Validate Process Id Range [%d]", value)
					if len(configs) <= value {
						log.Printf("Process [%d] Not Found", value)
						message := model.HttpResponse{
							Code:    http.StatusNotFound,
							Message: http.StatusText(http.StatusNotFound),
						}
						json.NewEncoder(w).Encode(message)
					} else {
						log.Printf("Recovery of data for Process [%d]", value)
						sites, error1 := model.RetrieveSites(configs[value].File)
						if error1 != nil {
							log.Printf("Process [%d] Error On RetrieveSites", value)
							message := model.HttpResponse{
								Code:    http.StatusInternalServerError,
								Message: http.StatusText(http.StatusInternalServerError),
							}
							json.NewEncoder(w).Encode(message)
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
		message := model.HttpResponse{
			Code:    http.StatusUnauthorized,
			Message: http.StatusText(http.StatusUnauthorized),
		}
		json.NewEncoder(w).Encode(message)
	}
}

func GateWayIndexServer(config model.IndexConfig, fileName string, waitGroup *sync.WaitGroup) (*http.Server, error) {
	listenAddress := fmt.Sprintf("%s:%d", config.Address, config.Port)
	log.Println("GateWay Index Port - Listen address : " + listenAddress)
	myHandler := new(ServerRestHandler)
	myHandler.file = fileName
	myHandler.token = config.SecurityToken
	myHandler.useToken = config.UseToken
	var err error
	var server *http.Server
	server = &http.Server{
		Addr:           listenAddress,
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func(server *http.Server, config model.IndexConfig, waitGroup *sync.WaitGroup) {
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
