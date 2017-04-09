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
			log.Printf("All sites : %d", len(val.Sites))
			Balancer = &val
		} else {
			// No text All Sites ...
			log.Println("Filtering Root and Caching")
			filteredSites := FilterSites("", sites, FilterSingleSite)
			log.Printf("All sites : %d", len(filteredSites))
			lB := model.Balancer{
				Valid:     true,
				Enabled:   true,
				Sites:     filteredSites,
				Diff:      -1,
				QueryText: "__root__",
			}
			Balancer = &lB
			registryMap["__root__"] = lB
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
				lB := model.Balancer{
					Valid:     false,
					Enabled:   false,
					Sites:     filteredSites,
					Diff:      0,
					QueryText: "",
				}
				Balancer = &lB
				registryMap["__invlid__"] = lB
			}
		} else if len(filteredSites) == 1 {
			log.Println("One proxy found ...")
			lB := model.Balancer{
				Valid:     true,
				Enabled:   false,
				Sites:     filteredSites,
				Diff:      -1,
				QueryText: text,
			}
			Balancer = &lB
			registryMap[text] = lB
		} else {
			log.Println("More proxies found ...")
			lB := model.Balancer{
				Valid:     true,
				Enabled:   true,
				Sites:     filteredSites,
				Diff:      -1,
				QueryText: text,
			}
			Balancer = &lB
			registryMap[text] = lB
		}
	}
	log.Println("Selected Balancer : " + Balancer.QueryText)
	log.Println("Selected Balancer Sites : " + strconv.Itoa(len(Balancer.Sites)))
	log.Println(Balancer)
	log.Println("URL Path : " + nextPath)
	return Balancer, nextPath
}

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
func HostRewriteReverseProxy(sites []model.Site, config *model.Configuration, procIndex int, indexConfig model.IndexSite, useToken bool, securityToken string) *httputil.ReverseProxy {
	registryMap := make(map[string]model.Balancer)
	director := func(req *http.Request) {
		XToken := req.Header.Get("X-GATEWAY-TOKEN")
		if !useToken || securityToken == XToken {
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
					req.Header.Add("X-GATEWAY-TOKEN", indexConfig.SecurityToken)
					log.Print("Rewriting URL : ")
					log.Println(req.URL)
				} else {
					log.Println("No index service active ...")
					req.Response.Status = http.StatusText(http.StatusNotFound)
					req.Response.StatusCode = http.StatusNotFound
					req.Body.Close()
					req.Response.Close = true
				}
				log.Println("Sending Request ...")
			} else {
				site := nextBalancedSite(Balancer)
				if site.Name == "__invalid__" {
					log.Println("No proxy found ...")
					req.Response.Status = http.StatusText(http.StatusNotFound)
					req.Response.StatusCode = http.StatusNotFound
					buffer := bytes.NewBufferString("")
					buffer.WriteString(indexConfig.ServiceAddress)
					if indexConfig.Port > 0 {
						buffer.WriteString(":")
						buffer.WriteString(strconv.FormatInt(indexConfig.Port, 10))
					}
					req.URL.Host = buffer.String()
					req.URL.Scheme = indexConfig.Protocol
					req.URL.Path = "/error"
					req.Header.Add("X-GATEWAY-TOKEN", indexConfig.SecurityToken)
					log.Print("Rewriting URL : ")
					log.Println(req.URL)
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
			}
		} else {
			log.Printf("Gateway Port %d request not authorized", config.Port)
			buffer := bytes.NewBufferString("")
			buffer.WriteString(indexConfig.ServiceAddress)
			if indexConfig.Port > 0 {
				buffer.WriteString(":")
				buffer.WriteString(strconv.FormatInt(indexConfig.Port, 10))
			}
			req.URL.Host = buffer.String()
			req.URL.Scheme = indexConfig.Protocol
			req.URL.Path = "/error?code="+strconv.Itoa(http.StatusUnauthorized)+"&message="+http.StatusText(http.StatusUnauthorized)
			req.Header.Add("X-GATEWAY-TOKEN", indexConfig.SecurityToken)
			log.Print("Rewriting URL : ")
			log.Println(req.URL)
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
