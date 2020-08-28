package bustime

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Key      string `yaml:"key"`
	Interval int    `yaml:"interval"`
	Url      string `yaml:"url"`
}

func GetConfig() *Config {
	c := Config{}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &c
}
