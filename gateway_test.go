package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gateway/model"
)

//var (
//)

func TestLoadSites(t *testing.T) {
	testSiteFile := "./test/data.json"
	//my_data := getFileAsString("data.xml")

	sites, err := model.RetrieveSites(testSiteFile)
	assert.Nil(t, err)
	assert.NotNil(t, sites)
	assert.Equal(t, 3, len(sites))
	assert.Equal(t, "Site1", sites[0].Name)
	assert.Equal(t, "mysite1.org", sites[0].Address)
	assert.Equal(t, 8080, sites[0].Port)
	assert.Equal(t, "http", sites[0].Protocol)
	assert.Equal(t, "http", sites[0].Scheme)
	assert.Equal(t, "json", sites[0].Type)
	assert.Equal(t, false, sites[0].Override)
	assert.Equal(t, "/v2", sites[0].APIUri)
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
}

func TestLoadPorts(t *testing.T) {
	testPortFile := "./test/config.json"
	
	ports, err := model.RetrieveConfig(testPortFile)
	assert.Nil(t, err)
	assert.NotNil(t, ports)
	assert.Equal(t, 2, len(ports))
	assert.Equal(t, "", ports[0].Address)
	assert.Equal(t, int64(10099), ports[0].Port)
	assert.Equal(t, "http", ports[0].Protocol)
	assert.Equal(t, "./data/data.json", ports[0].File)
	assert.Equal(t, "", ports[0].User)
	assert.Equal(t, "", ports[0].Password)
	assert.Equal(t, "/api/json?pretty=true", ports[0].APIUrl)
	assert.Equal(t, true, ports[0].Concat)
	assert.Equal(t, true, ports[0].BeforeApi)
	assert.Equal(t, "", ports[1].Address)
	assert.Equal(t, int64(10100), ports[1].Port)
	assert.Equal(t, "http", ports[1].Protocol)
	assert.Equal(t, "./data/data2.json", ports[1].File)
	assert.Equal(t, "", ports[1].User)
	assert.Equal(t, "", ports[1].Password)
	assert.Equal(t, "/api/json?pretty=true", ports[1].APIUrl)
	assert.Equal(t, true, ports[1].Concat)
	assert.Equal(t, true, ports[1].BeforeApi)
}


func getFileAsString(path string) string {
	buf, err := ioutil.ReadFile("_tests/" + path)
	if err != nil {
		panic(err)
	}

	return string(buf)
}

func getRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getFileReader(file string) *os.File {
	r, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	return r
}

