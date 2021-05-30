package caps

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

type Tracker struct {
	Ctx            context.Context
	SentPacket     chan *Packet
	ReceivedPacket chan *Packet
	TrackerMap     map[string]*Packet
	InputChan      chan *Packet
}

func (p *Tracker) PacketTracker() {

	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case sp := <-p.SentPacket:
			p.TrackerMap[sp.Uuid.String()] = sp
		case rp := <-p.ReceivedPacket:
			delete(p.TrackerMap, rp.Uuid.String())
			p.InputChan <- rp
		case <-ticker.C:
			for u, v := range p.TrackerMap {
				// loss is considered >3 sec since sent timestamp on packet
				if time.Since(time.Unix(0, v.Timestamp)) > time.Second*3 {
					log.Warnf("Packet has been lost for UUID: %v", v)
					message := v
					message.Receivedtimestamp = 0
					p.InputChan <- message
					delete(p.TrackerMap, u)
				}

			}
		}
	}
}
