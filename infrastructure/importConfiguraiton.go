package infrastructure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// Configuration of neTTool via Json
type Configuration struct {
	Pcapfile    string        `json:"file"`
	StoredData  bool          `json:"storedData"`
	InterfaceID string        `json:"interfaceID"`
	SnapshotLen int32         `json:"snapshotLen"`
	Promiscuous bool          `json:"promiscuous"`
	Timeout     time.Duration `json:"timeout"`
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
	defer jsonFile.Close()

	//fmt.Println("Successfully Opened json")
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &e.Config)

	//fmt.Println(e.Config.Pcapfile)
	//fmt.Println(e.Config.StoredData)
	//fmt.Println(e.Config.InterfaceID)

	return e.Config
}
