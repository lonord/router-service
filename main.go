package main

import (
	"log"
)

func main() {
	// config
	err := CheckDepedence()
	if err != nil {
		log.Fatalln(err)
	}
	//
}
