package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	esapi "github.com/elastic/go-elasticsearch/v7/esapi"
	guuid "github.com/google/uuid"
)

func main() {
	appName, appNamePresent := os.LookupEnv("ES_CUSTOM_APP_NAME")

	if !appNamePresent {
		appName = "default-app"
	}

	esIndexName, esIndexNamePresent := os.LookupEnv("ES_CUSTOM_INDEX_NAME")

	if !esIndexNamePresent {
		esIndexName = "default-student"
	}

	logLevels := []string{
		"INFO",
		"WARN",
		"ERROR",
		"UNKNOWN",
	}

	cfg := elasticsearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Addresses: []string{
			os.Getenv("ES_URL"),
		},
		Username: os.Getenv("ES_USERNAME"),
		Password: os.Getenv("ES_PASSWORD"),
	}

	es, _ := elasticsearch.NewClient(cfg)

	var (
		wg sync.WaitGroup
	)

	for true {

		// Generate some uid to be used inside the log message.
		someUUID := guuid.New()

		now := time.Now()
		// this is very insecure for cryptographic purposes.
		rand.Seed(now.Unix())
		n := rand.Int() % len(logLevels)

		// Generate log message
		logMessage := fmt.Sprintf(`{"message" : "%v;%v->ExecutionID:%v@%v"}`, appName, logLevels[n], someUUID.String(), now.Unix())

		// also print to sdtout what we are doing.
		log.Println(logMessage)
		wg.Add(1)

		go func(message string) {
			defer wg.Done()

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:   esIndexName,
				Body:    strings.NewReader(logMessage),
				Refresh: "true",
				// We assume here index and pipeline have the same name.
				Pipeline: esIndexName,
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document", res.Status())
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				}
			}
		}(logMessage)
		time.Sleep(time.Second)
	}
	wg.Wait()

}
