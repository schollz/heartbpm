package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	movingaverage "github.com/RobinUS2/golang-moving-average"
	"github.com/dustin/go-humanize"
	"github.com/schollz/httpfileserver"
	"github.com/schollz/logger"
	"github.com/tarm/serial"
)

var flagCOM string
var flagDebug bool

func init() {
	flag.StringVar(&flagCOM, "com", "", "set the com port (e.g. COM6)")
	flag.BoolVar(&flagDebug, "debug", false, "debug mode")
}

func main() {
	flag.Parse()
	if flagDebug {
		logger.SetLevel("debug")
	} else {
		logger.SetLevel("info")
	}
	run()
}

var ma, ma2 *movingaverage.MovingAverage

func run() (err error) {
	ma = movingaverage.New(20)
	ma2 = movingaverage.New(3)
	var s *serial.Port

	if flagCOM != "" {
		c := &serial.Config{Name: flagCOM, Baud: 9600, ReadTimeout: time.Second * 10}
		s, err = serial.OpenPort(c)
		if err != nil {
			err = fmt.Errorf("no com port: %s", err.Error())
			logger.Error(err)
			return
		}
		defer s.Close()

		csig := make(chan os.Signal, 1)
		signal.Notify(csig, os.Interrupt)
		go func() {
			for sig := range csig {
				logger.Debug("shutdown")
				logger.Debug(sig)
				s.Close()
				os.Exit(1)
			}
		}()

		lastTime := time.Now()
		go func() {
			for {
				var reply string
				reply, err = read(s)
				if strings.Contains(reply, "b") {
					diff := time.Since(lastTime).Seconds()
					if diff > 0.06 && diff < 2.0 {
						ma.Add(diff)
					}
					lastTime = time.Now()
				}
				for _, s := range strings.Fields(reply) {
					if len(s) < 3 {
						continue
					}
					dataPoint0, errC := strconv.ParseFloat(s, 64)
					if errC != nil {
						continue
					}
					ma2.Add(dataPoint0)
				}
				if err != nil {
					logger.Error(err)
					time.Sleep(1 * time.Second)
					continue
				}
			}
		}()

	}

	loadTemplates()
	port := 8054
	logger.Infof("listening on :%d", port)
	http.HandleFunc("/static/", httpfileserver.New("/static/", "static/", httpfileserver.OptionNoCache(true)).Handle())
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now().UTC()
	err := handle(w, r)
	if err != nil {
		logger.Error(err)
		t["main"].Execute(w, Render{HaveCom: flagCOM != ""})
	}
	logger.Debugf("%v %v %v (%s) %s\n", r.RemoteAddr, r.Method, r.URL.Path, r.Header.Get("Accept-Language"), time.Since(t0))
}

var t map[string]*template.Template
var mu sync.Mutex

type Render struct {
	HaveCom bool
}

func loadTemplates() {
	mu.Lock()
	defer mu.Unlock()
	t = make(map[string]*template.Template)
	funcMap := template.FuncMap{
		"beforeFirstComma": func(s string) string {
			ss := strings.Split(s, ",")
			if len(ss) == 1 {
				return s
			}
			if len(ss[0]) > 8 {
				return strings.TrimSpace(ss[0])
			}
			return strings.TrimSpace(ss[0] + ", " + ss[1])
		},
		"humanizeTime": func(t time.Time) string {
			return humanize.Time(t)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"removeSlashes": func(s string) string {
			return strings.TrimPrefix(strings.TrimSpace(strings.Replace(s, "/", "-", -1)), "-location-")
		},
		"removeDots": func(s string) string {
			return strings.TrimSpace(strings.Replace(s, ".", "", -1))
		},
		"minusOne": func(s int) int {
			return s - 1
		},
		"mod": func(i, j int) bool {
			return i%j == 0
		},
		"urlbase": func(s string) string {
			uparsed, _ := url.Parse(s)
			return filepath.Base(uparsed.Path)
		},
		"filebase": func(s string) string {
			_, base := filepath.Split(s)
			base = strings.Replace(base, ".", "", -1)
			return base
		},
		"roundfloat": func(f float64) string {
			return fmt.Sprintf("%2.1f", f)
		},
	}
	b, err := ioutil.ReadFile("templates/main.html")
	if err != nil {
		panic(err)
	}
	t["main"] = template.Must(template.New("base").Funcs(funcMap).Delims("///", "///").Parse(string(b)))
}

func handle(w http.ResponseWriter, r *http.Request) (err error) {
	if logger.GetLevel() == "debug" || logger.GetLevel() == "trace" {
		loadTemplates()
	}

	if r.URL.Path == "/bpm" {
		avg := ma.Avg()
		if avg == 0 {
			avg = 1
		}
		w.Write([]byte(fmt.Sprintf(`{"success":true,"bpm":` + fmt.Sprint(math.Round(60.0/avg*10)/10) + `,"point":` + fmt.Sprint(ma2.Avg()) + `}`)))
	} else {
		t["main"].Execute(w, Render{HaveCom: flagCOM != ""})
	}

	return
}

func write(s *serial.Port, data string) (err error) {
	logger.Debugf("writing '%s'", data)
	_, err = s.Write([]byte(data + "\n"))
	if err != nil {
		return
	}
	err = s.Flush()
	return
}

func read(s *serial.Port) (reply string, err error) {
	logger.Trace("reading")
	for {
		buf := make([]byte, 128)
		var n int
		n, err = s.Read(buf)
		reply += string(buf[:n])
		if bytes.Contains(buf[:n], []byte("\n")) {
			break
		}
		if err != nil {
			break
		}
	}
	logger.Tracef("read '%s'", strings.TrimSpace(reply))
	return
}
