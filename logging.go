package main

import (
	"io/ioutil"
	"log"
)

func quietLogging(enable bool) {
	if !enable {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
}
