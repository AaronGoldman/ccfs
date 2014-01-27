// multicastservice.go
package main

import (
	"fmt"
	"log"
	"net"
)

type multicastservice struct {
	conn            *net.UDPConn
	mcaddr          *net.UDPAddr
	responsechannel chan response
	waiting         map[hid]chan response
}

func (m multicastservice) getBlob(h HCID) (b blob, err error) {
	message := fmt.Sprintf("{\"type\":\"blob\", \"HCID\": \"%s\"}", h.Hex())
	m.sendmessage(message)
	return b, err

}

func (m multicastservice) getCommit(h HKID) (c commit, err error) {
	message := fmt.Sprintf("{\"type\":\"commit\",\"HKID\": \"%s\"}", h.Hex())
	m.sendmessage(message)
	return c, err
}

func (m multicastservice) getTag(h HKID, namesegment string) (t tag, err error) {
	message := fmt.Sprintf("{\"type\":\"tag\", \"HKID\": \"%s\", \"namesegment\": \"%s\"}", h.Hex(), namesegment)
	m.sendmessage(message)
	return t, err

}

func (m multicastservice) getKey(h HKID) (b blob, err error) {
	message := fmt.Sprintf("{\"type\":\"key\",\"HKID\": \"%s\"}", h.Hex())
	m.sendmessage(message)
	//add channel to waiting map
	//wait on channel
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
			m.receivemessage(string(b))
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

func (m multicastservice) receivemessage(message string) (err error) {
	log.Printf("Received message, %s,\n", message)
	//parse message
	//if in waiting map send on channel

	return err
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

	responsechannel := make(chan response)

	return multicastservice{conn, mcaddr}
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