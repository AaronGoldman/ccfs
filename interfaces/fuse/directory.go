//Directory

package fuse

import (
	"fmt"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

type Dir struct {
	leaf         objects.HID
	permission   os.FileMode
	content_type string
	parent       *Dir
	name         string
	openHandles  map[string]bool
	//nodeMap      map[fuse.NodeID]*Dir
	inode fuse.NodeID
}

func (d Dir) Fsync(r *fuse.FsyncRequest, intr fs.Intr) fuse.Error {
	r.Respond() // ????
	return nil
}

//constructor
func (d Dir) newDir(Name string) *Dir {
	log.Printf("newdir called:%s", Name)
	p := Dir{
		leaf:         objects.Blob{}.Hash(),
		permission:   d.permission,
		content_type: "list",
		parent:       &d,
		name:         Name,
		openHandles:  map[string]bool{},
		//nodeMap:      map[fuse.NodeID]*Dir{},
		inode: GenerateInode(d.inode, Name),
	}
	return &p
}

//func (d Dir) Remove(*fuse.RemoveRequest, Intr) fuse.Error
func (d Dir) String() string {
	return fmt.Sprintf(
		"[%s]%s %s\nmode:%s parent:%s\nid:%v",
		d.content_type,
		d.name,
		d.leaf,
		d.permission,
		d.parent,
		d.inode,
	)
}

func logRequestObject(r, o fmt.Stringer) {
	log.Printf("request: %+v", r)
	log.Printf("object: %+v", o)
	return
}

func (d Dir) Rename(
	r *fuse.RenameRequest,
	newDir fs.Node,
	intr fs.Intr,
) fuse.Error {
	logRequestObject(r, d)
	//find content_type
	if r.OldName != r.NewName {
		d.name = r.NewName
	}
	d.name = r.OldName

	switch d.content_type {

	case "list":
		l, listErr := services.GetList(d.leaf.(objects.HCID))
		if listErr != nil {
			return listErr
		}
		newList := l.Add(r.NewName, l[r.OldName].Hash, l[r.OldName].TypeString)
		newList = l.Remove(r.OldName)
		log.Printf(
			"Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------",
			newList.Hash(),
			newList,
		)
		el := services.PostList(newList)
		if el != nil {
			return listErr
		}
		d.Publish(d.leaf.(objects.HCID), d.name, d.content_type)

	case "commit":
		c, CommitErr := services.GetCommit(d.leaf.(objects.HKID))
		if CommitErr != nil {
			return CommitErr
		}
		l, ListErr := services.GetList(c.ListHash)
		if ListErr != nil {
			return ListErr
		}
		newList := l.Add(r.NewName, l[r.OldName].Hash, l[r.OldName].TypeString)
		newList = l.Remove(r.OldName)

		newCommit := c.Update(newList.Hash())
		log.Printf(
			"Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------",
			newList.Hash(),
			newList,
		)
		el := services.PostList(newList)
		if el != nil {
			return ListErr
		}
		log.Printf(
			"Posting commit %s\n-----BEGIN COMMIT-----\n%s\n-------END COMMIT-----",
			newCommit.Hash(),
			newCommit,
		)
		ec := services.PostCommit(newCommit)
		if ec != nil {
			return ListErr
		}

	case "tag":

		oldTag, tagErr := services.GetTag(d.leaf.(objects.HKID), r.OldName)
		var newTag objects.Tag
		if tagErr == nil {
			newTag = objects.NewTag(
				oldTag.HashBytes,
				oldTag.TypeString,
				r.NewName,
				[]objects.HCID{oldTag.Hash()},
				d.leaf.(objects.HKID),
			)
		} else {
			log.Printf("Tag %s\\%s Not Found", d.leaf, d.name)
			return fuse.ENOENT
		}
		log.Printf(
			"Posting tag %s\n-----BEGIN TAG--------\n%s\n-------END TAG--------",
			newTag.Hash(),
			newTag,
		)
		et := services.PostTag(newTag)
		if et != nil {
			return tagErr
		}
	} //end switch

	return nil
}

func (d Dir) Attr() fuse.Attr {
	log.Printf("Directory attributes requested\n\tName:%s", d.name)
	return fuse.Attr{Inode: uint64(d.inode), Mode: os.ModeDir | d.permission}
}

func (d Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	log.Printf("Directory Lookup:\n\tName: %s\n\tHID: %s", name, d.leaf.Hex())
	log.Printf("%v", d)
	new_nodeID := fuse.NodeID(fs.GenerateDynamicInode(uint64(d.inode), name))

	switch d.content_type {
	default:
		return nil, nil
	case "commit":
		node, CommitErr := d.LookupCommit(name, intr, new_nodeID)
		return node, CommitErr
	case "list":
		node, listErr := d.LookupList(name, intr, new_nodeID)
		return node, listErr
	case "tag":
		node, tagErr := d.LookupTag(name, intr, new_nodeID)
		return node, tagErr
	}

}

//creates new file only
func (d Dir) Create(
	request *fuse.CreateRequest,
	response *fuse.CreateResponse,
	intr fs.Intr,
) (fs.Node, fs.Handle, fuse.Error) {
	logRequestObject(request, d)

	//   O_RDONLY int = os.O_RDONLY // open the file read-only.
	//   O_WRONLY int = os.O_WRONLY // open the file write-only.
	//   O_RDWR   int = os.O_RDWR   // open the file read-write.
	//   O_APPEND int = os.O_APPEND // append data to the file when writing.
	//   O_CREATE int = os.O_CREAT  // create a new file if none exists.
	//   O_EXCL   int = os.O_EXCL   // used with O_CREATE, file must not exist
	//   O_SYNC   int = os.O_SYNC   // open for synchronous I/O.
	//   O_TRUNC  int = os.O_TRUNC  // if possible, truncate file when opened.

	node := File{
		contentHash: objects.Blob{}.Hash(),
		permission:  request.Mode, //os.FileMode(0777)
		parent:      &d,
		name:        request.Name,
		inode:       fuse.NodeID(fs.GenerateDynamicInode(uint64(d.inode), request.Name)),
		size:        0,
	}
	handle := OpenFileHandle{
		buffer: []byte{},
		parent: &d,
		name:   request.Name,
		inode:  node.inode,
	}

	switch {
	default:
		log.Printf("unexpected requests.flag error:%+v", request.Flags)
		return nil, nil, fuse.ENOSYS
	// case O_RDONLY, O_RDWR, O_APPEND, O_SYNC: //OPEN FILE if it exists
	// 		if err != nil {
	// 		return nil, fuse.ENOENT
	// 		}
	case os.O_EXCL&int(request.Flags) == os.O_EXCL:
		_, LookupErr := d.Lookup(request.Name, intr)
		if LookupErr == nil {
			return nil, nil, fuse.EEXIST
		}
		fallthrough //-- file doesn't exist
	case os.O_CREATE&int(request.Flags) == os.O_CREATE:
		d.openHandles[handle.name] = true
		return node, handle, nil
		// case O_WRONLY, O_APPEND: //OPEN AN EMPTY FILE
		// 		if err == nil {
		// 	 	return nil, fuse.EEXIST
		// 		 }
		//case O_TRUNC:
		//return ENOSYS
	}
	//return node, handle, nil
}

func (d Dir) Publish(h objects.HCID, name string, typeString string) (err error) {
	switch d.content_type {

	default:
		log.Printf("unknown type: %s", d.content_type)
		return fmt.Errorf("unknown type: %s", d.content_type)
	case "commit":
		c, CommitErr := services.GetCommit(d.leaf.(objects.HKID))
		if CommitErr != nil {
			return CommitErr
		}
		l, listErr := services.GetList(c.ListHash)
		if listErr != nil {
			return listErr
		}
		newList := l.Add(name, h, typeString)

		newCommit := c.Update(newList.Hash())
		log.Printf(
			"Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------",
			newList.Hash(),
			newList,
		)
		el := services.PostList(newList)
		if el != nil {
			return listErr
		}
		log.Printf(
			"Posting commit %s\n-----BEGIN COMMIT-----\n%s\n-------END COMMIT-----",
			newCommit.Hash(),
			newCommit,
		)
		ec := services.PostCommit(newCommit)
		if ec != nil {
			return listErr
		}
		return nil
	case "tag":
		t, tagErr := services.GetTag(d.leaf.(objects.HKID), name)
		var newTag objects.Tag
		if tagErr == nil {
			newTag = t.Update(h, typeString)
		} else {
			log.Printf("Tag %s\\%s Not Found", d.leaf, name)
			newTag = objects.NewTag(h, typeString, name, nil, d.leaf.(objects.HKID))
		}
		log.Printf(
			"Posting tag %s\n-----BEGIN TAG--------\n%s\n-------END TAG--------",
			newTag.Hash(),
			newTag,
		)
		et := services.PostTag(newTag)
		if et != nil {
			return tagErr
		}
		return nil
	case "list":
		l, listErr := services.GetList(d.leaf.(objects.HCID))
		if listErr != nil {
			return listErr
		}
		newList := l.Add(name, h, typeString)
		log.Printf(
			"Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------",
			newList.Hash(),
			newList,
		)
		el := services.PostList(newList)
		if el != nil {
			return listErr
		}
		d.parent.Publish(newList.Hash(), d.name, "list")
		return nil
	}
}

func (d Dir) LookupCommit(name string, intr fs.Intr, nodeID fuse.NodeID) (fs.Node, fuse.Error) {

	c, CommitErr := services.GetCommit(d.leaf.(objects.HKID))
	if CommitErr != nil {
		log.Printf("commit %s:", CommitErr)
		return nil, nil
	}
	//get list hash
	l, listErr := services.GetList(c.ListHash) //l is the list object
	if listErr != nil {
		log.Printf("commit list retieval error %s:", listErr)
		return nil, nil
	}

	list_entry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
	if !present {
		return nil, fuse.ENOENT
	}
	//getKey to figure out permissions of the child
	_, keyErr := services.GetKey(c.Hkid)
	//perm := fuse.Attr{Mode: 0555}//default read permissions
	perm := os.FileMode(0777)

	if keyErr != nil {
		log.Printf("error not nil; change file Mode %s:", keyErr)
		//perm =  fuse.Attr{Mode: 0755}
		perm = os.FileMode(0555)
	}

	if list_entry.TypeString == "blob" {
		b, blobErr := services.GetBlob(list_entry.Hash.(objects.HCID))
		sizeBlob := 0
		if blobErr == nil {
			sizeBlob = len(b)
		}
		return File{
			contentHash: list_entry.Hash.(objects.HCID),
			permission:  perm,
			name:        name,
			parent:      &d,
			inode:       nodeID,
			size:        uint64(sizeBlob),
		}, nil
	}
	ino := fuse.NodeID(1)
	if d.parent != nil {
		ino = GenerateInode(d.parent.inode, name)
	}

	return Dir{
		leaf:         list_entry.Hash,
		permission:   perm,
		content_type: list_entry.TypeString,
		parent:       &d,
		name:         name,
		openHandles:  map[string]bool{},
		inode:        ino,
	}, nil

}

func (d Dir) LookupList(name string, intr fs.Intr, nodeID fuse.NodeID) (fs.Node, fuse.Error) {
	l, listErr := services.GetList(d.leaf.(objects.HCID))
	if listErr != nil {
		log.Printf("get list %s:", listErr)
		return nil, nil
	}
	list_entry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
	if !present {
		return nil, fuse.ENOENT
	}

	b, blobErr := services.GetBlob(list_entry.Hash.(objects.HCID))
	sizeBlob := 0
	if blobErr == nil {
		sizeBlob = len(b)
	}

	if list_entry.TypeString == "blob" {
		return File{
			contentHash: list_entry.Hash.(objects.HCID),
			permission:  d.permission,
			parent:      &d,
			name:        name,
			inode:       nodeID,
			size:        uint64(sizeBlob),
		}, nil
	}
	return Dir{
		leaf:         list_entry.Hash,
		permission:   d.permission,
		content_type: list_entry.TypeString,
		parent:       &d,
		openHandles:  map[string]bool{},
		name:         name,
		inode:        GenerateInode(d.parent.inode, name),
	}, nil

}

func (d Dir) LookupTag(name string, intr fs.Intr, nodeID fuse.NodeID) (fs.Node, fuse.Error) {

	t, tagErr := services.GetTag(d.leaf.(objects.HKID), name) //leaf is HID
	// no blobs because blobs are for file structure

	if tagErr != nil {
		log.Printf("not a tag Err:%s", tagErr)
		return nil, fuse.ENOENT
	}
	//getKey to figure out permissions of the child
	_, keyErr := services.GetKey(t.Hkid)
	perm := os.FileMode(0555) //default read permissions
	if keyErr == nil {
		perm = os.FileMode(0755)
	} else {
		log.Printf("no private key %s:", keyErr)
	}
	if t.TypeString == "blob" {
		b, blobErr := services.GetBlob(t.HashBytes.(objects.HCID))
		sizeBlob := 0
		if blobErr == nil {
			sizeBlob = len(b)
		}
		return File{
			contentHash: t.HashBytes.(objects.HCID),
			permission:  perm,
			name:        name,
			parent:      &d,
			inode:       nodeID,
			size:        uint64(sizeBlob),
		}, nil
	}
	return Dir{
		leaf:         t.HashBytes,
		permission:   perm,
		content_type: t.TypeString,
		parent:       &d,
		openHandles:  map[string]bool{},
		name:         name,
		inode:        GenerateInode(d.parent.inode, name),
	}, nil
}

func (d Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	//log.Printf("ReadDir requested:\n\tName:%s", d.name)
	var l objects.List
	var listErr error
	var dirDirs = []fuse.Dirent{}
	switch d.content_type {
	case "tag":
		//if d.content_type == "tag" {
		tags, tagErr := services.GetTags(d.leaf.(objects.HKID))
		if tagErr != nil {
			log.Printf("tag %s:", tagErr)
			return nil, fuse.ENOENT
		}
		for _, tag := range tags {
			name := tag.NameSegment
			enttype := fuse.DT_Dir
			switch tag.TypeString {
			case "blob":
				enttype = fuse.DT_File
				fallthrough

			case "list", "commit", "tag":
				dirDirs = append(dirDirs, fuse.Dirent{
					Inode: fs.GenerateDynamicInode(uint64(d.inode), name),
					Name:  name,
					Type:  enttype,
				})
			default:
			}
		}
		return dirDirs, nil
	case "commit":
		//} else if d.content_type == "commit" {
		c, CommitErr := services.GetCommit(d.leaf.(objects.HKID))
		if CommitErr != nil {
			log.Printf("commit %s:", CommitErr)
			return nil, fuse.ENOENT
		}
		l, listErr = services.GetList(c.ListHash)
		if listErr != nil {
			log.Printf("commit list %s:", listErr)
			return nil, fuse.ENOENT
		}
	case "list":
		l, listErr = services.GetList(d.leaf.(objects.HCID))
		if listErr != nil {
			log.Printf("list %s:", listErr)
			return nil, fuse.ENOENT
		}
	default:
		return nil, fuse.ENOENT
	}

	for name, entry := range l {
		if entry.TypeString == "blob" {
			append_to_list := fuse.Dirent{
				Inode: fs.GenerateDynamicInode(uint64(d.inode), name),
				Name:  name,
				Type:  fuse.DT_File,
			}
			dirDirs = append(dirDirs, append_to_list)
		} else {
			append_to_list := fuse.Dirent{
				Inode: fs.GenerateDynamicInode(uint64(d.inode), name),
				Name:  name,
				Type:  fuse.DT_Dir}
			dirDirs = append(dirDirs, append_to_list)
		}
	}

	//loop through openHandles
	for openHandle, _ := range d.openHandles {
		inList := false
		for _, dir_entry := range dirDirs {
			if openHandle == dir_entry.Name {
				inList = true
				break
			}
		}
		if !inList {
			dirDirs = append(
				dirDirs,
				fuse.Dirent{
					Inode: fs.GenerateDynamicInode(uint64(d.inode), openHandle),
					Name:  openHandle,
					Type:  fuse.DT_Dir,
				})
		}
	}
	return dirDirs, nil
}
