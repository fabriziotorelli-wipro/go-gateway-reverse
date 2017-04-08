package process

import (
	"errors"
	"gateway/model"
	"net/http"
	"io/ioutil"
	"bytes"
	"strconv"
	"strings"
	"log"
)

func CollectResults(request model.Response, path []string, config *model.Configuration) (model.Result, error) {
	var result model.Result
	if request.SiteObj.Address == "" {
		return result, errors.New("Empty Address for : " + request.Name)
	}
	buffer := bytes.NewBufferString("")
	buffer.WriteString(config.Protocol)
	buffer.WriteString("://")
	buffer.WriteString(request.SiteObj.Address)
	if request.SiteObj.Port > 0 {
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(request.SiteObj.Port))
	}
	if request.SiteObj.Override {
		// Site Request Values override the default ones
		if len(path) > 1 && request.SiteObj.Concat && request.SiteObj.BeforeApi {
			buffer.WriteString("/")
			buffer.WriteString(strings.Join(path[1:], "/"))
		}
		buffer.WriteString(request.SiteObj.APIUri)
		if len(path) > 1 && request.SiteObj.Concat && ! request.SiteObj.BeforeApi {
			buffer.WriteString("/")
			buffer.WriteString(strings.Join(path[1:], "/"))
		}
	} else {
		if len(path) > 1 && config.Concat && config.BeforeApi {
			buffer.WriteString("/")
			buffer.WriteString(strings.Join(path[1:], "/"))
		}
		buffer.WriteString(config.APIUrl)
		if len(path) > 1 && config.Concat && ! config.BeforeApi {
			buffer.WriteString("/")
			buffer.WriteString(strings.Join(path[1:], "/"))
		}
	}
	buffer.WriteString(config.APIUrl)
	if request.SiteObj.Override {
		if len(path) > 1 && config.Concat && ! config.BeforeApi {
			buffer.WriteString("/")
			buffer.WriteString(strings.Join(path[1:], "/"))
		}
		
	} else {

 }
	log.Println("URL: " + buffer.String())
	response, error := http.Get(buffer.String())
	if error != nil {
		return result, error
	}
	defer response.Body.Close()
	
	if response.StatusCode != 200 {
		return result, errors.New("Error => HTTP Response code : " + string(response.StatusCode))
	}
	
	responseData, error2 := ioutil.ReadAll(response.Body)
	if error2 != nil {
		return result, error2
	}
	
	responseString := string(responseData)
	result.Content = responseString
	result.Process = "[" + request.Type + "] " + request.Name
	return result, nil
}