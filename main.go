package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/soh335/nasne"
	"github.com/soh335/nasne/xsrs"
	"github.com/xconstruct/go-pushbullet"
)

var (
	host   = flag.String("host", "", "host of nasne")
	port   = flag.String("port", "64230", "port of nasne for request xsrs")
	apikey = flag.String("apikey", "", "apikey of pushbullet")
)

func main() {
	flag.Parse()

	if err := _main(); err != nil {
		log.Fatal(err)
	}
}

func _main() error {
	if *host == "" {
		return fmt.Errorf("host is required")
	}

	if *port == "" {
		return fmt.Errorf("port is required")
	}

	if *apikey == "" {
		return fmt.Errorf("apikey is required")
	}

	root, err := nasne.GetRecordScheduleList(net.JoinHostPort(*host, *port))
	if err != nil {
		return err
	}

	for _, item := range root.Items {
		if item.ConflictID != "0" {
			if err := _notify(&item); err != nil {
				return err
			}
		}
	}

	return nil
}

func _notify(i *xsrs.Item) error {
	layout1 := "2006-01-02T15:04:05-0700"
	layout2 := "2006-01-02 15:04:05"

	start, err := time.Parse(layout1, i.ScheduledStartDateTime)
	if err != nil {
		return err
	}
	duration, err := strconv.Atoi(i.ScheduledDuration)
	if err != nil {
		return err
	}
	end := start.Add(time.Second * time.Duration(duration))

	n := &pushbullet.Note{
		Body: fmt.Sprintf("[nasne] %s (%s ~ %s) seems to be conflict", i.Title, start.Format(layout2), end.Format(layout2)),
		Type: "note",
	}
	return pushbullet.New(*apikey).Push("/pushes", n)
}
