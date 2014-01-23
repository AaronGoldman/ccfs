// multicastservice.go
package main

import (
	"fmt"
	"log"
	"net"
)

type multicastservice struct {
	conn   *net.UDPConn
	mcaddr *net.UDPAddr
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
	return multicastservice{conn, mcaddr}

}

func init() {
	multicastserviceInstance := multicastservicefactory()
	multicastserviceInstance.listenmessage()
	multicastserviceInstance.getBlob(HCID{})
	multicastserviceInstance.getCommit(HKID{})
	multicastserviceInstance.getTag(HKID{}, "Hey")
	multicastserviceInstance.getKey(HKID{})

}
