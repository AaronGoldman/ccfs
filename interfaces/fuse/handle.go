//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file
//handle

package fuse

import (
	"fmt"
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

type readHandle struct {
	data []byte
}

func ReadHandle(data []byte) fs.Handle {
	return &readHandle{data}
}

func (h *readHandle) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	log.Printf("ReadHandle ReadAll %s", h)
	return h.data, nil
}

type OpenFileHandle struct {
	buffer []byte
	parent *Dir
	name   string
	inode  fuse.NodeID //we're not using this field yet
}

func (o OpenFileHandle) String() string {
	return fmt.Sprintf(
		"[%d] %s\nid: %v parent: %s\nbuffer: %q",
		len(o.buffer),
		o.name,
		o.inode,
		o.parent,
		o.buffer,
	)
}

//handleReader interface
func (o OpenFileHandle) Read(request *fuse.ReadRequest, response *fuse.ReadResponse, intr fs.Intr) fuse.Error {
	logRequestObject(request, o)
	start := request.Offset
	stop := start + int64(request.Size)
	bufptr := o.buffer
	if stop > int64(len(bufptr)-1) {
		stop = int64(len(bufptr))
	}
	if stop == start {
		log.Printf("no bytes read")
		response.Data = []byte{} //new gives you a pointer
		return nil
	}

	slice := bufptr[start:stop]
	response.Data = slice //address of buffer goes to response
	log.Printf("FileHandle data:%s", response.Data)
	return nil
}

func (o *OpenFileHandle) Write(request *fuse.WriteRequest, response *fuse.WriteResponse, intr fs.Intr) fuse.Error {
	logRequestObject(request, o)
	start := request.Offset
	writeData := request.Data

	if writeData == nil {
		return fuse.ENOENT
	}
	lenData := int(start) + (len(writeData))
	if lenData > int(len(o.buffer)) {
		//set length and capacity of buffer
		var newbfr = make([]byte, (lenData), (lenData))
		copy(newbfr, o.buffer)
		response.Size = copy(newbfr[start:lenData], writeData)
		//log.Printf("before copying to o.buffer: %s", newbfr)
		o.buffer = newbfr

	} else {
		num := copy(o.buffer[start:lenData], writeData)
		response.Size = num
	}
	log.Printf("buffer: %s", o.buffer)
	log.Printf("write response size: %v", response.Size)
	publishErr := o.Publish() //get into loop on parent object
	if publishErr != nil {
		log.Printf("error publish in write(): %+v", publishErr)
		return fuse.EIO
	}
	return nil
}

func (o OpenFileHandle) Release(request *fuse.ReleaseRequest, intr fs.Intr) fuse.Error {
	logRequestObject(request, o)
	request.Respond()
	return nil //fuse.ENOENT
}

//func (o OpenfileHandle)
func (o OpenFileHandle) Flush(request *fuse.FlushRequest, intr fs.Intr) fuse.Error {
	logRequestObject(request, o)
	//o.Publish()
	request.Respond()
	return nil
}

//write out file using postblob

func (o OpenFileHandle) Publish() error { //name=file name
	//log.Printf("buffer contains: %s", o.buffer)
	bfrblob := objects.Blob(o.buffer)
	log.Printf("Posting blob %s\n-----BEGIN BLOB-------\n%s\n-------END BLOB-------", bfrblob.Hash(), bfrblob)
	postblobErr := services.PostBlob(bfrblob)
	if postblobErr != nil {
		return postblobErr
	}
	o.parent.Publish(bfrblob.Hash(), o.name, "blob")
	return postblobErr
}
