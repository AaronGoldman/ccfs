//files

package fuse

import (
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

type File struct {
	contentHash objects.HCID
	permission  os.FileMode
	parent      *Dir
	name        string
	inode       fuse.NodeID
}

func (f File) Attr() fuse.Attr {
	log.Printf("File attributes requested: %s", f.name)
	return fuse.Attr{Inode: uint64(f.inode), Mode: f.permission}
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
	log.Printf("Open File")
	//request.Flags = fuse.OpenFlags(os.O_RDWR)
	log.Printf("request: %+v", request)
	//request.dir = 0
	//   O_RDONLY int = os.O_RDONLY // open the file read-only.
	//   O_WRONLY int = os.O_WRONLY // open the file write-only.
	//   O_RDWR   int = os.O_RDWR   // open the file read-write.
	//   O_APPEND int = os.O_APPEND // append data to the file when writing.
	//   O_CREATE int = os.O_CREAT  // create a new file if none exists.
	//   O_EXCL   int = os.O_EXCL   // used with O_CREATE, file must not exist
	//   O_SYNC   int = os.O_SYNC   // open for synchronous I/O.
	//   O_TRUNC  int = os.O_TRUNC  // if possible, truncate file when opened.

	b, err := services.GetBlob(f.contentHash) //
	if err != nil {
		return nil, fuse.ENOENT
	}
/*switch request.Flags{
		default:

		case 
		case 

} */

	handle := OpenFileHandle{
		buffer: b,
		parent: f.parent,
		name:   f.name,
		inode:  f.inode,
	}
	f.parent.openHandles[handle.name] = true
	return handle, nil
}
