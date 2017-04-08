package service

import (
	"bytes"
	"gateway/model"
	"html"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func FilterSites(lookup string, sites []model.Site, f func(string, model.Site) bool) []model.Response {
	vsf := make([]model.Response, 0)
	for _, v := range sites {
		if f(lookup, v) {
			res := model.Response{
				Name:    v.Name,
				Type:    v.Type,
				SiteObj: v,
			}
			vsf = append(vsf, res)
		}
	}
	return vsf

}

func FilterSingleSite(text string, site model.Site) bool {
	if text == "" {
		return true
	} else {
		//return strings.Index(strings.ToLower(site.Name), strings.ToLower(text)) >= 0
		return strings.ToLower(site.Name) == strings.ToLower(text)
	}

}

func BalancerDiscovery(urlPath string, config *model.Configuration, sites []model.Site, registryMap map[string]model.Balancer) (*model.Balancer, string) {
	Balancer := new(model.Balancer)
	log.Printf("Path : %s", urlPath)
	urlTokens := strings.Split(urlPath, "/")
	text := urlTokens[1]
	log.Printf("Text : [%s]", text)
	nextPath := ""
	if len(urlTokens) > 1 {
		nextPath = strings.Join(urlTokens[2:], "/")
	}
	log.Printf("Rest Path : [%s]", nextPath)
	if len(strings.Trim(strings.Trim(text, ""), "")) == 0 {
		if val, ok := registryMap["__root__"]; ok {
			log.Println("Recovering Root Cache")
			val = registryMap["__root__"]
			Balancer = &val
		} else {
			// No text All Sites ...
			log.Println("Filtering Root and Caching")
			filteredSites := FilterSites(text, sites, FilterSingleSite)
			Balancer := model.Balancer{
				Valid:     true,
				Enabled:   true,
				Sites:     filteredSites,
				Diff:      -1,
				QueryText: "__root__",
			}
			registryMap["__root__"] = Balancer
		}
	} else if val, ok := registryMap[text]; ok {
		log.Printf("Recovering Cache for : %s", text)
		val = registryMap[text]
		Balancer = &val
	} else {
		log.Printf("Looking for : %s", text)
		filteredSites := FilterSites(text, sites, FilterSingleSite)
		log.Printf("Results : %d", len(filteredSites))
		Balancer = new(model.Balancer)
		if len(filteredSites) == 0 {
			log.Println("No proxies found ...")
			if val, ok := registryMap["__invlid__"]; ok {
				val = registryMap["__root__"]
				Balancer = &val
			} else {
				Balancer := model.Balancer{
					Valid:     false,
					Enabled:   false,
					Sites:     filteredSites,
					Diff:      0,
					QueryText: "",
				}
				registryMap["__invlid__"] = Balancer
			}
		} else if len(filteredSites) == 1 {
			log.Println("One proxy found ...")
			Balancer := model.Balancer{
				Valid:     true,
				Enabled:   false,
				Sites:     filteredSites,
				Diff:      -1,
				QueryText: text,
			}
			registryMap[text] = Balancer
		} else {
			log.Println("More proxies found ...")
			Balancer := model.Balancer{
				Valid:     true,
				Enabled:   true,
				Sites:     filteredSites,
				Diff:      -1,
				QueryText: text,
			}
			registryMap[text] = Balancer
		}
	}
	log.Println("Selected Balancer : " + Balancer.QueryText)
	log.Println(Balancer)
	log.Println("URL Path : " + nextPath)
	return Balancer, nextPath
}

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
//			singleOutElem := model.Result{
//				Process: oneResponse.Name,
//				Content: oneResponse.Type,
//			}
//			resultOutList = append(resultOutList, singleOutElem)
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

func nextBalancedSite(Balancer *model.Balancer) model.Site {
	if Balancer.Enabled && Balancer.Valid {
		// Balanced and Valid
		Balancer.Diff++
		Balancer.Diff = Balancer.Diff % len(Balancer.Sites)
		return Balancer.Sites[Balancer.Diff].SiteObj
	} else if Balancer.Valid {
		// No Balanced Valid
		return Balancer.Sites[0].SiteObj
	} else {
		// Invalid Balancer
		site := model.Site{Name: "__invalid__"}
		return site
	}
}

// MultiHostReverseProxy returns a new ReverseProxy that rewrites
// URLs to the scheme, host, and base path provided in target. If the
// target's path is "/proxyLabel/base" and the incoming request was for "/dir",
// the target request will be for /base/dir (proxyLabel causes a rewrite of the call to the
// registered proxies in LoadBalancing for multiple addressing of the label, the root list or
// a 404 Error if the required proxy doesn't exists or is empty). No content check is performed on the
// base attributes required : Address, Port[optional], Scheme, API[Optional], overwrite[Optional], etc...
func HostRewriteReverseProxy(sites []model.Site, config *model.Configuration, procIndex int, indexConfig model.IndexSite) *httputil.ReverseProxy {
	registryMap := make(map[string]model.Balancer)
	director := func(req *http.Request) {
		Balancer, newPath := BalancerDiscovery(html.EscapeString(req.URL.Path), config, sites, registryMap)
		if Balancer.QueryText == "__root__" {
			log.Println("Root requested ...")

			if indexConfig.Enabled {
				log.Println("Try an active index service ...")
				buffer := bytes.NewBufferString("")
				buffer.WriteString(indexConfig.ServiceAddress)
				if indexConfig.Port > 0 {
					buffer.WriteString(":")
					buffer.WriteString(strconv.FormatInt(indexConfig.Port, 10))
				}
				req.URL.Host = buffer.String()
				req.URL.Scheme = indexConfig.Protocol
				req.URL.Path = "/" + strconv.Itoa(procIndex)
				log.Print("Rewriting URL : ")
				log.Println(req.URL)
			} else {
				log.Println("No index service active ...")
				req.Response.Status = http.StatusText(http.StatusNotFound)
				req.Response.StatusCode = http.StatusNotFound
				req.Body.Close()
				req.Response.Body.Close()
				req.Response.Close = true
			}
			log.Println("Sending Request ...")
		} else {
			site := nextBalancedSite(Balancer)
			if site.Name == "__invalid__" {
				log.Println("No proxy found ...")
				req.Response.Status = http.StatusText(http.StatusNotFound)
				req.Response.StatusCode = http.StatusNotFound
				req.Body.Close()
				req.Response.Body.Close()
				req.Response.Close = true
			} else {
				log.Println("At least one proxy found ...")
				if len(newPath) > 0 && newPath[len(newPath)-1:] == "/" {
					newPath = newPath[0 : len(newPath)-1]
				}
				log.Printf("New Path [%s]", newPath)
				// Concatenate Api to Path without checking values [Configure Site Responsible!!] ...
				Path := config.APIUrl
				if site.Override {
					Path = site.APIUri
					if strings.Index(Path, "?") >= 0 {
						tokens := strings.Split(Path, "?")
						Path = tokens[0]
						if len(req.URL.RawQuery) > 0 {
							req.URL.RawQuery = tokens[1] + "&" + req.URL.RawQuery
						} else {
							req.URL.RawQuery = tokens[1]
						}
					}
					if site.BeforeApi && len(newPath) > 0 {
						Path = "/" + newPath + Path
					} else if len(newPath) > 0 {
						Path = Path + "/" + newPath
					}
				} else {
					if strings.Index(Path, "?") >= 0 {
						tokens := strings.Split(Path, "?")
						Path = tokens[0]
						if len(req.URL.RawQuery) > 0 {
							req.URL.RawQuery = tokens[1] + "&" + req.URL.RawQuery
						} else {
							req.URL.RawQuery = tokens[1]
						}
					}
					if config.BeforeApi && len(newPath) > 0 {
						Path = "/" + newPath + Path
					} else if len(newPath) > 0 {
						Path = Path + "/" + newPath
					}
				}
				// Upstreaming to the server site ...
				req.URL.Path = Path
				if site.Override {
					if len(site.Scheme) > 0 {
						req.URL.Scheme = site.Scheme
					}
				} else {
					if len(site.Scheme) > 0 {
						req.URL.Scheme = site.Scheme
					}
				}
				buffer := bytes.NewBufferString("")
				buffer.WriteString(site.Address)
				if site.Port > 0 {
					buffer.WriteString(":")
					buffer.WriteString(strconv.Itoa(site.Port))
				}
				host := buffer.String()
				req.URL.Host = host
				log.Print("Rewriting URL : ")
				log.Println(req.URL)
			}
			//req.URL.Scheme = target.Scheme
			//req.URL.Host = target.Host
			//if targetQuery == "" || req.URL.RawQuery == "" {
			//  req.URL.RawQuery = targetQuery + req.URL.RawQuery
			//} else {
			//  req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			//}
		}
	}
	return &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				log.Println("CALLING PROXY")
				return http.ProxyFromEnvironment(req)
			},
			Dial: func(network, addr string) (net.Conn, error) {
				log.Println("CALLING DIAL")
				conn, err := (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial(network, addr)
				if err != nil {
					log.Println("Error during DIAL:", err.Error())
				}
				return conn, err
			},
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}
