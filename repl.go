//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

//repl = Read, Evaluate, Print, Loop

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
	"status":           {"status prints the status page", status},
}

type commander interface {
	ID() string
	Command(string)
}

//Registercommand adds a commander to commands
func Registercommand(service commander, usage string) {
	commands[service.ID()] = command{usage, service.Command}
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
				ch <- line[:len(line)-1] //sends everything but the last char to the channel
			}
		}()
	label:
		for continueCLI {

			select {
			case line := <-ch:
				tokens := strings.SplitN(line, " ", 2)
				cmd, found := commands[tokens[0]]

				if len(tokens) < 2 {
					tokens = append(tokens, " ")
				}

				if found {
					cmd.cmdFunc(tokens[1])
				} else {
					for token, cmd := range commands {
						fmt.Printf("%s\n\t-%s\n", token, cmd.usage)
					}
				}
				break label
			case <-time.After(time.Millisecond * 250):
			}
		}
	}
	return
}
