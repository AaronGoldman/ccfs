package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)
	go BlobServerStart()
	go RepoServerStart()
	//hashfindwalk()
	in := bufio.NewReader(os.Stdin)
	_, _ = in.ReadString('\n')
	return
}
