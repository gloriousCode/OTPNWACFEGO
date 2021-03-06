package main

import (
	"sync"
	"time"

	"github.com/zserge/lorca"
)

const (
	filePath string = "data.json"
)

var shutdown chan (interface{})
var cfg jsonData
var codes []*code
var mtx sync.Mutex
var timer *countdownTimer
var password string
var ui lorca.UI
var wg sync.WaitGroup
var key string
var isLoaded bool

type jsonData struct {
	PromptEncrypt bool `json:"prompt-encrypt"`
	Codes []entry `json:"codes"`
}

type entry struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type code struct {
	Name    string
	Code    string
	Counter string
}

type countdownTimer struct {
	timer *time.Timer
	end   time.Time
}
