package infrastructure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Configuration of neTTool via Json
type Configuration struct {
	Pcapfile string `json:"file"`
}

// ConfigurationFromFS load config from FS
type ConfigurationFromFS struct {
	Config Configuration
}

// LoadConfig do the loading
func (e ConfigurationFromFS) LoadConfig() Configuration {
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func(jsonFile *os.File) {
		err = jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	//fmt.Println("Successfully Opened json")
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &e.Config)
	if err != nil {
		return Configuration{}
	}

	//fmt.Println(e.Config.Pcapfile)
	//fmt.Println(e.Config.UseStoredData)
	//fmt.Println(e.Config.InterfaceID)

	return e.Config
}
