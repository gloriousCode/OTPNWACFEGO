//go:generate go run -tags generate gen.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	Counter string
}

const (
	filePath string = "data.json"
)

var shutdown chan (interface{})
var entries []entry
var codes []*code
var mtx sync.Mutex
var timer *SecondsTimer

func main() {
	shutdown = make(chan (interface{}))
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
	defer ui.Close()
	timer = NewSecondsTimer(time.Second * 30)
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
			for i := range entries {
				generatedCode, err := totp.GenerateCode(entries[i].Secret, time.Now())
				if err != nil {
					panic(err)
				}
				var codeFound bool
				timerReset := false
				for j := range codes {
					if codes[j].Name == entries[i].Name {
						if codes[j].Code != "" && codes[j].Code != generatedCode && !timerReset {
							// timer has elapsed! Now reset your timer
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
					codes = append(codes, &code{Name: entries[i].Name, Code: generatedCode, Counter: fmt.Sprintf("%ds", timer.TimeRemaining())})
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
		update := []string{codes[i].Name, codes[i].Code, codes[i].Counter}
		resp = append(resp, update)
	}
	mtx.Unlock()
	return resp
}

type SecondsTimer struct {
	timer *time.Timer
	end   time.Time
}

func NewSecondsTimer(t time.Duration) *SecondsTimer {
	return &SecondsTimer{time.NewTimer(t), time.Now().Add(t)}
}

func (s *SecondsTimer) Reset(t time.Duration) {
	s.timer.Reset(t)
	s.end = time.Now().Add(t)
}

func (s *SecondsTimer) Stop() {
	s.timer.Stop()
}

func (s *SecondsTimer) TimeRemaining() time.Duration {
	return s.end.Sub(time.Now()) / 1000000000
}
