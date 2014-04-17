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
	//	action, path, flagged := parseFlags()
	//takeActions(action, path)
	in := bufio.NewReader(os.Stdin)
	_, _ = in.ReadString('\n')
	return
}
