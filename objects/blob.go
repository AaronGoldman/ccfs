package objects

import (
	"crypto/sha256"
	"log"
)

//Blob holds the Data for CCFS content
type Blob []byte

//Hash gets the sha256 of the Blob
func (b Blob) Hash() HCID {
	var h = sha256.New()
	h.Write(b)
	return HCID(h.Sum(make([]byte, 0)))
}

//Bytes gets the data as a []byte
func (b Blob) Bytes() []byte {
	return []byte(b)
}

//Log sends a go string escaped blob to the log
func (b Blob) Log() {
	log.Printf(
		"Posting Blob %s\n-----BEGIN BLOB--------\n%q\n-------END BLOB--------",
		b.Hash(),
		b,
	)
}
