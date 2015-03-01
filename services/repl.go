//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.

package services

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

//repl = Read, Evaluate, Print, Loop

var ContinueCLI = true

type command struct {
	usage   string
	cmdFunc func(string)
}

var commands = map[string]command{
	"quit":             {"Exits the program", func(s string) { ContinueCLI = false }},
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

type runner interface {
	ID() string
	Running() bool
}

var runners = map[string]runner{}

//Registerrunner adds a runner to runners
func Registerrunner(service runner) {
	runners[service.ID()] = service
}

//DeRegisterrunner removes a runner from runners
func DeRegisterrunner(service runner) {
	delete(runners, service.ID())
}

func Repl() {
	in := bufio.NewReader(os.Stdin)
	ch := make(chan string, 1)
	for ContinueCLI {
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
		for ContinueCLI {
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
						fmt.Printf("%20s: %s\n", token, cmd.usage)
					}
				}
				break label
			case <-time.After(time.Millisecond * 250):
			}
		}
	}
	return
}

type running bool

func (r running) String() string {
	if r {
		return "Running"
	} else {
		return "Not Running"
	}
}

func status(string) {
	for id, service := range runners {
		fmt.Printf("%13s: %s\n", id, running(service.Running()))
	}
}
