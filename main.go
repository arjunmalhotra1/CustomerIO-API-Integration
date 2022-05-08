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

func readConfigFile(fileString string) {
	//fmt.Println(fileString)
	contentConfig, err := ioutil.ReadFile(fileString)
	if err != nil {
		log.Fatal("Error while opening the file - ", err)
	}
	err = json.Unmarshal(contentConfig, &configData)
	if err != nil {
		log.Fatal("Error while un-marshalling the data", err)
	}
	fmt.Println("Read Config file completed")
}

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
	fmt.Println("Read Data file completed")
}

func sendUsingHttp() {
	defer wg.Done()
	for ds := range readChan {
		// initialize http client
		client := &http.Client{}
		jsonBody, _ := json.Marshal(ds.Body)
		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://track.customer.io/api/v1/customers/%s", ds.id), bytes.NewBuffer(jsonBody))
		if err != nil {
			panic(err)
		}
		key := fmt.Sprintf("%s:%s", configData.SiteId, configData.APIKey)
		bearer := "Basic " + base64.StdEncoding.EncodeToString([]byte(key))
		// add authorization header to the req
		req.Header.Add("Authorization", bearer)
		// set the request header Content-Type for json
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		fmt.Println("resp:", resp.Status)
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
