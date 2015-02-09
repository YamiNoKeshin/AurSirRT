package processors

import (
	"log"
	"aursirrt/src/config"
)

func print(msg string) {
	log.Println("PROCESSORS", msg)

}

func printDebug(msg ...interface {}) {
	if config.Debug {
		log.Println("DEBUG PROCESSORS", msg)
	}
}
