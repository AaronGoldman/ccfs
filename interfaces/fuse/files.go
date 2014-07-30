//files

package fuse

import (
	"log"
	"os"
	//"time"
	//"fmt"
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
	//Mtime		time.
	flags fuse.OpenFlags
	size  uint64
}

// func (f File) String() string {
// 		return fmt.Sprintf("%s", f)
// }

func (f File) Attr() fuse.Attr {
	log.Printf("File attributes requested: %+v", f)
	att := fuse.Attr{
		Inode:  uint64(f.inode),
		Size:   f.size,
		Blocks: f.size / 4096,
		// 	Atime:0001-01-01 00:00:00 +0000 UTC
		// 	Mtime:0001-01-01 00:00:00 +0000 UTC
		// 	Ctime:0001-01-01 00:00:00 +0000 UTC
		// 	Crtime:0001-01-01 00:00:00 +0000 UTC
		Mode: f.permission,
		Uid:  1000, //TODO Uid and Gid shouldn't be hardcoded .CCFS_store
		Gid:  1000,
		// 	Rdev:0
		// 	Nlink:0
		Flags: uint32(f.flags),
	}

	log.Printf("file atributes: %+v", att)
	return att
	// files.go:31: file atributes:
	// {
	// 	Inode:10526737836144204806
	// 	Size:0 Blocks:0 				// size should be
	// 	Atime:0001-01-01 00:00:00 +0000 UTC
	// 	Mtime:0001-01-01 00:00:00 +0000 UTC
	// 	Ctime:0001-01-01 00:00:00 +0000 UTC
	// 	Crtime:0001-01-01 00:00:00 +0000 UTC
	// 	Mode:-rw-r--r--
	// 	Nlink:0
	// 	Uid:0
	// 	Gid:0
	// 	Rdev:0
	// 	Flags:0
	// }
}

func (f File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	log.Println("File ReadAll requested")
	if f.contentHash == nil { // default file
		return []byte("hello, world\n"), nil
	}
	b, blobErr := services.GetBlob(f.contentHash) //
	if blobErr != nil {
		return nil, fuse.ENOENT
	}
	return b, nil
}

//nodeopener interface contains open(). Node may be used for file or directory
func (f File) Open(request *fuse.OpenRequest, response *fuse.OpenResponse, intr fs.Intr) (fs.Handle, fuse.Error) {
	//logRequestObject(request, f)
	//request.dir = 0
	//   O_RDONLY int = os.O_RDONLY // open the file read-only.
	//   O_WRONLY int = os.O_WRONLY // open the file write-only.
	//   O_RDWR   int = os.O_RDWR   // open the file read-write.
	//   O_APPEND int = os.O_APPEND // append data to the file when writing.
	//   O_CREATE int = os.O_CREAT  // create a new file if none exists.
	//   O_EXCL   int = os.O_EXCL   // used with O_CREATE, file must not exist
	//   O_SYNC   int = os.O_SYNC   // open for synchronous I/O.
	//   O_TRUNC  int = os.O_TRUNC  // if possible, truncate file when opened.

	b, blobErr := services.GetBlob(f.contentHash) //
	if blobErr != nil {
		log.Printf("get blob error in opening handel %s", blobErr)
		return nil, fuse.ENOENT
	}

	handle := OpenFileHandle{
		buffer: b,
		parent: f.parent,
		name:   f.name,
		inode:  f.inode,
	}
	f.parent.openHandles[handle.name] = true
	return handle, nil
}
