package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
	Mappings    []Mapping `json:"mappings"`
}

type dataSend struct {
	Body map[string]interface{}
	id   string
}

var configData Configuration

var data map[string]interface{}

var readChan chan (dataSend)

func readConfigFile(fileString string) {
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
	//fmt.Println(tken)

	//count := 0

	for decoder.More() {
		//fmt.Println(count)
		//count++
		temp := make(map[string]interface{})
		var Cid string
		decoder.Decode(&data)
		for _, k := range configData.Mappings {
			if k.From == configData.UserId {
				Cid = fmt.Sprint(data[k.From])
				//log.Println(Cid)
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
	//close(readChan)
}

func main() {
	readChan = make(chan dataSend)
	var wg sync.WaitGroup
	readConfigFile("configuration.json")
	children := configData.Parallelism
	fmt.Println("children", children)
	// track := customerio.NewTrackClient("e7192b0752ef138df135", "89f446c3ada96eb5c73b", customerio.WithRegion(customerio.RegionUS))
	for c := 0; c < children; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ds := range readChan {
				//time.Sleep(1000)
				//fmt.Println(ds)
				// if err := track.Identify(ds.id, ds.Body); err != nil {
				// 	log.Println(err)
				// }
				// initialize http client
				client := &http.Client{}
				jsonBody, _ := json.Marshal(ds.Body)
				// set the HTTP method, url, and request body
				req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://track.customer.io/api/v1/customers/%s", ds.id), bytes.NewBuffer(jsonBody))
				if err != nil {
					panic(err)
				}
				site_id_base_64 := base64.StdEncoding.EncodeToString([]byte("e7192b0752ef138df135"))
				api_key_base_64 := base64.StdEncoding.EncodeToString([]byte("89f446c3ada96eb5c73b"))
				bearer := "Bearer " + site_id_base_64 + ":" + api_key_base_64
				// add authorization header to the req
				req.Header.Add("Authorization", bearer)
				// set the request header Content-Type for json
				req.Header.Add("Content-Type", "application/json; charset=utf-8")
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				//fmt.Println("Here")
				fmt.Println("resp:", resp)
			}
			//fmt.Println("Shut down signal received")
		}()
	}

	readDataFile("data.json")
	close(readChan)
	wg.Wait()
	fmt.Println("Waiting")
}
