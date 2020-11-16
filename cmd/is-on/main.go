package main

import (
	"log"
	"time"

	"github.com/Sharsie/tv-status-rpio/cmd/is-on/config"
	"github.com/Sharsie/tv-status-rpio/cmd/is-on/providers/sony/bravia"
	"github.com/Sharsie/tv-status-rpio/cmd/is-on/version"
	"github.com/stianeikeland/go-rpio/v4"

	"fmt"
)

func main() {
	if version.Tag != "" {
		fmt.Printf("Running version %s\n", version.Tag)
	}

	tvStatusCheck := time.NewTicker(config.TvCheckPeriod)

	err := rpio.Open()

	var tvStatusPin rpio.Pin

	if err != nil {
		log.Println(err)
		if config.Debug {
			log.Println("Proceeding in testing mode without GPIO support")
		} else {
			log.Fatal("Cannot continue without GPIO memory allocation")
		}
	} else {
		defer rpio.Close()
		if config.TvActivePinNumber > 0 {
			tvStatusPin = rpio.Pin(config.TvActivePinNumber)
			tvStatusPin.Output()
		}
	}

	for {
		select {
		case <-tvStatusCheck.C:
			status, err := bravia.IsOn()

			if err != nil {
				log.Println(err)
				break
			}

			if tvStatusPin > 0 {
				if status {
					tvStatusPin.High()
				} else {
					tvStatusPin.Low()
				}
			}

			fmt.Printf("TV is on: %t\n", status)
		}
	}

}
