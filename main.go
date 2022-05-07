package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	//fmt.Println("%+v", configData)
	fmt.Println("Read Config file completed")
}

// func readDataFile(fileString string) []byte {
// 	content, err := ioutil.ReadFile("data.json")
// 	if err != nil {
// 		log.Fatal("Error while opening the file - ", err)
// 	}

// 	err = json.Unmarshal(content, &data)
// 	if err != nil {
// 		log.Fatal("Error while un-marshalling the data", err)
// 	}
// 	return content
// }

func readDataFile(filePath string) {
	readChan = make(chan dataSend)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Error while opening the file: ", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	tken, err := decoder.Token()
	if err != nil {
		log.Fatalln("Error while decoding a token: ", err)
	}
	fmt.Println(tken)

	count := 0

	for decoder.More() {
		fmt.Println(count)
		count++
		temp := make(map[string]interface{})
		var id string
		decoder.Decode(&data)
		for _, k := range configData.Mappings {
			if k.From == configData.UserId {
				id = fmt.Sprint(data[k.From])
				log.Println(id)
			}
			temp[k.To] = data[k.From]
		}
		ds := dataSend{
			Body: temp,
			id:   id,
		}
		//fmt.Println(ds)
		readChan <- ds
	}
	fmt.Println("Read Data file completed")
	//close(readChan)
}

func main() {
	// var wg sync.WaitGroup
	readConfigFile("configuration.json")
	children := configData.Parallelism
	fmt.Println("children", children)
	for c := 0; c < children; c++ {
		// 	wg.Add(1)
		go func() {
			for ds := range readChan {
				//fmt.Println(ds)
				//defer wg.Done()
				// initialize http client
				client := &http.Client{}
				jsonBody, _ := json.Marshal(ds.Body)
				// set the HTTP method, url, and request body
				req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://track.customer.io/api/v1/customers/%s", ds.id), bytes.NewBuffer(jsonBody))
				if err != nil {
					panic(err)
				}
				var bearer = "Bearer " + "e7192b0752ef138df135:89f446c3ada96eb5c73b"
				// add authorization header to the req
				req.Header.Add("Authorization", bearer)
				// set the request header Content-Type for json
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				fmt.Println(resp.StatusCode)
			}
		}()
	}

	readDataFile("data.json")
	//close(readChan)
	//	wg.Wait()

	//fmt.Println(data[0]["bio"])
	//fmt.Println(len(data))

	// track := customerio.NewTrackClient("e7192b0752ef138df135", "89f446c3ada96eb5c73b", customerio.WithRegion(customerio.RegionUS))

	// for _, v := range data {
	// 	temp := make(map[string]interface{})
	// 	var id string
	// 	for _, k := range configData.Mappings {
	// 		if k.From == configData.UserId {
	// 			id = fmt.Sprint(v[k.From])
	// 			fmt.Println(id)
	// 		}
	// 		temp[k.To] = v[k.From]
	// 	}

	// 	if err := track.Identify(id, temp); err != nil {
	// 		log.Println(err)
	// 	}
	// }

}
