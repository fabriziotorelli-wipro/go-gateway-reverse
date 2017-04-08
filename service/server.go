package service

import (
	"bytes"
	"gateway/model"
	"log"
	"net/http"
	"strconv"
	"sync"
)

//func FilterSites(lookup string, sites []model.Site, f func(string, model.Site) bool) ([]model.Response) {
//	vsf := make([]model.Response, 0)
//	for _, v := range sites {
//		if f(lookup, v) {
//			res := model.Response{
//				Name: v.Name,
//				Type: v.Type,
//				SiteObj: v,
//			}
//				vsf = append(vsf, res)
//		}
//	}
//	return vsf
//
//}
//
//func FilterSingleSite(text string, site model.Site) (bool) {
//	if text=="" {
//		return true
//	} else {
//		return strings.Index(strings.ToLower(site.Name), strings.ToLower(text)) >= 0
//	}
//
//}
//
//func handleGatewayRequest(w http.ResponseWriter, urlPath string, config *model.Configuration, sites []model.Site) {
//	log.Printf("Path : %s", urlPath)
//	urlTokens := strings.Split(urlPath, "/")
//	text := urlTokens[1]
//	log.Printf("Looking for : %s", text)
//	filteredSites := FilterSites(text, sites, FilterSingleSite)
//	if len(filteredSites) == 0 {
//		fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", 404, "Not Found")
//	} else if len(filteredSites) == 1 && (len(strings.Trim(text, ""))>0) {
//		singleResult, error1 := process.CollectResults(filteredSites[0], urlTokens[1:], config)
//		if error1 == nil {
//			fmt.Fprintf(w, singleResult.Content)
//		} else {
//			fmt.Fprintf(w, "{\"code\": %d, \"message\":\"%s\"}", 404, "Unreachable Service : " + filteredSites[0].Name)
//		}
//	} else if (len(strings.Trim(text, ""))==0) {
//		resultOutList := make([]model.Result, 0)
//		for _, oneResponse := range filteredSites {
//				singleOutElem := model.Result{
//					Process: oneResponse.Name,
//					Content: oneResponse.Type,
//				}
//				resultOutList = append(resultOutList, singleOutElem)
//		}
//		json.NewEncoder(w).Encode(resultOutList)
//	} else {
//		resultOutList := make([]model.JSonResult, 0)
//		for _, oneResponse := range filteredSites {
//			singleResult, error2 := process.CollectResults(oneResponse, urlTokens[1:], config)
//
//			if error2 == nil {
//				jenkinsMap := make(map[string]interface{})
//				err := json.Unmarshal([]byte(singleResult.Content), &jenkinsMap)
//				if err == nil {
//					singleOutElem := model.JSonResult{
//						Name: oneResponse.Name,
//						Body: jenkinsMap,
//					}
//					resultOutList = append(resultOutList, singleOutElem)
//				} else {
//					errorResp := "Wrong SERVICE Json : "+html.UnescapeString(strings.Replace(singleResult.Content, "\n", "", len(singleResult.Content)))
//					errorMap := make(map[string]interface{})
//					errorMap["code"] = 403
//					errorMap["message"] = html.UnescapeString(strings.Replace(errorResp, "\n", "", 0))
//					singleOutElem := model.JSonResult{
//						Name: oneResponse.Name,
//						Body: errorMap,
//					}
//					resultOutList = append(resultOutList, singleOutElem)
//				}
//			} else {
//				errorResp := "{\"404\": %d, \"message\":\"Unreachable Service : "+html.UnescapeString(strings.Replace(oneResponse.Name, "\n", "", 0))+"\"}"
//				errorMap := make(map[string]interface{})
//				errorMap["code"] = 403
//				errorMap["message"] = html.UnescapeString(strings.Replace(errorResp, "\n", "", 0))
//				singleOutElem := model.JSonResult{
//					Name: oneResponse.Name,
//					Body: errorMap,
//				}
//				resultOutList = append(resultOutList, singleOutElem)
//			}
//		}
//		json.NewEncoder(w).Encode(resultOutList)
//
//	}
//
//}
//type ServerRestHandler struct{
//	config *model.Configuration
//	sites []model.Site
//}
//
//func (h ServerRestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	handleGatewayRequest(w, html.EscapeString(r.URL.Path), h.config, h.sites)
//}

func RestServer(config model.Configuration, sites []model.Site, waitGroup *sync.WaitGroup, procIndex int, indexConfig model.IndexSite) {
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	handleGatewayRequest(w, html.EscapeString(r.URL.Path), config, sites)
	//})

	buffer := bytes.NewBufferString("")
	buffer.WriteString(config.Address)
	buffer.WriteString(":")
	buffer.WriteString(strconv.FormatInt(config.Port, 10))
	listenAddress := buffer.String()
	log.Println("GateWay Port - Listen address : " + listenAddress)
	//myHandler := new (ServerRestHandler)
	//myHandler.config = &config
	//myHandler.sites = sites
	//server := &http.Server{
	//	Addr:           listenAddress,
	//	Handler:        nil,
	//	ReadTimeout:    10 * time.Second,
	//	WriteTimeout:   10 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
	reverseProxy := HostRewriteReverseProxy(sites, &config, procIndex, indexConfig)
	//	log.Fatal(server.ListenAndServe(), reverseProxy)
	log.Fatal(http.ListenAndServe(listenAddress, reverseProxy))
	waitGroup.Done()

}
