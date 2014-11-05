package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/AaronGoldman/ccfs/interfaces/crawler"
	"github.com/AaronGoldman/ccfs/interfaces/fuse"
	"github.com/AaronGoldman/ccfs/interfaces/web"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"github.com/AaronGoldman/ccfs/services/appsscript"
	"github.com/AaronGoldman/ccfs/services/googledrive"
	"github.com/AaronGoldman/ccfs/services/kademliadht"
	"github.com/AaronGoldman/ccfs/services/localfile"
	"github.com/AaronGoldman/ccfs/services/multicast"
	"github.com/AaronGoldman/ccfs/services/timeout"
)

func start() {
	go func() { //defining, calling and throwing to a different thread
		ch := make(chan os.Signal, 1) //ch is the name of the channel.
		signal.Notify(ch, os.Interrupt, os.Kill)
		sig := <-ch
		log.Printf("Got signal: %s", sig)
		log.Printf("Stopping...")
		stopAll()
	}()

	//This Block registers the services for the object module to use
	objects.RegisterGeterPoster(
		services.GetPublicKeyForHkid,
		services.GetPrivateKeyForHkid,
		services.PostKey,
		services.PostBlob,
	)

	Flags, Command := parseFlags()
	go localfile.Start()
	go timeout.Start()
	go multicast.Start()

	if *Flags.serve {
		web.Start()
		crawler.Start()
	}
	if *Flags.mount {
		fuse.Start()
	}
	if *Flags.dht {
		kademliadht.Start()
	}
	if *Flags.apps {
		appsscript.Start()
	}
	if *Flags.drive {
		googledrive.Start()
	}

	addCurators(Command)
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
	fmt.Printf("%13s: %s\n", localfile.Instance.ID(), running(localfile.Instance.Running()))
	fmt.Printf("%13s: %s\n", googledrive.Instance.ID(), running(googledrive.Instance.Running()))
	fmt.Printf("%13s: %s\n", appsscript.Instance.ID(), running(appsscript.Instance.Running()))
	fmt.Printf("%13s: %s\n", kademliadht.Instance.ID(), running(kademliadht.Instance.Running()))
	fmt.Printf("%13s: %s\n", multicast.Instance.ID(), running(multicast.Instance.Running()))
}
