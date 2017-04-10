package ifaces

import (
	"net/http"
	"sync"
)

type WebProcess struct {
	ServerError error
	ServerRef   *http.Server
}

type GateWayModel struct {
	ConfigFile   string
	IndexFile    string
	Status       bool
	Processes    []WebProcess
	IndexProcess WebProcess
	WaitGroup    *sync.WaitGroup
}

