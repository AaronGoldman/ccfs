package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"github.com/AaronGoldman/ccfs/services/appsscript"
	//"github.com/AaronGoldman/ccfs/services/googledrive"
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
		log.Printf("Stoping...")
		stopAll()
	}()

	//This Block registers the services for the object module to use
	objects.RegisterGeterPoster(
		services.GetPublicKeyForHkid,
		services.GetPrivateKeyForHkid,
		services.PostKey,
		services.PostBlob,
	)

	localfile.Start()
	timeout.Start()
	multicast.Start()
	kademliadht.Start()
	_ = appsscript.Start
	//googledrive.Start()
	parseFlagsAndTakeAction()
}
