package main

import (
	"bufio"
	golist "container/list"

	"fmt"
	"os"
)

func main() {
	//go BlobServerStart()
	//hashfindwalk()
	in := bufio.NewReader(os.Stdin)
	input, err := in.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println(input)
	return
}

type hkid []byte

func regen(objectsToRegen *golist.List, objecthash hkid, b blob) error {
	return nil
}

func initRepo(objecthash hkid, path string) (repoHkid []byte) {
	return nil
}

func initDomain(objecthash hkid, path string) (domainHkid []byte) {
	return nil
}
