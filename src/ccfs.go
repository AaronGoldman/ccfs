package main

import (
	//"bufio"
	"crypto/elliptic"
	"log"
	//"os"
)

func main() {
	log.SetFlags(log.Lshortfile)
	//go BlobServerStart()
	//hashfindwalk()
	repoHkid := hkidFromDString("46298148238932964800164113348087938361861245597"+
		"2320097996217675372497646408870646300138355611242482091187065042115198890"+
		"6751710824965155500230480521264034469", 10)
	hash, err := Post(repoHkid, "postedBlob", blob([]byte("Posted Blob")))
	//in := bufio.NewReader(os.Stdin)
	//_, err = in.ReadString('\n')
	//if err != nil {
	//	log.Panic(err)
	//}
	log.Println(hash, err)
	return
}

func InsertRepo(h HKID, path string, foreign_hkid HKID) error {
	//splitpath path -> path/last_segment
	_, err := post(
		h,
		path,
		"commit",
		foreign_hkid,
		"commit")
	return err
}

func InsertDomain(h HKID, path string, foreign_hkid HKID) error {
	_, err := post(
		h,
		path,
		"commit",
		foreign_hkid,
		"tag")
	return err
}

func InitRepo(h HKID, path string) error {
	foreign_hkid := GenHKID()
	err := InsertRepo(h, path, foreign_hkid)
	return err
}

func InitDomain(h HKID, path string) error {
	foreign_hkid := GenHKID()
	err := InsertDomain(h, path, foreign_hkid)
	return err
}

func GenHKID() HKID {
	privkey := KeyGen()
	PostBlob(elliptic.Marshal(privkey.PublicKey.Curve,
		privkey.PublicKey.X, privkey.PublicKey.Y))
	return GenerateHKID(privkey)
}
