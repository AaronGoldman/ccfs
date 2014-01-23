//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.
package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)
	action, path, flagged := parseFlags()
	takeActions(action, path)

	if flagged == false {
		go BlobServerStart()
		go RepoServerStart()
		//hashfindwalk()
		in := bufio.NewReader(os.Stdin)
		_, _ = in.ReadString('\n')
	}
	return
}
