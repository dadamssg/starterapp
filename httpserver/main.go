package main

import (
	"flag"
	"github.com/dadamssg/starterapp/app"
	"github.com/dadamssg/starterapp/httpserver/user"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	app := starterapp.New(config())
	router := mux.NewRouter().StrictSlash(true)

	user.Connect(app, router)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func config() starterapp.Config {

	config := readConfig()

	return starterapp.Config{
		DbConfig: starterapp.DatabaseConfig{
			User:     config["database_user"],
			Password: config["database_password"],
			Host:     config["database_host"],
			Database: config["database_name"],
		},
	}
}

func readConfig() map[string]string {

	configFile := flag.String("config-path", "/config/starterapp.yml", "location of config file")
	flag.Parse()

	data, err := ioutil.ReadFile(*configFile)
	panicIf(err)

	config := make(map[string]string)

	err = yaml.Unmarshal([]byte(data), &config)
	panicIf(err)

	return config
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
