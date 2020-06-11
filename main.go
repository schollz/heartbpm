package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"time"

	movingaverage "github.com/RobinUS2/golang-moving-average"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/schollz/logger"
	"github.com/tarm/serial"
)

var flagCOM string

func init() {
	flag.StringVar(&flagCOM, "com", "", "set the com port (e.g. COM6)")
}

func main() {
	flag.Parse()
	if flagCOM == "" {
		fmt.Println("must add com, --com")
		os.Exit(1)
	}
	logger.SetLevel("trace")
	run()
}

func run() (err error) {
	c := &serial.Config{Name: flagCOM, Baud: 9600, ReadTimeout: time.Second * 10}
	s, err := serial.OpenPort(c)
	if err != nil {
		err = errors.Wrap(err, "no com port")
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

	ma := movingaverage.New(5)
	lastTime := time.Now()

	go func() {
		for {
			var reply string
			reply, err = read(s)
			if strings.TrimSpace(reply) == "b" {
				diff := time.Since(lastTime).Seconds()
				if diff > 0.06 && diff < 2.0 {
					ma.Add(diff)
					fmt.Println(ma.Avg())
				}
				lastTime = time.Now()
			}
			if err != nil {
				logger.Error(err)
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}()

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "hello",
		})
	})
	router.GET("/bpm", func(c *gin.Context) {
		avg := ma.Avg()
		if avg == 0 {
			avg = 1
		}
		c.JSON(200, gin.H{
			"success": true,
			"bpm":     math.Round(60.0/avg*10) / 10,
		})
	})
	router.Run()
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
