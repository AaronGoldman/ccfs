//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	//"github.com/AaronGoldman/ccfs/interfaces/crawler"
	//"github.com/AaronGoldman/ccfs/objects"
	//"github.com/AaronGoldman/ccfs/services"
	//"github.com/AaronGoldman/ccfs/services/localfile"
	//"github.com/AaronGoldman/ccfs/services/timeout"
	//"github.com/AaronGoldman/ccfs/services/appsscript"
	//"github.com/AaronGoldman/ccfs/services/googledrive"
	//"github.com/AaronGoldman/ccfs/services/kademliadht"
	//"github.com/AaronGoldman/ccfs/services/multicast"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	log.SetFlags(log.Lshortfile)
	start()

	//services.Registerblobgeter(appsscript.Instance)
	//services.Registerblobgeter(googledrive.Instance)
	//services.Registerblobgeter(kademliadht.Instance)
	//services.Registerblobgeter(multicast.Instance)

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
