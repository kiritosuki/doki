package main

import (
	"github.com/kiritosuki/doki/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := app.RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
