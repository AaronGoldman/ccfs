//CCFS a Cryptographically Curated File System binds a cryptographic chain of trust into content names.
package main

import (
	"bufio"
	"github.com/AaronGoldman/ccfs/interfaces/crawler"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"github.com/AaronGoldman/ccfs/services/localfile"
	"github.com/AaronGoldman/ccfs/services/timeout"
	"log"
	"os"
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

	crawler.Start()
	services.Registercontentservice(localfile.Instance)
	services.Registerblobgeter(timeout.Instance)
	//services.Registerblobgeter(appsscript.Instance)
	//services.Registerblobgeter(googledrive.Instance)
	//services.Registerblobgeter(kademliadht.Instance)
	//services.Registerblobgeter(multicast.Instance)

	objects.RegisterGeterPoster(
		services.GetPublicKeyForHkid,
		services.GetPrivateKeyForHkid,
		services.PostKey,
		services.PostBlob,
	)
	parseFlagsAndTakeAction()
	in := bufio.NewReader(os.Stdin)
	_, _ = in.ReadString('\n')
	return
}
