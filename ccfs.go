//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	log.SetFlags(log.Lshortfile)
	start()
	commandLineInterface()
	return
}

var continueCLI = true

func commandLineInterface() {
	in := bufio.NewReader(os.Stdin)
	ch := make(chan string, 1)
	for continueCLI {

		go func() {
			line, err := in.ReadString('\n')
			if err != nil {
				fmt.Printf("[CLI] %s", err)
			} else {
				ch <- line
			}
		}()
	label:
		for continueCLI {
			select {
			case line := <-ch:
				switch line {
				case "quit\n":
					continueCLI = false
					stopAll()

				default:
					fmt.Printf("Type quit to quit\n")
				}
				break label
			case <-time.After(time.Millisecond * 250):
			}
		}
	}
	return
}
