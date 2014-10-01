//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	log.SetFlags(log.Lshortfile)
	start()
	command_line_interface()
	return
}

func command_line_interface() {
	in := bufio.NewReader(os.Stdin)

	continueCLI := true

	for continueCLI {
		line, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("[CLI] %s", err)
		}

		switch line {
		case "quit\n":
			continueCLI = false
			stopAll()

		default:
			fmt.Printf("Type quit to quit\n")
		}

	}
	return
}
