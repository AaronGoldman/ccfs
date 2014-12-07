//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

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

//Verbosity
const verbosity = 0 //Change this for more or less output when running CCFS

//Repository is a collection defined by a commit
type repository struct {
	contents   objects.HKID
	inode      fuse.NodeID
	name       string
	parent     *directory
	permission os.FileMode
}

//Domain is a collection defined by a tag
type domain struct {
	contents   objects.HKID
	inode      fuse.NodeID
	name       string
	parent     *directory
	permission os.FileMode
}

//Folder is a collection defined by a list
type folder struct {
	contents   objects.HCID
	inode      fuse.NodeID
	name       string
	parent     *directory
	permission os.FileMode
}

func (f folder) Attr() fuse.Attr {
	return fuse.Attr{
		Inode: uint64(f.inode),
		Mode:  os.ModeDir | f.permission,
	}
}

type directory interface {
	fs.Node
	fs.NodeCreater
	fs.NodeStringLookuper
	fs.NodeRenamer
}

type dir struct {
	leaf        objects.HID
	permission  os.FileMode
	contentType string
	parent      *dir
	name        string
	openHandles map[string]bool
	//nodeMap      map[fuse.NodeID]*Dir
	inode fuse.NodeID
}

func (d dir) Fsync(r *fuse.FsyncRequest, intr fs.Intr) fuse.Error {
	select {
	case <-intr:
		return fuse.EINTR
	default:
	}
	r.Respond() // ????
	return nil
}

//constructor
func (d dir) newDir(Name string) *dir {
	log.Printf("newdir called:%s", Name)
	p := dir{
		leaf:        objects.Blob{}.Hash(),
		permission:  d.permission,
		contentType: "list",
		parent:      &d,
		name:        Name,
		openHandles: map[string]bool{},
		//nodeMap:      map[fuse.NodeID]*Dir{},
		inode: generateInode(d.inode, Name),
	}
	return &p
}

//func (d Dir) Remove(*fuse.RemoveRequest, Intr) fuse.Error
func (d dir) String() string {
	parentName := "nil"
	if d.parent != nil {
		parentName = d.parent.name
	}
	return fmt.Sprintf(
		"Dir. Type:\t%s\nDir. Name:\t%q\nDir. Leaf:\t%s\nDir. Mode:\t%s\nDir. Parent:\t%q\nDir. ID:\t%v\n",
		d.contentType,
		d.name,
		d.leaf,
		d.permission,
		parentName,
		d.inode,
	)
}

func (d dir) Rename(
	r *fuse.RenameRequest,
	newDir fs.Node,
	intr fs.Intr,
) fuse.Error {
	log.Printf("request: %+v\nobject: %+v", r, d)
	select {
	case <-intr:
		return fuse.EINTR
	default:
	}
	//find content_type
	if r.OldName != r.NewName {
		d.name = r.NewName
	}
	d.name = r.OldName

	switch d.contentType {

	case "list":
		l, listErr := services.GetList(d.leaf.(objects.HCID))
		if listErr != nil {
			return listErr
		}
		newList := l.Add(r.NewName, l[r.OldName].Hash, l[r.OldName].TypeString)
		newList = l.Remove(r.OldName)
		//=========================================================================
		if verbosity == 1 {
			log.Printf(
				"Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------",
				newList.Hash(),
				newList,
			)
		}
		//=========================================================================
		el := services.PostList(newList)
		if el != nil {
			return listErr
		}
		d.Publish(d.leaf.(objects.HCID), d.name, d.contentType)

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
		//=========================================================================
		if verbosity == 1 {
			log.Printf(
				"Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------",
				newList.Hash(),
				newList,
			)
		}
		//=========================================================================
		el := services.PostList(newList)
		if el != nil {
			return ListErr
		}
		//=========================================================================
		if verbosity == 1 {
			log.Printf(
				"Posting commit %s\n-----BEGIN COMMIT-----\n%s\n-------END COMMIT-----",
				newCommit.Hash(),
				newCommit,
			)
		}
		//=========================================================================
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
		//=========================================================================
		if verbosity == 1 {
			log.Printf(
				"Posting tag %s\n-----BEGIN TAG--------\n%s\n-------END TAG--------",
				newTag.Hash(),
				newTag,
			)
		}
		//=========================================================================
		et := services.PostTag(newTag)
		if et != nil {
			return tagErr
		}
	} //end switch

	return nil
}

func (d dir) Attr() fuse.Attr {
	log.Printf("Directory Attribute Request: %q\n", d.name)
	return fuse.Attr{
		Inode: uint64(d.inode),
		Mode:  os.ModeDir | d.permission,
		Uid:   uint32(os.Getuid()),
		Gid:   uint32(os.Getgid()),
	}
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

func (d dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	log.Printf("Lookup %q in Directory: %q\n%v\n", name, d.name, d)
	select {
	case <-intr:
		return nil, fuse.EINTR
	default:
	}
	newNodeID := fuse.NodeID(fs.GenerateDynamicInode(uint64(d.inode), name))

	switch d.contentType {
	default:
		return nil, nil
	case "commit":
		node, CommitErr := d.LookupCommit(name, intr, newNodeID)
		return node, CommitErr
	case "list":
		node, listErr := d.LookupList(name, intr, newNodeID)
		return node, listErr
	case "tag":
		node, tagErr := d.LookupTag(name, intr, newNodeID)
		return node, tagErr
	}

}

func (d dir) RemoveHandle(name string) {
	delete(d.openHandles, name)
}

//creates new file only
func (d dir) Create(
	request *fuse.CreateRequest,
	response *fuse.CreateResponse,
	intr fs.Intr,
) (fs.Node, fs.Handle, fuse.Error) {
	log.Printf("Request:\n %+v\nObject: %+v", request, d)
	select {
	case <-intr:
		return nil, nil, fuse.EINTR
	default:
	}
	//   O_RDONLY int = os.O_RDONLY // open the file read-only.
	//   O_WRONLY int = os.O_WRONLY // open the file write-only.
	//   O_RDWR   int = os.O_RDWR   // open the file read-write.
	//   O_APPEND int = os.O_APPEND // append data to the file when writing.
	//   O_CREATE int = os.O_CREAT  // create a new file if none exists.
	//   O_EXCL   int = os.O_EXCL   // used with O_CREATE, file must not exist
	//   O_SYNC   int = os.O_SYNC   // open for synchronous I/O.
	//   O_TRUNC  int = os.O_TRUNC  // if possible, truncate file when opened.

	node := file{
		contentHash: objects.Blob{}.Hash(),
		permission:  os.FileMode(0777), //request.Mode,
		parent:      &d,
		name:        request.Name,
		inode:       generateInode(d.inode, request.Name),
		size:        0,
		flags:       fuse.OpenReadWrite,
	}
	handle := openFileHandle{
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
		response.Flags = 1 << 2
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

func (d dir) Publish(h objects.HCID, name string, typeString string) (err error) {
	switch d.contentType {

	default:
		log.Printf("unknown type: %s", d.contentType)
		return fmt.Errorf("unknown type: %s", d.contentType)
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
		//=========================================================================
		if verbosity == 1 {
			log.Printf(
				"Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------",
				newList.Hash(),
				newList,
			)
		}
		//=========================================================================
		el := services.PostList(newList)
		if el != nil {
			return listErr
		}
		//=========================================================================
		if verbosity == 1 {
			log.Printf(
				"Posting commit %s\n-----BEGIN COMMIT-----\n%s\n-------END COMMIT-----",
				newCommit.Hash(),
				newCommit,
			)
		}
		//=========================================================================
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
		//=========================================================================
		if verbosity == 1 {
			log.Printf(
				"Posting tag %s\n-----BEGIN TAG--------\n%s\n-------END TAG--------",
				newTag.Hash(),
				newTag,
			)
		}
		//=========================================================================
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
		//=========================================================================
		if verbosity == 1 {
			log.Printf(
				"Posting list %s\n-----BEGIN LIST-------\n%s\n-------END LIST-------",
				newList.Hash(),
				newList,
			)
		}
		//=========================================================================
		el := services.PostList(newList)
		if el != nil {
			return listErr
		}
		d.parent.Publish(newList.Hash(), d.name, "list")
		return nil
	}
}

func (d dir) LookupCommit(name string, intr fs.Intr, nodeID fuse.NodeID) (fs.Node, fuse.Error) {
	select {
	case <-intr:
		return nil, fuse.EINTR
	default:

	}
	ino := fuse.NodeID(1)
	if d.parent != nil {
		ino = generateInode(d.parent.inode, name)
	}
	c, CommitErr := services.GetCommit(d.leaf.(objects.HKID))

	if CommitErr != nil {
		return nil, fuse.EIO
		/*
			log.Printf("commit %s:", CommitErr)
			_, err := services.GetKey(d.leaf.(objects.HKID))
			perm := os.FileMode(0555)
			if err == nil {
				perm = 0777
			}
			return dir{
				permission:  perm,
				contentType: "commit",
				leaf:        d.leaf.(objects.HKID),
				parent:      &d,
				name:        name,
				openHandles: map[string]bool{},
				inode:       ino,
			}, nil
		*/
	}
	//get list hash
	l, listErr := services.GetList(c.ListHash) //l is the list object
	if listErr != nil {
		log.Printf("commit list retrieval error %s:", listErr)
		return nil, nil
	}

	listEntry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
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

	if listEntry.TypeString == "blob" {
		b, blobErr := services.GetBlob(listEntry.Hash.(objects.HCID))
		sizeBlob := 0
		if blobErr == nil {
			sizeBlob = len(b)
		}
		return file{
			contentHash: listEntry.Hash.(objects.HCID),
			permission:  perm,
			name:        name,
			parent:      &d,
			inode:       nodeID,
			size:        uint64(sizeBlob),
		}, nil
	}

	ino = fuse.NodeID(1)
	if d.parent != nil {
		ino = generateInode(d.parent.inode, name)
	}

	return dir{
		leaf:        listEntry.Hash,
		permission:  perm,
		contentType: listEntry.TypeString,
		parent:      &d,
		name:        name,
		openHandles: map[string]bool{},
		inode:       ino,
	}, nil

}

func (d dir) LookupList(name string, intr fs.Intr, nodeID fuse.NodeID) (fs.Node, fuse.Error) {
	select {
	case <-intr:
		return nil, fuse.EINTR
	default:
	}
	l, listErr := services.GetList(d.leaf.(objects.HCID))
	if listErr != nil {
		log.Printf("get list %s:", listErr)
		return nil, nil
	}
	listEntry, present := l[name] //go through list entries and is it maps to the string you passed in present == 1
	if !present {
		return nil, fuse.ENOENT
	}

	b, blobErr := services.GetBlob(listEntry.Hash.(objects.HCID))
	sizeBlob := 0
	if blobErr == nil {
		sizeBlob = len(b)
	}

	if listEntry.TypeString == "blob" {
		return file{
			contentHash: listEntry.Hash.(objects.HCID),
			permission:  d.permission,
			parent:      &d,
			name:        name,
			inode:       nodeID,
			size:        uint64(sizeBlob),
		}, nil
	}
	return dir{
		leaf:        listEntry.Hash,
		permission:  d.permission,
		contentType: listEntry.TypeString,
		parent:      &d,
		openHandles: map[string]bool{},
		name:        name,
		inode:       generateInode(d.parent.inode, name),
	}, nil

}

func (d dir) LookupTag(name string, intr fs.Intr, nodeID fuse.NodeID) (fs.Node, fuse.Error) {
	select {
	case <-intr:
		return nil, fuse.EINTR
	default:
	}
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
		return file{
			contentHash: t.HashBytes.(objects.HCID),
			permission:  perm,
			name:        name,
			parent:      &d,
			inode:       nodeID,
			size:        uint64(sizeBlob),
		}, nil
	}
	return dir{
		leaf:        t.HashBytes,
		permission:  perm,
		contentType: t.TypeString,
		parent:      &d,
		openHandles: map[string]bool{},
		name:        name,
		inode:       generateInode(d.parent.inode, name),
	}, nil
}

func (d dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	//log.Printf("ReadDir requested:\n\tName:%s", d.name)
	select {
	case <-intr:
		return nil, fuse.EINTR
	default:
	}
	var l objects.List
	var listErr error
	var dirDirs = []fuse.Dirent{}
	switch d.contentType {
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
			appendToList := fuse.Dirent{
				Inode: fs.GenerateDynamicInode(uint64(d.inode), name),
				Name:  name,
				Type:  fuse.DT_File,
			}
			dirDirs = append(dirDirs, appendToList)
		} else {
			appendToList := fuse.Dirent{
				Inode: fs.GenerateDynamicInode(uint64(d.inode), name),
				Name:  name,
				Type:  fuse.DT_Dir}
			dirDirs = append(dirDirs, appendToList)
		}
	}

	//loop through openHandles
	for openHandle := range d.openHandles {
		inList := false
		for _, dirEntry := range dirDirs {
			if openHandle == dirEntry.Name {
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
