package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	guuid "github.com/google/uuid"
)

func main() {
	now := time.Now()
	logLevels := []string{
		"INFO",
		"WARN",
		"ERROR",
		"UNKNOWN",
	}
	someUUID := guuid.New()
	appName, appNamePresent := os.LookupEnv("DUMMY_APP_NAME")

	if !appNamePresent {
		appName = "default"
	}
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(logLevels)

	//somestring1, _ := regen.Generate("[a-z0-9]{32}")
	log.Printf("%v;%v->ExecutionID:%v@%v\n", appName, logLevels[n], someUUID.String(), now.Unix())
}
