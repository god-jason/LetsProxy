package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Cache    string
	Email    string
	ForceRSA bool

	Proxies map[string]string
}

var config = Config{
	Cache:    "certs",
	Email:    "",
	ForceRSA: false,
	Proxies:  nil,
}

func LoadConfig() error{
	log.Println("加载配置")
	filename := configPath
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return SaveConfig()
	} else {
		y, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer y.Close()

		d := yaml.NewDecoder(y)
		return d.Decode(&config)
	}
}

func SaveConfig() error {
	log.Println("保存配置")
	filename := configPath
	y, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755) //os.Create(filename)
	if err != nil {
		return err
	}
	defer y.Close()

	e := yaml.NewEncoder(y)
	defer e.Close()

	return e.Encode(&config)
}
