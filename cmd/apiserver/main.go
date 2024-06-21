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
	flag.StringVar(&configPath, "config-path", "D:/User/smth/go/goProjects/flesh-cards-backup/configs/apiserver/config.yaml", "path to config file")
	// flag.StringVar(&configPath, "config-path", "C:/Users/Alex/Desktop/cardsNotFlesh/configs/apiserver/config.yaml", "path to config file")
	//flag.StringVar(&configPath, "config-path", "config.yaml", "path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	err := cleanenv.ReadConfig(configPath, &config)
	//err := cleanenv.ReadConfig("config.yaml", &config)
	if err != nil {
		log.Fatal("couldn't read config: ", err)
	}
	log.Println(config.LocalHostMode, config.DatabaseURL.FullName, config.DatabaseURL.Host, config.DatabaseURL.Password, config.RedisURL, config.BindAddr)
	err = apiserver.Start(&config)
	if err != nil {
		log.Fatal(err)
	}
}
