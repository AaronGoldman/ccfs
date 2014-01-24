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
		go startFSintegration()
		//hashfindwalk()
		in := bufio.NewReader(os.Stdin)
		_, _ = in.ReadString('\n')
	}

	return
}
