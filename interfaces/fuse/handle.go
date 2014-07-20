//handle

package fuse

import (
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/AaronGoldman/ccfs/objects"
	"github.com/AaronGoldman/ccfs/services"
)

type OpenFileHandle struct {
	buffer []byte
	parent *Dir
	name   string
	inode  fuse.NodeID //we're not using this field yet
}

//handleReader interface
func (o OpenFileHandle) Read(request *fuse.ReadRequest, response *fuse.ReadResponse, intr fs.Intr) fuse.Error {
	log.Printf("FileHandle Read requested")
	log.Printf("Read request is: %+v", request)
	log.Printf("Read request handle is: %+v", o)
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

	//log.Printf("start:%d", start)
	//log.Printf("stop:%d", stop)
	//log.Printf("length of buffer:%d", len(bufptr))
	slice := bufptr[start:stop]
	response.Data = slice //address of buffer goes to response
	log.Printf("FileHandle data:%s", response.Data)
	return nil
}

func (o *OpenFileHandle) Write(request *fuse.WriteRequest, response *fuse.WriteResponse, intr fs.Intr) fuse.Error {
	log.Printf("Handle Write handle: %+v", o)
	log.Printf("Handle Write request: %+v", request)
	start := request.Offset
	writeData := request.Data
	//log.Printf("start:%d", start)
	//log.Printf("length of write data:%d", len(writeData))
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
	err := o.Publish() //get into loop on parent object
	if err != nil {
		log.Printf("error publish in write(): %+v", err)
		return fuse.EIO
	}
	return nil
}

func (o OpenFileHandle) Release(request *fuse.ReleaseRequest, intr fs.Intr) fuse.Error {
	log.Printf("FileHandle Release requested:\n\tHandle: %+v", o)
	log.Printf("FileHandle Release requested:\n\tRequest: %+v", request)
	request.Respond()
	return nil //fuse.ENOENT
}

//func (o OpenfileHandle)
func (o OpenFileHandle) Flush(request *fuse.FlushRequest, intr fs.Intr) fuse.Error {
	log.Printf("FileHandle Flush requested:\n\tHandle: %+v", o)
	log.Printf("FileHandle Flush requested:\n\tRequest: %+v", request)
	//o.Publish()
	request.Respond()
	return nil
}

//write out file using postblob

func (o OpenFileHandle) Publish() error { //name=file name
	//log.Printf("buffer contains: %s", o.buffer)
	bfrblob := objects.Blob(o.buffer)
	log.Printf("Posting blob %s\n-----BEGIN BLOB-------\n%s\n-------END BLOB-------", bfrblob.Hash(), bfrblob)
	err := services.PostBlob(bfrblob)
	if err != nil {
		return err
	}
	o.parent.Publish(bfrblob.Hash(), o.name, "blob")
	return err
}
