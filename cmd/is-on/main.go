package main

import (
	"log"
	"time"

	"github.com/Sharsie/tv-status-rpio/cmd/is-on/config"
	"github.com/Sharsie/tv-status-rpio/cmd/is-on/logger"
	"github.com/Sharsie/tv-status-rpio/cmd/is-on/providers/sony/bravia"
	"github.com/Sharsie/tv-status-rpio/cmd/is-on/version"
	"github.com/stianeikeland/go-rpio/v4"

	"fmt"
)

func main() {
	if version.Tag != "" {
		fmt.Printf("Running version %s\n", version.Tag)
	}

	l := logger.Log{}

	var tvStatusPin rpio.Pin

	// Open memory range for GPIO access in /dev/mem
	err := rpio.Open()

	if err != nil {
		log.Println(err)
		if config.Debug {
			log.Println("Proceeding in DEBUG mode without GPIO support")
		} else {
			log.Fatal("Cannot continue without GPIO memory allocation")
		}
	} else {
		defer rpio.Close()
		if config.GPIOPinNumber > 0 {
			// Set the GPIO pin number
			tvStatusPin = rpio.Pin(config.GPIOPinNumber)

			// Switch the pin number fo Output mode
			tvStatusPin.Output()

			l.Debug("Using GPIO pin number %d in OUTPUT mode", tvStatusPin)
		}
	}

	// TODO: Logic should be part of the provider (other television sets might use subscriptions instead of polling)
	tvStatusCheck := time.NewTicker(config.StatusCheckPeriod)

	l.Debug("Running status checks every %s", config.StatusCheckPeriod)
	failedStatusChecks := 0

	// Poll the status periodically
	for {
		select {
		case <-tvStatusCheck.C:
			// Get the status from the TV
			status, err := bravia.IsOn(&l)

			if err != nil {
				l.Debug("Failed to get TV Status: %s", err)

				failedStatusChecks += 1
				if config.SwitchOffFailedAttemptsThreshold > 0 && tvStatusPin > 0 && failedStatusChecks >= config.SwitchOffFailedAttemptsThreshold {
					// Switch off the pin if we cannot reach the TV (completely switched off for example)
					tvStatusPin.Low()
					l.Debug(
						"Received %d failed status checks, keeping the GPIO pin number %d off",
						failedStatusChecks,
						tvStatusPin,
					)
				}
				break
			}

			failedStatusChecks = 0

			if tvStatusPin > 0 {
				if status {
					if tvStatusPin.Read() != rpio.High {
						l.Debug("Switching the GPIO pin %d to HIGH", tvStatusPin)
						tvStatusPin.High()
					}
				} else {
					if tvStatusPin.Read() != rpio.Low {
						l.Debug("Switching the GPIO pin %d to LOW", tvStatusPin)
						tvStatusPin.Low()
					}
				}
			}
		}
	}

}
