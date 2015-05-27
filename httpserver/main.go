package main

import (
	"flag"
	"github.com/codegangsta/negroni"
	"github.com/dadamssg/starterapp/app"
	"github.com/dadamssg/starterapp/httpserver/user"
	"github.com/gorilla/mux"
	"log"
	"os"
)

func main() {
	configFile := flag.String("config-path", "/config/starterapp.yml", "location of config file")
	flag.Parse()

	config := starterapp.ReadConfig(*configFile)
	app := starterapp.New(config)
	router := mux.NewRouter().StrictSlash(true)

	user.Connect(app, router)

	logger := &negroni.Logger{log.New(os.Stdout, "[starterapp] ", 0)}
	n := negroni.New(negroni.NewRecovery(), logger)
	n.UseHandler(router)
	n.Run(":8080")
}
