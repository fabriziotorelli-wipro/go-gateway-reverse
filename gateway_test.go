package main

import (
	//"io/ioutil"
	//"math/rand"
	//"os"
	"testing"

	"encoding/json"
	"gateway/ifaces"
	"gateway/model"
	"gateway/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
)

var (
	gateway = ifaces.GateWayModel{
		ConfigFile:   "./test/config.json",
		IndexFile:    "./test/indexservice.json",
		Status:       false,
		Processes:    []ifaces.WebProcess{},
		IndexProcess: ifaces.WebProcess{},
	}
	server *http.Server
	err    error
	Token  string
)

func TestInit(t *testing.T) {
	gateway.Start()
	for _, proc := range gateway.Processes {
		assert.Nil(t, proc.ServerError, "Gateway port should start without errors")
	}
	assert.Nil(t, gateway.IndexProcess.ServerError, "Gateway index should start without errors")
	server, err = test.GateWayTestServer("", 10111)
	assert.Nil(t, err, "Test Service should start")
}

func TestLoadSites(t *testing.T) {
	testSiteFile := "./test/data.json"
	sites, err := model.RetrieveSites(testSiteFile)
	assert.Nil(t, err)
	assert.NotNil(t, sites)
	assert.Equal(t, 3, len(sites))
	assert.Equal(t, "Site1", sites[0].Name)
	assert.Equal(t, "localhost", sites[0].Address)
	assert.Equal(t, 10111, sites[0].Port)
	assert.Equal(t, "http", sites[0].Protocol)
	assert.Equal(t, "http", sites[0].Scheme)
	assert.Equal(t, "json", sites[0].Type)
	assert.Equal(t, false, sites[0].Override)
	assert.Equal(t, "/", sites[0].APIUri)
	assert.Equal(t, false, sites[0].Concat)
	assert.Equal(t, false, sites[0].BeforeApi)
	assert.Equal(t, "Site2", sites[1].Name)
	assert.Equal(t, "mysite2.org", sites[1].Address)
	assert.Equal(t, 8081, sites[1].Port)
	assert.Equal(t, "http", sites[1].Protocol)
	assert.Equal(t, "http", sites[1].Scheme)
	assert.Equal(t, "json", sites[1].Type)
	assert.Equal(t, false, sites[1].Override)
	assert.Equal(t, "/v2", sites[1].APIUri)
	assert.Equal(t, false, sites[1].Concat)
	assert.Equal(t, false, sites[1].BeforeApi)
	assert.Equal(t, "Site3", sites[2].Name)
	assert.Equal(t, "mysite3.org", sites[2].Address)
	assert.Equal(t, 8082, sites[2].Port)
	assert.Equal(t, "http", sites[2].Protocol)
	assert.Equal(t, "http", sites[2].Scheme)
	assert.Equal(t, "json", sites[2].Type)
	assert.Equal(t, true, sites[2].Override)
	assert.Equal(t, "/v2", sites[2].APIUri)
	assert.Equal(t, true, sites[2].Concat)
	assert.Equal(t, true, sites[2].BeforeApi)
}

func TestIndexServiceConfig(t *testing.T) {
	testIndexFile := "./test/indexservice.json"

	index, err := model.RetrieveIndex(testIndexFile)
	assert.Nil(t, err)
	assert.NotNil(t, index)
	assert.Equal(t, true, index.Enabled)
	assert.Equal(t, "", index.Address)
	assert.Equal(t, int64(10098), index.Port)
	assert.Equal(t, "http", index.Protocol)
	assert.Equal(t, "localhost", index.ServiceAddress)
	assert.Equal(t, "J1qK1c18UUGJFAzz9xnH56584l4", index.SecurityToken)
}

func TestLoadPorts(t *testing.T) {
	testPortFile := "./test/config.json"

	ports, err := model.RetrieveConfig(testPortFile)
	assert.Nil(t, err)
	assert.NotNil(t, ports)
	assert.Equal(t, 3, len(ports))
	assert.Equal(t, "", ports[0].Address)
	assert.Equal(t, int64(10099), ports[0].Port)
	assert.Equal(t, "http", ports[0].Protocol)
	assert.Equal(t, "./test/data2.json", ports[0].File)
	assert.Equal(t, "", ports[0].User)
	assert.Equal(t, "", ports[0].Password)
	assert.Equal(t, false, ports[0].UseToken)
	assert.Equal(t, "", ports[0].SecurityToken)
	assert.Equal(t, "/", ports[0].APIUrl)
	assert.Equal(t, true, ports[0].Concat)
	assert.Equal(t, true, ports[0].BeforeApi)
	assert.Equal(t, "", ports[1].Address)
	assert.Equal(t, int64(10100), ports[1].Port)
	assert.Equal(t, "http", ports[1].Protocol)
	assert.Equal(t, "./test/data.json", ports[1].File)
	assert.Equal(t, "", ports[1].User)
	assert.Equal(t, "", ports[1].Password)
	assert.Equal(t, true, ports[1].UseToken)
	assert.Equal(t, "J1qK1c18UUGJFAzz9xnH56584l4", ports[1].SecurityToken)
	assert.Equal(t, "/", ports[1].APIUrl)
	assert.Equal(t, true, ports[1].Concat)
	assert.Equal(t, true, ports[1].BeforeApi)
	Token = ports[1].SecurityToken
	assert.Equal(t, "", ports[2].Address)
	assert.Equal(t, int64(10101), ports[2].Port)
	assert.Equal(t, "http", ports[2].Protocol)
	assert.Equal(t, "./test/data3.json", ports[2].File)
	assert.Equal(t, "", ports[2].User)
	assert.Equal(t, "", ports[2].Password)
	assert.Equal(t, false, ports[2].UseToken)
	assert.Equal(t, "", ports[2].SecurityToken)
	assert.Equal(t, "/", ports[2].APIUrl)
	assert.Equal(t, true, ports[2].Concat)
	assert.Equal(t, true, ports[2].BeforeApi)
}

func TestGatewayServicePortPlain(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:10101/Site1", nil)
	resp, error1 := client.Do(req)
	assert.Nil(t, error1, "Test Service should be reachable")
	var respData test.ServerTestHandler
	err = json.NewDecoder(resp.Body).Decode(&respData)
	resp.Body.Close()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 200, respData.Code)
	assert.Equal(t, "Test service answer", respData.Message)
}

func TestGatewayServicePortToken(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:10100/Site1", nil)
	req.Header.Add("X-GATEWAY-TOKEN", Token)
	resp, error1 := client.Do(req)
	assert.Nil(t, error1, "Test Service should be reachable")
	var respData test.ServerTestHandler
	err = json.NewDecoder(resp.Body).Decode(&respData)
	resp.Body.Close()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 200, respData.Code)
	assert.Equal(t, "Test service answer", respData.Message)
}

func TestGatewayServiceIndex(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:10100/", nil)
	req.Header.Add("X-GATEWAY-TOKEN", Token)
	resp, error1 := client.Do(req)
	assert.Nil(t, error1, "Test Service should be reachable")
	var respData []model.Response
	err := json.NewDecoder(resp.Body).Decode(&respData)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 3, len(respData))
	assert.Equal(t, "Site1", respData[0].Name)
	assert.Equal(t, "json", respData[0].Type)
	assert.Equal(t, "Site2", respData[1].Name)
	assert.Equal(t, "json", respData[1].Type)
	assert.Equal(t, "Site3", respData[2].Name)
	assert.Equal(t, "json", respData[2].Type)
}

func TestGatewayServicePortCertificate(t *testing.T) {
	cert, err := tls.LoadX509KeyPair("./test/server.pem", "./test/server.key")
	caCert, err := ioutil.ReadFile("./test/server.pem")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	
	resp, error1 := client.Get("https://localhost:10099/Site1")
	assert.Nil(t, error1, "Test Service should be reachable")
	var respData test.ServerTestHandler
	err = json.NewDecoder(resp.Body).Decode(&respData)
	resp.Body.Close()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 200, respData.Code)
	assert.Equal(t, "Test service answer", respData.Message)
}

func TestEnd(t *testing.T) {
	go func () {
		server.Close()
		gateway.Stop()
	}()
}

//func getFileAsString(path string) string {
//	buf, err := ioutil.ReadFile("_tests/" + path)
//	if err != nil {
//		panic(err)
//	}
//
//	return string(buf)
//}

//func getRandomString(n int) string {
//	rand.Seed(time.Now().UnixNano())
//	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
//
//	b := make([]rune, n)
//	for i := range b {
//		b[i] = letters[rand.Intn(len(letters))]
//	}
//	return string(b)
//}

//func getFileReader(file string) *os.File {
//	r, err := os.Open(file)
//	if err != nil {
//		panic(err)
//	}
//	return r
//}
