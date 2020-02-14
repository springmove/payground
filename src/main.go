package main

import (
	"flag"
	"github.com/linshenqi/sptty"
)

func main() {
	cfg := flag.String("config", "./config.yml", "--config")
	flag.Parse()

	app := sptty.GetApp()
	app.ConfFromFile(*cfg)

	services := sptty.Services{}

	configs := sptty.Configs{}

	app.AddServices(services)
	app.AddConfigs(configs)

	app.Sptting()
}
