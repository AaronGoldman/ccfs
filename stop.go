package main

import (
	"github.com/AaronGoldman/ccfs/interfaces/fuse"
	//"github.com/AaronGoldman/ccfs/objects"
	//"github.com/AaronGoldman/ccfs/services"
	"github.com/AaronGoldman/ccfs/services/appsscript"
	//"github.com/AaronGoldman/ccfs/services/googledrive"
	"github.com/AaronGoldman/ccfs/services/kademliadht"
	"github.com/AaronGoldman/ccfs/services/localfile"
	"github.com/AaronGoldman/ccfs/services/multicast"
	"github.com/AaronGoldman/ccfs/services/timeout"
)

func stopAll() {
	fuse.Stop()
	localfile.Stop()
	timeout.Stop()
	multicast.Stop()
	kademliadht.Stop()
	//googledrive.Stop()
	appsscript.Stop()
}
