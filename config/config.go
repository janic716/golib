package config

import (
	"encoding/json"
	"io/ioutil"
)

func LoadJsonConf(configFile string, conf interface{}) (err error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, conf)
	return
}
