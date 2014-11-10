//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.

package main

import (
	"log"

	"github.com/AaronGoldman/ccfs/services"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
}

func main() {
	start()
	services.Repl()
	return
}
