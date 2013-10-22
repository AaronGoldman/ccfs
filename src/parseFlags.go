package main

import (
	"flag"
	"log"
)

func parseFlags() (*string, *string, *string) {
	var action = flag.String("action", "",
		"createDomain, createRepository, insertDomain, insertRepository")

	var hkidString = flag.String("hkid", "",
		"hkid to be inserted as hex")
	var path = flag.String("newPath", "/", "path where inserting happens")
	flag.Parse()
	return action, hkidString, path
}

func takeActions(action *string, hkidString *string, path *string) {
	if *action != "" {
		log.Println("action", *action)
	}
	if *hkidString != "" {
		log.Println("hkidString", *hkidString)
	}
	if *path != "/" {
		log.Println("path", *path)
	}
}
