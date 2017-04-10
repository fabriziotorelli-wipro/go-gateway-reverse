package ifaces

import (
	"log"
	"gateway/model"
	"sync"
	"gateway/service"
)


func StartGateWay(gateway *GateWayModel) {
	file := gateway.ConfigFile
	indexFile := gateway.IndexFile
	log.Println("Starting Gateway ...")
	log.Printf("Config file : [%s]", file)
	log.Printf("Index file : [%s]", indexFile)
	
	indexConfig, error := model.RetrieveIndex(indexFile)
	
	if error != nil {
		log.Fatal(error)
	}
	
	configs, error0 := model.RetrieveConfig(file)
	if error0 != nil {
		log.Fatal(error0)
	}
	var waitGroup sync.WaitGroup
	if indexConfig.Enabled {
		waitGroup.Add(len(configs) + 1)
	} else {
		waitGroup.Add(len(configs))
	}
	gateway.WaitGroup = &waitGroup
	
	counter := 0
	for _, config := range configs {
		log.Printf("Server Configuration #%d ", counter)
		log.Printf("Address : %s", config.Address)
		log.Printf("Port : %d", config.Port)
		log.Printf("Protocol : %s", config.Protocol)
		log.Printf("APIUrl : %s", config.APIUrl)
		log.Printf("Concatenate : %t", config.Concat)
		log.Printf("Data File Location : %s", config.File)
		log.Printf("User : %s", config.User)
		log.Printf("Password : %s", config.Password)
		sites, error1 := model.RetrieveSites(config.File)
		if error1 != nil {
			log.Fatal(error1)
		}
		log.Println("List of Sites")
		for _, site := range sites {
			log.Printf("[%s]: [%s:%d] (type: [%s])", site.Name, site.Address, site.Port, site.Type)
		}
		go func(config model.Configuration, sites []model.Site, procIndex int, indexConfig model.IndexSite, gateway *GateWayModel) {
			server, err := service.GateWayPortServer(config, sites, &waitGroup, procIndex, indexConfig)
			gateway.Processes = append(gateway.Processes, WebProcess{
				ServerError: err,
				ServerRef: server,
			})
		}(config, sites, counter, indexConfig, gateway)
		counter++
	}
	if indexConfig.Enabled {
		go func(indexConfig model.IndexSite, file string, gateway *GateWayModel) {
			server, err := service.GateWayIndexServer(indexConfig, file, &waitGroup)
			gateway.IndexProcess = WebProcess{
				ServerError: err,
				ServerRef: server,
			}
		}(indexConfig, file, gateway)
	}
	gateway.Status=true
}

func (gateway *GateWayModel) Wait() {
	gateway.WaitGroup.Wait()
	gateway.Status=false
}

func (gateway *GateWayModel) Start() {
	if ! gateway.Status {
		if len(gateway.Processes) == 0 {
			StartGateWay(gateway)
		} else {
			if gateway.IndexProcess.ServerRef != nil {
				gateway.WaitGroup.Add(len(gateway.Processes) + 1)
			} else {
				gateway.WaitGroup.Add(len(gateway.Processes))
			}
			for _, val := range gateway.Processes {
				go func() {
					val.ServerError = val.ServerRef.ListenAndServe()
					gateway.WaitGroup.Done()
				}()
			}
			if gateway.IndexProcess.ServerRef != nil {
				go func() {
					gateway.IndexProcess.ServerError = gateway.IndexProcess.ServerRef.ListenAndServe()
					gateway.WaitGroup.Done()
				}()
			} else {
				log.Println("Gateway Index not planned ...")
			}
		}
	} else {
		log.Println("Gateway already started ...")
	}
}

func (gateway *GateWayModel) Stop() {
	if gateway.Status {
		for index, val := range gateway.Processes {
			if val.ServerError == nil {
				val.ServerRef.Close()
				gateway.WaitGroup.Done()
			} else {
				log.Printf("Gateway Port [%d] not started ...", index)
			}
		}
		if gateway.IndexProcess.ServerRef != nil {
			if gateway.IndexProcess.ServerError == nil {
				gateway.IndexProcess.ServerRef.Close()
				gateway.WaitGroup.Done()
			} else {
				log.Println("Gateway Index not started ...")
			}
		} else {
			log.Println("Gateway Index not planned ...")
		}
		gateway.Status = false
	} else {
		log.Println("Gateway already stopped ...")
	}
}
