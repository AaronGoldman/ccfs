//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.

package main

import (
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	log.SetFlags(log.Lshortfile)
	start()
	repl()
	return
}
