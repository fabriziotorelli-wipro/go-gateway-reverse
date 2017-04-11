package main

import (
	"gateway/ifaces"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

//func GateWay(file string, indexFile string) {
//	log.Println("Starting Gateway ...")
//
//	indexConfig, error := model.RetrieveIndex(indexFile)
//
//	if error != nil {
//		log.Fatal(error)
//	}
//
//	configs, error0 := model.RetrieveConfig(file)
//	if error0 != nil {
//		log.Fatal(error0)
//	}
//	var waitGroup sync.WaitGroup
//	if indexConfig.Enabled {
//		waitGroup.Add(len(configs) + 1)
//	} else {
//		waitGroup.Add(len(configs))
//	}
//
//	counter := 0
//	for _, config := range configs {
//		log.Printf("Server Configuration #%d ", counter)
//		log.Printf("Address : %s", config.Address)
//		log.Printf("Port : %d", config.Port)
//		log.Printf("Protocol : %s", config.Protocol)
//		log.Printf("APIUrl : %s", config.APIUrl)
//		log.Printf("Concatenate : %t", config.Concat)
//		log.Printf("Data File Location : %s", config.File)
//		log.Printf("User : %s", config.User)
//		log.Printf("Password : %s", config.Password)
//		sites, error1 := model.RetrieveSites(config.File)
//		if error1 != nil {
//			log.Fatal(error1)
//		}
//		log.Println("List of Sites")
//		for _, site := range sites {
//			log.Printf("[%s]: [%s:%d] (type: [%s])", site.Name, site.Address, site.Port, site.Type)
//		}
//		go func(config model.Configuration, sites []model.Site, procIndex int, indexConfig model.IndexConfig) {
//			service.RestServer(config, sites, &waitGroup, procIndex, indexConfig)
//		}(config, sites, counter, indexConfig)
//		counter++
//	}
//	if indexConfig.Enabled {
//		go func(indexConfig model.IndexConfig, file string) {
//			service.IndexServer(indexConfig, file, &waitGroup)
//		}(indexConfig, file)
//	}
//
//	waitGroup.Wait()
//}

func CreateGateway(configFile string, indexFile string) ifaces.GateWay {
	gateway := ifaces.GateWayModel{
		ConfigFile:   configFile,
		IndexFile:    indexFile,
		Status:       false,
		Processes:    []ifaces.WebProcess{},
		IndexProcess: ifaces.WebProcess{},
	}
	return &gateway
}

func PrintHelp() {
	println("gateway [OPTIONS]")
	println("OPTIONS")
	println("--config-dir   [path]      Define the default configuration path (we assume as file names : config.json and indexservice.json)")
	println("--ports-file [file path] Alternatively you can define full qualified port services configuration file path")
	println("--index-file [file path] Alternatively you can define full qualified index service configuration file path")
	os.Exit(0)
}

func main() {
	Dir := string("")
	File := string("")
	IndexFile := string("")
	for index, arg := range os.Args {
		if index % 2 == 1 {
			switch arg {
			case "--help":
				PrintHelp()
			case "--config-dir":
				if index + 1 < len(os.Args) {
					Dir = os.Args[index+1]
				}
			case "--ports-file":
				if index + 1 < len(os.Args) {
					File = os.Args[index+1]
				}
			case "--index-file":
				if index + 1 < len(os.Args) {
					File = os.Args[index+1]
				}
			}
		}
	}
	if Dir != "" {
		if (File != "" || IndexFile != "") {
			println("Warning: You cannot use config dir with specific config files")
			PrintHelp()
		} else {
			if Dir[len(Dir) - 1] == "/"[0] {
				Dir = Dir[:len(Dir) - 1]
			} else if Dir[len(Dir) - 1] == "\\"[0] {
				Dir = Dir[:len(Dir) - 1]
			}
			File = Dir + "/" + "config.json"
			IndexFile = Dir  + "/" + "indexservice.json"
		}
	}
	gateway := CreateGateway(File, IndexFile)
	gateway.Start()
	gateway.Wait()

}
