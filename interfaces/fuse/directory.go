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

//constructor
func (d Dir) newDir(Name string) *Dir {
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

func (d Dir) Rename(
	r *fuse.RenameRequest,
	newDir fs.Node,
	intr fs.Intr,
) fuse.Error {
	log.Printf("Rename request: %+v", r)
	log.Printf("rename request: rename dir - %+v", d)
	//find content_type
	if r.OldName != r.NewName {
		d.name = r.NewName
	}
	d.name = r.OldName

	switch d.content_type {

	case "list":
		l, err := services.GetList(d.leaf.(objects.HCID))
		if err != nil {
			return err
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
			return err
		}
		d.Publish(d.leaf.(objects.HCID), d.name, d.content_type)

	case "commit":
		c, err := services.GetCommit(d.leaf.(objects.HKID))
		if err != nil {
			return err
		}
		l, err := services.GetList(c.ListHash)
		if err != nil {
			return err
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
			return err
		}
		log.Printf(
			"Posting commit %s\n-----BEGIN COMMIT-----\n%s\n-------END COMMIT-----",
			newCommit.Hash(),
			newCommit,
		)
		ec := services.PostCommit(newCommit)
		if ec != nil {
			return err
		}

	case "tag":

		oldTag, err := services.GetTag(d.leaf.(objects.HKID), r.OldName)
		var newTag objects.Tag
		if err == nil {
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
			return err
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
	new_nodeID := fuse.NodeID(fs.GenerateDynamicInode(uint64(d.inode), name))

	switch d.content_type {
	default:
		return nil, nil
	case "commit":
		node, err := d.LookupCommit(name, intr, new_nodeID)
		return node, err
	case "list":
		node, err := d.LookupList(name, intr, new_nodeID)
		return node, err
	case "tag":
		node, err := d.LookupTag(name, intr, new_nodeID)
		return node, err
	}

}

//creates new file only
func (d Dir) Create(
	request *fuse.CreateRequest,
	response *fuse.CreateResponse,
	intr fs.Intr,
) (fs.Node, fs.Handle, fuse.Error) {
	log.Printf("create node")
	log.Printf("permission: %s", request.Mode)
	log.Printf("name: %s", request.Name)
	request.Flags = fuse.OpenFlags(os.O_RDWR)
	//request.Dir = 1
	log.Printf("flags: %s", request.Flags)
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
		permission:  request.Mode,
		parent:      &d,
		name:        request.Name,
		inode:       fuse.NodeID(fs.GenerateDynamicInode(uint64(d.inode), request.Name)),
	}
	handle := OpenFileHandle{
		buffer: []byte{},
		parent: &d,
		name:   request.Name,
		inode:  node.inode,
	}
	d.openHandles[handle.name] = true
	return node, handle, nil
}

func (d Dir) Publish(h objects.HCID, name string, typeString string) (err error) {
	switch d.content_type {

	default:
		log.Printf("unknown type: %s", d.content_type)
		return fmt.Errorf("unknown type: %s", d.content_type)
	case "commit":
		c, err := services.GetCommit(d.leaf.(objects.HKID))
		if err != nil {
			return err
		}
		l, err := services.GetList(c.ListHash)
		if err != nil {
			return err
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
			return err
		}
		log.Printf(
			"Posting commit %s\n-----BEGIN COMMIT-----\n%s\n-------END COMMIT-----",
			newCommit.Hash(),
			newCommit,
		)
		ec := services.PostCommit(newCommit)
		if ec != nil {
			return err
		}
		return nil
	case "tag":
		t, err := services.GetTag(d.leaf.(objects.HKID), name)
		var newTag objects.Tag
		if err == nil {
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
			return err
		}
		return nil
	case "list":
		l, err := services.GetList(d.leaf.(objects.HCID))
		if err != nil {
			return err
		}
		newList := l.Add(name, h, typeString)
		log.Printf(
			"Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------",
			newList.Hash(),
			newList,
		)
		el := services.PostList(newList)
		if el != nil {
			return err
		}
		d.parent.Publish(newList.Hash(), d.name, "list")
		return nil
	}
}

func (d Dir) LookupCommit(name string, intr fs.Intr, nodeID fuse.NodeID) (fs.Node, fuse.Error) {

	c, err := services.GetCommit(d.leaf.(objects.HKID))
	if err != nil {
		log.Printf("commit %s:", err)
		return nil, nil
	}
	//get list hash
	l, err := services.GetList(c.ListHash) //l is the list object
	if err != nil {
		log.Printf("commit list retieval error %s:", err)
		return nil, nil
	}

	list_entry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
	if !present {
		return nil, fuse.ENOENT
	}
	//getKey to figure out permissions of the child
	_, err = services.GetKey(c.Hkid)
	//perm := fuse.Attr{Mode: 0555}//default read permissions
	perm := os.FileMode(0555)

	if err != nil {
		log.Printf("error not nil; change file Mode %s:", err)
		//perm =  fuse.Attr{Mode: 0755}
		perm = os.FileMode(0755)
	} else {
		//log.Printf("no private key %s:", err)
	}
	if list_entry.TypeString == "blob" {
		return File{
			contentHash: list_entry.Hash.(objects.HCID),
			permission:  perm,
			name:        name,
			parent:      &d,
			inode:       nodeID,
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
	l, err := services.GetList(d.leaf.(objects.HCID))
	if err != nil {
		log.Printf("get list %s:", err)
		return nil, nil
	}
	list_entry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
	if !present {
		return nil, fuse.ENOENT
	}
	if list_entry.TypeString == "blob" {
		return File{
			contentHash: list_entry.Hash.(objects.HCID),
			permission:  d.permission,
			parent:      &d,
			name:        name,
			inode:       nodeID,
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

	t, err := services.GetTag(d.leaf.(objects.HKID), name) //leaf is HID
	// no blobs because blobs are for file structure

	if err != nil {
		log.Printf("not a tag Err:%s", err)
		return nil, fuse.ENOENT
	}
	//getKey to figure out permissions of the child
	_, err = services.GetKey(t.Hkid)
	perm := os.FileMode(0555) //default read permissions
	if err == nil {
		perm = os.FileMode(0755)
	} else {
		log.Printf("no private key %s:", err)
	}
	if t.TypeString == "blob" {
		return File{
			contentHash: t.HashBytes.(objects.HCID),
			permission:  perm,
			name:        name,
			parent:      &d,
			inode:       nodeID,
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
	log.Printf("ReadDir requested:\n\tName:%s", d.name)
	var l objects.List
	var err error
	var dirDirs = []fuse.Dirent{}

	if d.content_type == "tag" {
		tags, err := services.GetTags(d.leaf.(objects.HKID))
		if err != nil {
			log.Printf("tag %s:", err)
			return nil, fuse.ENOENT
		}
		for _, tag := range tags {
			name := tag.NameSegment
			enttype := fuse.DT_Dir
			if tag.TypeString == "blob" {
				enttype = fuse.DT_File
			}
			dirDirs = append(dirDirs, fuse.Dirent{
				Inode: fs.GenerateDynamicInode(uint64(d.inode), name),
				Name:  name,
				Type:  enttype,
			})
		}
		return dirDirs, nil
	} else if d.content_type == "commit" {
		c, err := services.GetCommit(d.leaf.(objects.HKID))
		if err != nil {
			log.Printf("commit %s:", err)
			return nil, fuse.ENOENT
		}
		l, err = services.GetList(c.ListHash)
		if err != nil {
			log.Printf("commit list %s:", err)
			return nil, fuse.ENOENT
		}

	} else if d.content_type == "list" {
		l, err = services.GetList(d.leaf.(objects.HCID))
		if err != nil {
			log.Printf("list %s:", err)
			return nil, fuse.ENOENT
		}
	} else {
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
			// we need to append this to list + work on the next if(commit/list/tag? )
			// Type for the other one will be fuse.DT_DIR
		} else {
			append_to_list := fuse.Dirent{
				Inode: fs.GenerateDynamicInode(uint64(d.inode), name),
				Name:  name,
				Type:  fuse.DT_Dir}
			dirDirs = append(dirDirs, append_to_list)
		}
	} // end if range
	//	log.Printf("return dirDirs: %s", dirDirs)

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
