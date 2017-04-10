package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ServerTestHandler struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
}

func (h *ServerTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Gateway test service request ...")
	json.NewEncoder(w).Encode(h)
}

func GateWayTestServer(address string, port int) (*http.Server, error) {
	listenAddress := fmt.Sprintf("%s:%d", address, port)
	log.Println("GateWay Test B/E Port - Listen address : " + listenAddress)
	myHandler := new(ServerTestHandler)
	myHandler.Code = 200
	myHandler.Message = "Test service answer"
	var err error
	var server *http.Server
	server = &http.Server{
		Addr:           listenAddress,
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func(server *http.Server) {
		err = server.ListenAndServe()
		log.Fatal(err)
	}(server)
	
	return server, err
}
