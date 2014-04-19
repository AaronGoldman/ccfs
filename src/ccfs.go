//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.
package main

import (
	"bufio"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	log.SetFlags(log.Lshortfile)
	parseFlagsAndTakeAction()
	in := bufio.NewReader(os.Stdin)
	_, _ = in.ReadString('\n')
	return
}
