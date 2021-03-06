//Copyright 2014 Aaron Goldman. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE file

package fuse

import (
	"fmt"
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

type openFileHandle struct {
	buffer    []byte
	parent    *dir
	file      *file
	name      string
	inode     fuse.NodeID //we're not using this field yet
	publish   bool        //For dir struct only
	handleNum int         //Temp measure
}

func (o openFileHandle) String() string {
	return fmt.Sprintf(
		"Buf. Length\t[%d]\nHan. Name:\t%s\nHan. ID:\t%d\nBuffer:\t%q\nPublish: \t%v\nHandle Parent Information:\n%v\n",
		len(o.buffer),
		o.name,
		o.inode,
		o.buffer,
		o.publish,
		o.parent,
	)
}

//handleReader interface
func (o *openFileHandle) Read(request *fuse.ReadRequest, response *fuse.ReadResponse, intr fs.Intr) fuse.Error {
	log.Printf("request: %+v\nobject: %+v", request, o)
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

func (o *openFileHandle) Write(request *fuse.WriteRequest, response *fuse.WriteResponse, intr fs.Intr) fuse.Error {
	log.Printf("Request: %+v\nObject: %+v", request, o)
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
	log.Printf("Buffer: %s", o.buffer)
	log.Printf("write response size: %v", response.Size)

	return nil
}

func (o openFileHandle) Release(request *fuse.ReleaseRequest, intr fs.Intr) fuse.Error {
	log.Printf("Request: %+v\nObject: %+v\n", request, o)
	if o.parent.openHandles[o.name].handleNum == 1 {
		o.parent.RemoveHandle(o.name)
	} else {
		o.parent.openHandles[o.name].handleNum--
	}
	request.Respond()
	return nil //fuse.ENOENT
}

//func (o OpenfileHandle)
func (o openFileHandle) Flush(request *fuse.FlushRequest, intr fs.Intr) fuse.Error {
	log.Printf("Request: %+v\nHandle Information:\n%+v", request, o)
	if o.parent.openHandles[o.name].publish == true && o.publish == true {
		o.Publish()
		o.parent.openHandles[o.name].publish = false //Prevent multiple publishes
	}
	request.Respond()
	return nil
}

//write out file using postblob

func (o openFileHandle) Publish() error { //name=file name
	//log.Printf("buffer contains: %s", o.buffer)
	bfrblob := objects.Blob(o.buffer)
	log.Printf("Posting blob %s\n-----BEGIN BLOB-------\n%s\n-------END BLOB-------", bfrblob.Hash(), bfrblob)
	postblobErr := services.PostBlob(bfrblob)
	if postblobErr != nil {
		return postblobErr
	}

	//Protection against nil pointer(note that this should not be able to occur)
	if o.file != nil {
		o.file.Update(&o, bfrblob.Hash())
	}
	o.parent.Publish(bfrblob.Hash(), o.name, "blob")
	return postblobErr
}
