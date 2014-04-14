// multicastservice.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type tagfields struct {
	hkid        HKID
	namesegment string
}

type multicastservice struct {
	conn             *net.UDPConn
	mcaddr           *net.UDPAddr
	waitingforblob   map[string]chan blob
	waitingfortag    map[string]chan tag
	waitingforcommit map[string]chan commit
	waitingforkey    map[string]chan blob
}

func (m multicastservice) GetBlob(h HCID) (b blob, err error) {
	message := fmt.Sprintf("{\"type\":\"blob\", \"hcid\": \"%s\"}", h.Hex())
	blobchannel := make(chan blob, 1)
	m.waitingforblob[h.Hex()] = blobchannel
	m.sendmessage(message)
	select {
	case b = <-blobchannel:
		return b, err

	case <-time.After(150 * time.Millisecond):
		log.Printf("Timing out now")
		return b, fmt.Errorf("GetBlob on Multicast service timed out")
	}

}

func (m multicastservice) GetCommit(h HKID) (c commit, err error) {
	message := fmt.Sprintf("{\"type\":\"commit\",\"hkid\": \"%s\"}", h.String())

	commitchannel := make(chan commit, 1)
	m.waitingforcommit[h.String()] = commitchannel
	m.sendmessage(message)
	select {
	case c = <-commitchannel:
		return c, err

	case <-time.After(150 * time.Millisecond):
		log.Printf("Timing out now")
		return c, fmt.Errorf("GetCommit on Multicast service timed out")
	}

}

func (m multicastservice) GetTag(h HKID, namesegment string) (t tag, err error) {
	message := fmt.Sprintf("{\"type\":\"tag\", \"hkid\": \"%s\", \"namesegment\": \"%s\"}", h.Hex(), namesegment)

	tagchannel := make(chan tag, 1)
	m.waitingfortag[h.Hex()+namesegment] = tagchannel
	m.sendmessage(message)
	select {
	case t = <-tagchannel:
		return t, err

	case <-time.After(150 * time.Millisecond):
		log.Printf("Timing out now")
		return t, fmt.Errorf("GetTag on Multicast service timed out")
	}

}

func (m multicastservice) GetKey(h HKID) (b blob, err error) {
	message := fmt.Sprintf("{\"type\":\"key\",\"hkid\": \"%s\"}", h.Hex())
	m.sendmessage(message)
	keychannel := make(chan blob)
	m.waitingforkey[h.Hex()] = keychannel
	select {
	case b = <-keychannel:
		return b, err

	case <-time.After(12000 * time.Millisecond):
		log.Printf("Timing out now")
		return b, fmt.Errorf("GetKey on Multicast service timed out")
	}

}

func (m multicastservice) listenmessage() (err error) {
	log.Printf("Listenmessage is being called now")
	//It is taking only 256 bytes of data.
	go func() {
		log.Printf("gofunc in Listenmessage is being called now")
		for {
			b := make([]byte, 256)
			_, addr, err := m.conn.ReadFromUDP(b)
			if err != nil {
				log.Printf("multicasterror, %s, \n", err)
				return
			}
			msg := strings.Trim(string(b), "\x00")
			//log.Printf("%s", m.conn.LocalAddr())
			log.Printf("Message that is being called in listen message, %s", msg)
			m.receivemessage(msg, addr)
		}
	}()
	return
}

func (m multicastservice) sendmessage(message string) (err error) {
	b := make([]byte, 256)
	copy(b, message)
	log.Printf("Sent message, %s", message)
	_, err = m.conn.WriteToUDP(b, m.mcaddr)
	if err != nil {
		log.Printf("multicasterror, %s, \n", err)
		return
	}

	return err
}

func (m multicastservice) receivemessage(message string, addr net.Addr) (err error) {
	log.Printf("Received message, %s,\n", message)
	hkid, hcid, typestring, namesegment, url := parseMessage(message)
	log.Println(hkid, hcid, typestring, namesegment)
	if url == "" {
		checkAndRespond(hkid, hcid, typestring, namesegment)
		return nil
	}
	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return err
	}
	url = fmt.Sprintf("http://%s:%d%s", host, 8080, url)

	if typestring == "blob" {
		blobchannel, present := m.waitingforblob[hcid.String()]
		//log.Printf("url is %s", url)
		data, err := m.geturl(url)

		if err == nil {
			if present {
				log.Printf("The hcid is: %s", hcid.String())
				blobchannel <- data
				log.Printf("Now the data is: %s", data)
			} else {
				log.Printf("%s \nis not present in waiting map, \n%s", hcid.String(), m.waitingforblob)
			}
			if blob(data).Hash().Hex() == hcid.Hex() {
				localfileserviceInstance.PostBlob(data)
			}
		} else {
			log.Printf("error: %s", err)
		}
	}

	if typestring == "tag" {
		tagchannel, present := m.waitingfortag[hkid.Hex()+namesegment]
		data, err := m.geturl(url)
		if err != nil {
			log.Printf("Error from geturl in tag, %s", err)
		}
		t, err := TagFromBytes(data)
		if err == nil {
			if present {
				log.Printf("Tag is present")
				tagchannel <- t
			} else {
				log.Printf("%s \n is not present in map \n %s", hkid.Hex()+namesegment, m.waitingfortag)
			}
			if t.Verify() {
				localfileserviceInstance.PostTag(t)
			}
		} else {
			log.Printf("error, %s", err)
			log.Printf("Data, %s", data)
		}
	}

	if typestring == "commit" {
		commitchannel, present := m.waitingforcommit[hkid.String()]
		data, err := m.geturl(url)
		if err != nil {
			log.Printf("Error for geturl in commitchannel is, %s", err)
		} else {
			c, err := CommitFromBytes(data)

			if err == nil {
				if present {
					log.Printf("commit is present")
					commitchannel <- c
				} else {
					log.Printf("commit %s\n is not present, \n%v", hkid.String(), m.waitingforcommit)
				}
				if c.Verify() {
					localfileserviceInstance.PostCommit(c)
				}
			}
		}
	}
	if typestring == "key" {
		keychannel, present := m.waitingforkey[hkid.String()]
		data, err := m.geturl(url)
		if err == nil {
			if present {
				log.Printf("key is present")
				keychannel <- data
			} else {
				log.Printf("key is not present, %s", m.waitingforkey)
			}
			p, err := PrivteKeyFromBytes(data)

			if err != nil && p.Verify() && p.Hkid().Hex() == hkid.Hex() {
				localfileserviceInstance.PostKey(p)
			}

		}
	}

	return err
}

func (m multicastservice) geturl(url string) (data []byte, err error) {
	resp, err := http.Get(url) //Takes the http channel and makes it a channel object
	if err != nil {
		log.Printf("The HTTP Get error is %s", err)
		return data, err

	}
	defer resp.Body.Close() //Do this after return is called
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("The ReadAll error is %s", err)
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
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		return multicastservice{}
	}

	return multicastservice{
		conn:             conn,
		mcaddr:           mcaddr,
		waitingforblob:   map[string]chan blob{},
		waitingfortag:    map[string]chan tag{},
		waitingforcommit: map[string]chan commit{},
		waitingforkey:    map[string]chan blob{},
	}
}

var multicastserviceInstance multicastservice

func init() {
	multicastserviceInstance = multicastservicefactory()
	multicastserviceInstance.listenmessage()
}
