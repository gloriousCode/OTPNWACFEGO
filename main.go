//go:generate go run -tags generate gen.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/zserge/lorca"
)

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
type counter struct {
	sync.Mutex
	count int
}

type entry struct {
	Name   string
	Secret string
}

type code struct {
	Name    string
	Code    string
	Counter uint64
}

const (
	filePath string = "data.json"
)

var shutdown chan (interface{})
var entries []entry
var codes []*code
var mtx sync.Mutex

func main() {
	shutdown = make(chan (interface{}))
	entries = readJSONFile(filePath)

	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 480, 320, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	go generateCodes(ui)
	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	ui.Bind("getCodes", getAllCodes)

	// Create and bind Go object to the UI
	c := &counter{}
	ui.Bind("counterAdd", c.Add)
	ui.Bind("counterValue", c.Value)

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	ui.Eval(`
		console.log("Hello, world!");
		console.log('Multiple values:', [1, false, {"x":5}]);
	`)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
	close(shutdown)
	log.Println("exiting...")
}

// readJSONFile reads a file and converts the JSON to an Entry type
func readJSONFile(file string) []entry {
	plan, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var data []entry
	err = json.Unmarshal(plan, &data)
	if err != nil {
		panic(err)
	}

	return data
}

func (c *counter) Add(n int) {
	c.Lock()
	defer c.Unlock()
	c.count = c.count + n
}

func (c *counter) Value() int {
	c.Lock()
	defer c.Unlock()
	return c.count
}

func generateCodes(ui lorca.UI) {
	for {
		select {
		case <-shutdown:
			return
		default:
			for _, entry := range entries {
				timer := time.Now()
				counter := uint64(math.Floor(float64(timer.Unix()) / float64(30)))
				// Generate and display codes
				generatedCode, err := totp.GenerateCode(entry.Secret, timer)
				if err != nil {
					panic(err)
				}
				var codeFound bool
				for i := range codes {
					if codes[i].Name == entry.Name {
						codes[i].Code = generatedCode
						codes[i].Counter = counter
						codeFound = true
						break
					}
				}
				if !codeFound {
					codes = append(codes, &code{Name: entry.Name, Code: generatedCode, Counter: counter})
				}
			}
			time.Sleep(time.Second)
		}
	}
}


func getAllCodes() [][]string {
	mtx.Lock()
	var resp [][]string
	for i := range codes {
		update := []string{codes[i].Name, codes[i].Code}
		resp = append(resp,update)
	}
	mtx.Unlock()
	return resp
}
