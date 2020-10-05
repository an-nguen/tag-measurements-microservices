package structures

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type ResourceConfig struct {
	// Database info
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	DbName      string `json:"db_name"`
	DbNameUsers string `json:"db_name_users"`
	// Resource server params
	ServerPort  string `json:"server_port"`
	AllowOrigin string `json:"allow_origin"`
	HmacSecret  string `json:"hmac_secret"`
}

/*		Public function: ReadAppConfig
 *		---------------------------
 *      Read json data from a filename in a program filepath and fill ResourceConfig struct
 *
 *		returns: filled ResourceConfig struct
 */
func ReadAppConfig(filename string) ResourceConfig {
	// Get full path to config
	filepath, _ := os.Getwd()
	filepath = strings.ReplaceAll(filepath, "\\", "/")
	// Open that config file
	file, fileErr := os.Open(filepath + filename)
	if fileErr != nil {
		panic(fileErr)
	}
	// Read all from config file
	buffer, readErr := ioutil.ReadAll(file)

	if readErr != nil {
		panic("Read error!")
	}
	// Close file
	_ = file.Close()

	// Parse json string and assign to AppConfig struct
	var config ResourceConfig
	unmarshalErr := json.Unmarshal(buffer, &config)

	if unmarshalErr != nil {
		panic(unmarshalErr)
	}
	return config
}
