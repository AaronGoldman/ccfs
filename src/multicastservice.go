// multicastservice.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type tagfields struct {
	hkid        HKID
	namesegment string
}

type multicastservice struct {
	conn             *net.UDPConn
	mcaddr           *net.UDPAddr
	responsechannel  chan response
	waitingforblob   map[string]chan blob
	waitingfortag    map[string]chan tag
	waitingforcommit map[string]chan commit
	waitingforkey    map[string]chan blob
}

func (m multicastservice) GetBlob(h HCID) (b blob, err error) {
	message := fmt.Sprintf("{\"type\":\"blob\", \"hcid\": \"%s\"}", h.Hex())
	m.sendmessage(message)
	blobchannel := make(chan blob)
	m.waitingforblob[h.Hex()] = blobchannel
	b = <-blobchannel
	return b, err

}

func (m multicastservice) GetCommit(h HKID) (c commit, err error) {
	message := fmt.Sprintf("{\"type\":\"commit\",\"hkid\": \"%s\"}", h.Hex())
	m.sendmessage(message)
	commitchannel := make(chan commit)
	m.waitingforcommit[h.Hex()] = commitchannel
	c = <-commitchannel
	return c, err
}

func (m multicastservice) GetTag(h HKID, namesegment string) (t tag, err error) {
	message := fmt.Sprintf("{\"type\":\"tag\", \"hkid\": \"%s\", \"namesegment\": \"%s\"}", h.Hex(), namesegment)
	m.sendmessage(message)
	tagchannel := make(chan tag)
	m.waitingfortag[h.Hex()+namesegment] = tagchannel
	t = <-tagchannel
	return t, err

}

func (m multicastservice) GetKey(h HKID) (b blob, err error) {
	message := fmt.Sprintf("{\"type\":\"key\",\"hkid\": \"%s\"}", h.Hex())
	m.sendmessage(message)
	keychannel := make(chan blob)
	m.waitingforkey[h.Hex()] = keychannel
	b = <-keychannel
	return b, err
}

func (m multicastservice) listenmessage() (err error) {
	//It is taking only 256 bytes of data.
	go func() {
		for {
			b := make([]byte, 256)
			_, _, err := m.conn.ReadFromUDP(b)
			if err != nil {
				log.Printf("multicasterror, %s, \n", err)
				return
			}
			//log.Printf("%s", m.conn.LocalAddr())
			m.receivemessage(string(b), m.conn.LocalAddr())
		}
	}()
	return
}

func (m multicastservice) sendmessage(message string) (err error) {
	b := make([]byte, 256)
	copy(b, message)
	_, err = m.conn.WriteToUDP(b, m.mcaddr)
	if err != nil {
		log.Printf("multicasterror, %s, \n", err)
		return
	}

	return err
}

func (m multicastservice) receivemessage(message string,, addr net.Addr) (err error) {
	log.Printf("Received message, %s,\n", message)
	hkid, hcid, typestring, namesegment := parseMessage(message)
	url := "www.google.com"
	if typestring == "blob" {
		blobchannel, present := m.waitingforblob[hcid.String()]
		if present {

			data, err := m.geturl(url)
			if err == nil {
				blobchannel <- data
			}
		}
	}
	if typestring == "tag" {
		tagchannel, present := m.waitingfortag[hkid.String()+namesegment]
		if present {
			data, err := m.geturl(url)
			t, err := TagFromBytes(data)
			if err == nil {
				tagchannel <- t
			}
		}
	}
	if typestring == "commit" {
		commitchannel, present := m.waitingforcommit[hkid.String()]
		if present {
			data, err := m.geturl(url)
			c, err := CommitFromBytes(data)
			if err == nil {
				commitchannel <- c
			}
		}
	}
	if typestring == "key" {
		keychannel, present := m.waitingforkey[hcid.String()]
		if present {
			data, err := m.geturl(url)
			if err == nil {
				keychannel <- data
			}
		}
	}
	log.Printf("HCID message, %s,\n", hcid.String())
	log.Printf("HKID message, %s,\n", hkid.String())
	log.Printf("typestring message, %s,\n", typestring)
	log.Printf("namesegment message, %s,\n", namesegment)
	//parse message
	//if in waiting map send on channel

	return err
}

func (m multicastservice) geturl(url string) (data []byte, err error) {
	resp, err := http.Get(url) //Takes the http channel and makes it a channel object
	if err != nil {
		return data, err
	}
	defer resp.Body.Close() //Do this after return is called
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	} else {
		return data, nil
	}

}

func multicastservicefactory() (m multicastservice) {
	mcaddr, err := net.ResolveUDPAddr("udp", "224.0.1.20:5354")
	if err != nil {
		return multicastservice{}
	}

	conn, err := net.ListenMulticastUDP("udp", nil, mcaddr)
	if err != nil {
		return multicastservice{}
	}

	return multicastservice{conn: conn, mcaddr: mcaddr}
}

func init() {
	multicastserviceInstance := multicastservicefactory()
	multicastserviceInstance.listenmessage()

}

type response struct {
	typestring  string
	hkid        HKID
	hcid        HCID
	namesegment string
	url         string
}

type responseblob struct {
	hcid HCID
	url  string
}

type responsecommit struct {
	hkid HKID
	url  string
}

type responsetag struct {
	hkid        HKID
	namesegment string
	url         string
}

type responsekey struct {
	hkid HKID
	url  string
}
