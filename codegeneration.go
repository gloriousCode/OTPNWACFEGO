package main

import (
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

func generateCodes() {
	for {
		select {
		case <-shutdown:
			return
		default:
			if !isLoaded {
				time.Sleep(time.Second)
				continue
			}
			for i := range cfg.Codes {
				generatedCode, err := totp.GenerateCode(cfg.Codes[i].Secret, time.Now())
				if err != nil {
					panic(err)
				}
				var codeFound bool
				timerReset := false
				for j := range codes {
					if codes[j].Name == cfg.Codes[i].Name {
						if codes[j].Code != "" && codes[j].Code != generatedCode && !timerReset {
							timer.Reset(time.Second * 30)
							timerReset = true
						}
						codes[j].Code = generatedCode
						codes[j].Counter = fmt.Sprintf("%ds", timer.TimeRemaining())
						codeFound = true
						break
					}
				}
				if !codeFound {
					codes = append(codes, &code{Name: cfg.Codes[i].Name, Code: generatedCode, Counter: fmt.Sprintf("%ds", timer.TimeRemaining())})
				}
			}
			time.Sleep(time.Second)
		}
	}
}

func getAllCodes() [][]string {
	mtx.Lock()
	var resp [][]string
	if !isLoaded {
		mtx.Unlock()
		return resp
	}
	for i := range codes {
		update := []string{codes[i].Name, codes[i].Code, codes[i].Counter}
		resp = append(resp, update)
	}
	mtx.Unlock()
	return resp
}
