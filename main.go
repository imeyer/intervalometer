package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

var interval = flag.String("interval", "15s", "The interval to do something.")

func main() {
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	quit := make(chan struct{})
	numbers := make(chan float64)
	data := make([]float64, 0)

	go func() {
		for {
			value, err := reader.ReadString('\n')
			if err != nil {
				glog.Info("Nothing to be read.")
			}
			value = strings.TrimSpace(value)

			num, err := strconv.ParseFloat(value, 64)
			if err != nil {
				glog.Error("Can not convert", num, " to a float: ", err)
			}
			numbers <- num
		}
	}()

	interval, err := time.ParseDuration(*interval)
	if err != nil {
		glog.Fatal(err)
	}

	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			glog.Info(data)
			data = data[:0]
		case number := <-numbers:
			data = append(data, number)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
