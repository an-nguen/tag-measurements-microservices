package structures

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// Application config struct
type NotifyConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}

/*		Public function: ReadAppConfig
 *		---------------------------
 *      Read json data from a filename in a program filepath and fill FetchConfig struct
 *
 *		returns: filled FetchConfig struct
 */
func ReadAppConfig(filename string) NotifyConfig {
	// Get full path to config
	filepath, _ := os.Getwd()
	filepath = strings.ReplaceAll(filepath, "\\", "/")
	file, fileErr := os.Open(filepath + filename)
	if fileErr != nil {
		panic(fileErr)
	}
	buffer, readErr := ioutil.ReadAll(file)

	if readErr != nil {
		panic("Read error!")
	}
	_ = file.Close()

	var config NotifyConfig
	unmarshalErr := json.Unmarshal(buffer, &config)

	if unmarshalErr != nil {
		panic(unmarshalErr)
	}
	return config
}
