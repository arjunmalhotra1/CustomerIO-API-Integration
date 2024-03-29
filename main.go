package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Mapping struct {
	From string `json:"from"`
	To   string `json:"to"`
}
type Configuration struct {
	Parallelism int       `json:"parallelism"`
	UserId      string    `json:"userid"`
	SiteId      string    `json:"site_id"`
	APIKey      string    `json:"api_key"`
	Mappings    []Mapping `json:"mappings"`
}

type dataSend struct {
	Body map[string]interface{}
	id   string
}

var configData Configuration

var data map[string]interface{}
var wg sync.WaitGroup
var readChan chan (dataSend)

const retries = 3
const url = "https://track.customer.io/api/v1/customers/"

// readConfigFile reads the config file.
func readConfigFile(fileString string) {
	contentConfig, err := ioutil.ReadFile(fileString)
	if err != nil {
		log.Fatal("Error while opening the file - ", err)
	}
	err = json.Unmarshal(contentConfig, &configData)
	if err != nil {
		log.Fatal("Error while un-marshalling the data", err)
	}
	fmt.Println("Reading the Config file completed")
}

// readDataFile reads the data file.
func readDataFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Error while opening the file: ", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	_, err = decoder.Token()
	if err != nil {
		log.Fatalln("Error while reading the opening token: ", err)
	}
	for decoder.More() {
		temp := make(map[string]interface{})
		var Cid string
		decoder.Decode(&data)
		for _, k := range configData.Mappings {
			if k.From == configData.UserId {
				Cid = fmt.Sprint(data[k.From])
			}
			temp[k.To] = data[k.From]
		}
		ds := dataSend{
			Body: temp,
			id:   Cid,
		}
		readChan <- ds
	}
	_, err = decoder.Token()
	if err != nil {
		log.Fatalln("Error while reading the ending token: ", err)
	}
	fmt.Println("Reading the data file completed")
}

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

func sendUsingHttp() {
	defer wg.Done()
	for ds := range readChan {
		// initialize http client
		client := &http.Client{}
		jsonBody, _ := json.Marshal(ds.Body)
		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(url+ds.id), bytes.NewBuffer(jsonBody))
		if err != nil {
			panic(err)
		}
		// Encoding the key to base64.
		key := fmt.Sprintf("%s:%s", configData.SiteId, configData.APIKey)
		bearer := "Basic " + base64.StdEncoding.EncodeToString([]byte(key))
		// add authorization header to the req
		req.Header.Add("Authorization", bearer)
		// set the request header Content-Type for json
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		attempt := 0
		resp, err := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			// handling request time out.
			if resp.StatusCode == http.StatusRequestTimeout {
				for attempt < retries {
					resp, err = client.Do(req)
					if resp.StatusCode != http.StatusOK {
						resp.Body.Close()
						log.Printf("Attempt: %d Error for user %s, %v \n", attempt+1, ds.id, resp)
						time.Sleep(backoffSchedule[attempt])
						attempt++
					} else {
						resp.Body.Close()
						log.Printf("Request for user: %s succeeded ", ds.id)
						continue
					}
				}
				if attempt == retries {
					resp.Body.Close()
					log.Printf("Request for user: %s failed", ds.id)
					continue
				}
			}
			// Handling StatusUnauthorized error.
			if resp.StatusCode == http.StatusUnauthorized {
				resp.Body.Close()
				log.Fatalln("Please check the site_id and api_key in the config file.")

			}
			resp.Body.Close()
			log.Printf("Unexpected error for user: %s status code: %d error: %+v \n", ds.id, resp.StatusCode, err)
			continue
		}
		resp.Body.Close()
		log.Printf("Request for user: %s succeeded ", ds.id)
	}
}

var configFile *string
var dataFile *string

func init() {
	configFile = flag.String("config", "", "config file")
	dataFile = flag.String("data", "", "data file")
	flag.Parse()

}

func main() {
	readChan = make(chan dataSend)
	readConfigFile(*configFile)
	children := configData.Parallelism

	for c := 0; c < children; c++ {
		wg.Add(1)
		go sendUsingHttp()
	}

	readDataFile(*dataFile)
	close(readChan)
	wg.Wait()
}
