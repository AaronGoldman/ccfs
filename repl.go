//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var continueCLI = true

type command struct {
	usage   string
	cmdFunc func(string)
}

var commands = map[string]command{
	"quit":             {"Exits the program", func(s string) { continueCLI = false }},
	"createDomain":     {"createDomain Path", func(s string) {}},
	"createRepository": {"createRepository Path", func(s string) {}},
	"insertDomain":     {"insertDomain Path", func(s string) {}},
	"insertRepository": {"insertRepository Path", func(s string) {}},
	"status":           {"status prints the status page", func(s string) {}},
}

interface commander{
	ID func() string
	command func(string)
}

//Registercommand adds a commander to commands
func Registercommand(service commander, usage string) {
	commands[service.ID()] = command{usage, service.command}
}

//DeRegistercommand removes a commander from commands
func DeRegistercommand(service commander) {
	delete(commands, service.ID())
}

func repl() {
	in := bufio.NewReader(os.Stdin)
	ch := make(chan string, 1)
	for continueCLI {
		fmt.Printf("CCFS:$> ")
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
				tokens := strings.SplitN(line, " ", 2)
				cmd, found := commands[tokens[0]]

				if found && len(tokens) > 1 {
					cmd.cmdFunc(tokens[1])
				} else {
					for token, cmd := range commands {
						fmt.Printf("%s\n\t-%s\n", token, cmd.usage)
					}
				}

				//	switch line {
				//	case "quit\n":
				//		continueCLI = false
				//	case "createDomain\n":
				//		fmt.Printf("Usage: createDomain Path\n")
				//	case "createRepository\n":
				//		fmt.Printf("Usage: createRepository Path\n")
				//	case "insertDomain\n":
				//		// ID Path HKID (Hex)
				//		fmt.Printf("Usage: insertDomain Path HKID(Hex)\n")
				//	case "insertRepository\n":
				//		// IR Path HKID (Hex)
				//		fmt.Printf("Usage: insertRepository Path HKID(Hex)\n")
				//	case "insertKey\n":
				//		// Should print out HKID of the new key
				//		fmt.Printf("Usage: insertKey key(HEX)\n")
				//	case "status\n":
				//		// This prints out the status of the services
				//		fmt.Printf("Usage: status prints the status page\n")
				//	default:
				//		fmt.Printf(`Type quit to quit
				//createDomain Creates a new domain at path
				//createRepository Creates a new repository at path
				//insertDomain Inserts the domain HKID at path
				//insertRepository Inserts the repository HKID at path
				//`)
				//}
				break label
			case <-time.After(time.Millisecond * 250):
			}

		}
	}
	return
}
