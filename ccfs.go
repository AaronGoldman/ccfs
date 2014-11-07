//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.

package main

import (
	"log"

	"github.com/AaronGoldman/ccfs/services"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	log.SetFlags(log.Lshortfile)
	start()
	services.Repl()
	return
}
