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
				case "createDomain\n":
					fmt.Printf("Usage: createDomain Path")
				case "createRepository\n":
					fmt.Printf("Usage: createRepository Path")
				case "insertDomain\n":
					// ID Path HKID (Hex)
					fmt.Printf("Usage: insertDomain Path HKID(Hex)")
				case "insertRepository\n":
					// IR Path HKID (Hex)
					fmt.Printf("Usage: insertRepository Path HKID(Hex)")
				case "insertKey\n":
					// Should print out HKID of the new key
					fmt.Printf("Usage: insertKey key(HEX)")
				case "status\n":
					// This prints out the status of the services
					fmt.Printf("Usage: status prints the status page")
				default:
					fmt.Printf(`Type quit to quit
createDomain Creates a new domain at path
createRepository Creates a new repository at path
insertDomain Inserts the domain HKID at path
insertRepository Inserts the repository HKID at path
`)
				}
				break label
			case <-time.After(time.Millisecond * 250):
			}

		}
	}
	return
}
