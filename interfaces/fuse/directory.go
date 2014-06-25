//Directory

package fuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"fmt"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
	"log"
	"os"
	"os/signal"
)



type Dir struct {
	//path         string
	//trunc        HKID
	//branch       HKID
	leaf         objects.HID
	permission   os.FileMode
	content_type string
	parent       *Dir //should parent be a string
	name         string
	openHandles  map[string]bool
	nodeMap      map[fuse.NodeID]*Dir
	inode        fuse.NodeID
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
		nodeMap:      map[fuse.NodeID]*Dir{},
		inode:        GenerateInode(d.inode, Name),
	}
	return &p
}

func (d Dir) Attr() fuse.Attr {
	log.Printf("Directory attributes requested\n\tName:%s", d.name)
	return fuse.Attr{Inode: uint64(d.inode), Mode: os.ModeDir | d.permission}
}

func (d Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	log.Printf("Directory Lookup:\n\tName: %s\n\tHID: %s", name, d.leaf.Hex())
	new_nodeID := fs.GenerateDynamicInode(uint64(d.inode), name)
	if name == "hello" {
		return File{permission: os.FileMode(0444)}, nil
	}
	//in each case, call
	switch d.content_type {
	default:
		return nil, nil
	case "commit": // a commit has a list hash
		c, err := services.GetCommit(d.leaf.(objects.HKID))
		if err != nil {
			log.Printf("commit %s:", err)
			return nil, nil
		}
		//get list hash
		l, err := services.GetList(c.ListHash()) //l is the list object
		if err != nil {
			log.Printf("commit list retieval error %s:", err)
			return nil, nil
		}

		list_entry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
		if !present {
			return nil, fuse.ENOENT
		}
		//getKey to figure out permissions of the child
		_, err = services.GetKey(c.Hkid())
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
				inode:       new_nodeID,
			}, nil
		}

		return Dir{
			//path:         d.path + "/" + name,
			//trunc:        d.trunc,
			//branch:       d.leaf.(HKID),
			leaf:         list_entry.Hash,
			permission:   perm,
			content_type: list_entry.TypeString,
			parent:       &d,
			name:         name,
			openHandles:  map[string]bool{},
			inode:        GenerateInode(d.parent.inode, name),
		}, nil

	case "list":
		l, err := services.GetList(d.leaf.(objects.HCID))
		if err != nil {
			log.Printf("commit list %s:", err)
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
			}, nil
		}
		return Dir{
			//path: d.path + "/" + name,
			//trunc:        d.trunc,
			//branch:       d.branch,
			leaf:         list_entry.Hash,
			permission:   d.permission,
			content_type: list_entry.TypeString,
			parent:       &d,
			openHandles:  map[string]bool{},
			name:         name,
			inode:        GenerateInode(d.parent.inode, name),
		}, nil
	case "tag":
		t, err := services.GetTag(d.leaf.(objects.HKID), name) //leaf is HID
		// no blobs because blobs are for file structure

		if err != nil {
			log.Printf("not a tag Err:%s", err)
			return nil, fuse.ENOENT
		}
		//getKey to figure out permissions of the child
		_, err = services.GetKey(t.Hkid())
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
				inode:       new_nodeID,
			}, nil
		}
		return Dir{
			//path:         d.path + "/" + name,
			//trunc:        d.trunc,
			//branch:       t.hkid,
			leaf:         t.HashBytes,
			permission:   perm,
			content_type: t.TypeString,
			parent:       &d,
			openHandles:  map[string]bool{},
			name:         name,
			inode:        GenerateInode(d.parent.inode, name),
		}, nil

	}

}



func (d Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	log.Printf("ReadDir requested:\n\tName:%s", d.name)
	var l objects.List
	var err error
	var dirDirs = []fuse.Dirent{
		{Inode: fs.GenerateDynamicInode(uint64(d.inode), "hello"),
			Name: "hello",
			Type: fuse.DT_File,
		},
	}

	if d.content_type == "tag" {
		return dirDirs, nil
	} else if d.content_type == "commit" {
		c, err := services.GetCommit(d.leaf.(objects.HKID))
		//log.Printf("hash is: %s", d.leaf)
		if err != nil {
			log.Printf("commit %s:", err)
			return nil, nil
		}
		l, err = services.GetList(c.ListHash())
		if err != nil {
			log.Printf("commit list %s:", err)
			return nil, nil
		}

	} else if d.content_type == "list" {
		l, err = services.GetList(d.leaf.(objects.HCID))
		if err != nil {
			log.Printf("list %s:", err)
			return nil, nil
		}
	} else {
		return nil, nil
	}
	//log.Printf("list map: %s", l)

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



//2 types of nodes for files and directories. So call rename twice?
//Create node (directory)

//creates new file only
func (d Dir) Create(request *fuse.CreateRequest, response *fuse.CreateResponse, intr fs.Intr) (fs.Node, fs.Handle, fuse.Error) {
	log.Printf("create node")
	log.Printf("permission: %s", request.Mode)
	log.Printf("name: %s", request.Name)
	log.Printf("flags: %s", request.Flags)
	node := File{
		contentHash: objects.Blob{}.Hash(),
		permission:  request.Mode,
		parent:      &d,
		name:        request.Name,
		inode:       fs.GenerateDynamicInode(uint64(d.inode), request.Name),
	}
	handle := OpenFileHandle{
		buffer: []byte{},
		parent: &d,
		name:   request.Name,
	}
	d.openHandles[handle.name] = true
	return node, handle, nil
}

//For directory node
//(f File) ---> ?
func (d Dir) Rename(r *fuse.RenameRequest, newDir fs.Node, intr fs.Intr) fuse.Error {
	log.Printf("print request: %s", r)
	log.Printf("rename request: rename dir - %s", d)
	p := d.newDir(r.NewName) //generate pointer
	//p.name = r.OldName 		// what is newDir to be used for?
	//node = GenerateInode(d.inode, name)
	//p := nodeMap(node){}
	if r.OldName != r.NewName {
		//p = newDir(r.NewName)
		p.name = r.NewName
	}
	//p.name = r.NewName
	return nil //ERNOENT
}




func (d Dir) Publish(h objects.HCID, name string, typeString string) (err error) { //name=file name

	switch d.content_type {
	//-----BEGIN LIST-------
	//-------END LIST-------
	//-----BEGIN COMMIT-----
	//-------END COMMIT-----
	//-----BEGIN TAG--------
	//-------END TAG--------
	//-----BEGIN BLOB-------
	//-------END BLOB-------
	default:
		log.Printf("unknown type: %s", d.content_type)
		return fmt.Errorf("unknown type: %s", d.content_type)
	case "commit":
		c, err := services.GetCommit(d.leaf.(objects.HKID))
		if err != nil {
			return err
		}
		l, err := services.GetList(c.ListHash())
		if err != nil {
			return err
		}
		newList := l.Add(name, h, typeString)
		//log.Printf()
		newCommit := c.Update(newList.Hash())
		log.Printf("Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------", newList.Hash(), newList)
		el := services.PostList(newList)
		if el != nil {
			return err
		}
		log.Printf("Posting commit %s\n-----BEGIN COMMIT-----\n%s\n-------END COMMIT-----", newCommit.Hash(), newCommit)
		ec := services.PostCommit(newCommit)
		if ec != nil {
			return err
		}
		return nil
	case "tag":
		//log.Printf("entered tag block\n\thkid:%s\n\tnamesegment:%s", d.leaf, name)
		t, err := services.GetTag(d.leaf.(objects.HKID), name)
		var newTag objects.Tag
		if err == nil {
			newTag = t.Update(h, typeString)
		} else {
			log.Printf("Tag %s\\%s Not Found", d.leaf, name)
			newTag = objects.NewTag(h, typeString, name, d.leaf.(objects.HKID))
		}
		log.Printf("Posting tag %s\n-----BEGIN TAG--------\n%s\n-------END TAG--------", newTag.Hash(), newTag)
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
		log.Printf("Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------", newList.Hash(), newList)
		el := services.PostList(newList)
		if el != nil {
			return err
		}
		d.parent.Publish(newList.Hash(), d.name, "list")
		return nil
	}
}


