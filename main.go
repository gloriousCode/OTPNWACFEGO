//go:generate go run -tags generate gen.go

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/zserge/lorca"
)

func main() {
	shutdown = make(chan (interface{}))
	ui = setupLorca()
	defer ui.Close()

	timer = NewSecondsTimer(time.Second * 30)
	go generateCodes()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))
	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
	close(shutdown)
}

func setKey() {
	key = ui.Eval("$\"#key\").val()").String()
	wg.Done()
}

func setupLorca() lorca.UI {
	entries = readJSONFile(filePath)
	height := len(entries) * 100
	if height > 1280 {
		height = 1280
	}
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 360, height, args...)
	if err != nil {
		log.Fatal(err)
	}

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})
	ui.Bind("getCodes", getAllCodes)
	ui.Bind("setKey", setKey)

	return ui
}
