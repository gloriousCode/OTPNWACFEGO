package main

import (
	"sync"
	"time"
)

const (
	filePath string = "data.json"
)

var shutdown chan (interface{})
var entries []entry
var codes []*code
var mtx sync.Mutex
var timer *countdownTimer
var password string

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
