package caps


import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Sender struct {
	Ctx        context.Context
	SentPacket chan *Packet
}

var targets Targets

type Packet struct {
	Timestamp         int64
	Uuid              uuid.UUID
	Reflected         bool
	Sourcename        string
	Destinationname   string
	Destinationip     string
	Connectiontype    string
	Protocol          string
	Receivedtimestamp int64
}


func (s Sender) RunSender(conn *net.UDPConn) {

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			targets.lock.Lock()
			log.Infof("TARGETS ARRAY:%v", targets)
			for _, t := range targets.targets {
				target := t
				networkOne := &net.UDPAddr{
					IP:   net.ParseIP(target.Networkoneip),
					Port: target.Networkoneport,
				}

				networkTwo := &net.UDPAddr{
					IP:   net.ParseIP(target.Networktwoip),
					Port: target.Networktwoport,
				}

				go func() {
					if target.Networktwoip != ""{
						ssp := s.SendPacket(networkTwo, target.Destinationname, config.HostName, target.Networktwoip, "networkTwo", "udp", conn)
						s.SentPacket <- ssp
					}
				}()
				go func() {
					if target.Networkoneip != ""{
						psp := s.SendPacket(networkOne, target.Destinationname, config.HostName, target.Networkoneip, "networkOne", "udp", conn)
						s.SentPacket <- psp
					}

				}()
			}
			targets.lock.Unlock()
		}
	}
}

func (s Sender) SendPacket(n *net.UDPAddr, destHost string, srcHost string, destIp string, connType string, protocol string, conn *net.UDPConn) *Packet {

	//build the packet
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error("Error creating UUID: ", err)
	}

	//TODO: make packet binary, but text works for now without causing size issues
	p := &Packet{
		Timestamp:         time.Now().UnixNano(),
		Uuid:              id,
		Reflected:         false,
		Sourcename:        srcHost,
		Destinationname:   destHost,
		Destinationip:     destIp,
		Connectiontype:    connType,
		Protocol:          protocol,
		Receivedtimestamp: 0,
	}
	j, err := json.Marshal(p)
	if err != nil {
		log.Error("Error Marshaling packet: ", err)
	}

	//send the packet
	_, err = conn.WriteTo(j, n)
	if err != nil {
		log.Errorf("Sending packet:%s PACKET:%s SRCHOST:%s DESTHOST: %s OVER: %s", err, string(j), srcHost, destHost, connType)
	} else {
		log.Infof("Packet Sent - SOURCE:%s, DEST:%s", conn.LocalAddr().String(), n.String())
	}

	return p

}

func ReadFromSocket(conn *net.UDPConn, pCh chan *Packet, errCh chan error) {
	b := make([]byte, 2048)
	i, _, err := conn.ReadFromUDP(b)
	if err != nil {
		errCh <- err
		return
	}
	receivedtime := time.Now().UnixNano()

	//convert byte packet to string
	res := Packet{}
	err = json.Unmarshal(b[:i], &res)
	if err != nil {
		errCh <- err
		return
	}
	log.Infof("LISTENER READ:%v ON:%s FROM HOST:%s", res, conn.LocalAddr().String(), res.Destinationname)
	res.Receivedtimestamp = receivedtime
	pCh <- &res
}
