package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tiny-docker"
	app.Usage = `a tiny-docker using grpc`
	app.Commands = []cli.Command{
		run,
	}
	app.Before = func(ctx *cli.Context) error {
		//log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.DebugLevel)
		log.SetOutput(os.Stdout)
		log.SetReportCaller(true)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
