package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/customerio/go-customerio"
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

// type data struct {
// 	id         int
// 	created_at string
// 	first_name string
// 	last_name  string
// 	email      string
// 	location   string
// 	hirable    bool
// 	bio        string
// }

// type YourJson struct {
// 	YourSample []struct {
// 		data map[string]string
// 	}
// }

func main() {
	contentConfig, err := ioutil.ReadFile("configuration.json")
	if err != nil {
		log.Fatal("Error while opening the file - ", err)
	}

	var configData Configuration
	err = json.Unmarshal(contentConfig, &configData)
	if err != nil {
		log.Fatal("Error while un-marshalling the data", err)
	}
	fmt.Println(configData.Mappings)

	content, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal("Error while opening the file - ", err)
	}

	var data []map[string]interface{}
	// var data []map[string]interface{}
	//var data YourJson
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatal("Error while un-marshalling the data", err)
	}
	fmt.Println(data[0]["bio"])
	fmt.Println(len(data))

	track := customerio.NewTrackClient("e7192b0752ef138df135", "89f446c3ada96eb5c73b", customerio.WithRegion(customerio.RegionUS))

	// to create people in a workspace with default settings, the id (5) can also be an email address.
	// when creating people using an email address, do not include an email attribute.

	//var clientCustomer make(map[string]string)
	// clientCustomer := make(map[string]string)
	// for _, v := range configData.Mappings {
	// 	//fmt.Println(v.From, v.To)
	// 	clientCustomer[v.From] = v.To
	// }

	//fmt.Println(clientCustomer)

	for _, v := range data {
		temp := make(map[string]interface{})
		var id string
		for _, k := range configData.Mappings {
			if k.From == configData.UserId {
				id = k.From
			}
			temp[k.To] = v[k.From]
		}

		if err := track.Identify(id, temp); err != nil {
			log.Println(err)
		}
	}

	// if err := track.Identify("5", map[string]interface{}{
	// 	"emaila":     "bob@example.com",
	// 	"created_at": time.Now().Unix(),
	// 	"first_name": "Bob-2",
	// 	"plan":       "basic",
	// }); err != nil {
	// 	log.Println(err)
	// }
	//}

}
