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

var TvHostname = String("TV_HOSTNAME", "", true)

var TvCheckPeriod = Duration("TV_CHECK_PERIOD", 5*time.Second, false)

var TvActivePinNumber = Int("TV_ACTIVE_GPIO_NUMBER", 0, true)

var Debug = Bool("DEBUG", false, false)
