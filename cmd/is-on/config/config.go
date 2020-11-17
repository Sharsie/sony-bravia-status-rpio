package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

const requiredErrorF = `required variable "%s" not set`

func String(name string, defaultValue string, required bool) string {
	val, present := os.LookupEnv(name)

	if !present {

		if required {
			log.Fatalf(requiredErrorF, name)
		}

		val = defaultValue
	}
	return val
}

func Int(name string, defaultValue int, required bool) (val int) {
	sVal, present := os.LookupEnv(name)
	if !present {
		if required {
			log.Fatalf(requiredErrorF, name)
		}
		val = defaultValue
	} else {
		v, err := strconv.ParseInt(sVal, 0, 64)
		if err != nil {
			log.Fatalf(`cannot parse int variable "%s"`, name)
		}
		val = int(v)
	}
	return
}

func Bool(name string, defaultValue bool, required bool) (val bool) {
	sVal, present := os.LookupEnv(name)
	if !present {
		if required {
			log.Fatalf(requiredErrorF, name)
		}
		val = defaultValue
	} else {
		var err error
		val, err = strconv.ParseBool(sVal)
		if err != nil {
			log.Fatalf(`cannot parse bool variable "%s"`, name)
		}
	}
	return
}

func Duration(name string, defaultValue time.Duration, required bool) (val time.Duration) {
	sVal, present := os.LookupEnv(name)
	if !present {
		if required {
			log.Fatalf(requiredErrorF, name)
		}
		val = defaultValue
	} else {
		var err error
		val, err = time.ParseDuration(sVal)
		if err != nil {
			log.Fatalf(`cannot parse duration variable "%s"`, name)
		}
	}
	return
}

var TvApiURL = String("TV_API_URL", "", true)
var GPIOPinNumber = Int("GPIO_PIN_NUMBER", 0, true)

var StatusCheckPeriod = Duration("TV_STATUS_CHECK_PERIOD", 5*time.Second, false)
var SwitchOffFailedAttemptsThreshold = Int("GPIO_SWITCH_OFF_FAILED_ATTEMPTS", 10, false)

var Debug = Bool("DEBUG", false, false)
