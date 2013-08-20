package main

import (
	"crypto/elliptic"
	"crypto/sha256"
	"log"
	"strings"
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
	//input, err := in.ReadString('\n')
	//if err != nil {
	//	log.Panic(err)
	//}
	log.Println(hash, err)
	return
}

func InsertRopo(thereHKID HKID, myHKID HKID, path string) {
	c, err := GetCommit(myHKID)
	if err == nil {
		Post(thereHKID, path, c)
	}
}

func InsertDomain(thereHKID HKID, myHKID HKID, path string) {
	//split path in to tag_name and leading_path
	nameSegments := strings.Split(path, "/")
	tag_name := nameSegments[len(nameSegments)]
	leading_path := path[:(len(path) - len(tag_name) - 1)]
	//l := list/tag from leading_path
	b, err := Get(thereHKID, leading_path)
	l := NewListFromBytes(b.Bytes())
	//add new tag to tag/list
	l.add(tag_name, HID(myHKID), "tag")
	//post the modified list
	hid, err := Post(thereHKID, path, l)
	_, _ = hid, err
}

func InitRepo(hkid HKID) error {
	//InitCommit()
	c := NewCommit(sha256.New().Sum(nil), hkid)
	err := PostCommit(c)
	return err
}
func InitDomain(hkid HKID, nameSegment string) error {
	//GenHKID()
	t := NewTag(sha256.New().Sum(nil), "blob", nameSegment, hkid)
	err := PostTag(t)
	return err
}

func GenHKID() HKID {
	privkey := KeyGen()
	PostBlob(elliptic.Marshal(privkey.PublicKey.Curve,
		privkey.PublicKey.X, privkey.PublicKey.Y))
	return GenerateHKID(privkey)
}

/*func initRepo(objecthash HKID, path string) (repoHkid HKID) {
	c, privkey := InitCommit()
	err := PostKey(privkey)
	if err != nil {
		log.Panic(err)
	}
	return c.hkid
}*/
