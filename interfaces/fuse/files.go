//files

package fuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
//	"fmt"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"log"
	"os"
)


type File struct {
	contentHash objects.HCID
	permission  os.FileMode
	parent      *Dir
	name        string
	inode       uint64 //fuse.NodeID
}

func (f File) Attr() fuse.Attr {
	log.Println("File attributes requested")
	return fuse.Attr{Inode: 1, Mode: f.permission}
	//log.Println("Attr 0444")
	//return fuse.Attr{Mode: 0444}
}

func (f File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	log.Println("File ReadAll requested")
	if f.contentHash == nil { // default file
		return []byte("hello, world\n"), nil
	}
	b, err := services.GetBlob(f.contentHash) //
	if err != nil {
		return nil, fuse.ENOENT
	}
	return b, nil
}

//nodeopener interface contains open(). Node may be used for file or directory
func (f File) Open(request *fuse.OpenRequest, response *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	b, err := services.GetBlob(f.contentHash) //
	if err != nil {
		return nil, fuse.ENOENT
	}
	handle := OpenFileHandle{buffer: b, parent: f.parent, name: f.name}
	f.parent.openHandles[handle.name] = true
	return handle, nil
}