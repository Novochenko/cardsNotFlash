package main

import (
	"firstRestAPI/internal/apiserver"
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	configPath string
)

func init() {
	//flag.StringVar(&configPath, "config-path", "D:/User/smth/go/goProjects/flesh-carti/guideAPIGolang/configs/apiserver/config.yaml", "path to config file")
	flag.StringVar(&configPath, "config-path", "config.yaml", "path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	err := cleanenv.ReadConfig(configPath, &config)
	//err := cleanenv.ReadConfig("config.yaml", &config)
	if err != nil {
		log.Fatal("couldn't read config: ", err)
	}
	err = apiserver.Start(&config)
	if err != nil {
		log.Fatal(err)
	}

}
