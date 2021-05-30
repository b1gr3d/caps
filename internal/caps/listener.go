package caps

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
)

type Listener struct {
	Ctx            context.Context
	ReceivedPacket chan *Packet
}


func (l Listener) RunListener(conn *net.UDPConn) {

	for {
		pCh := make(chan *Packet)
		errCh := make(chan error)
		go ReadFromSocket(conn, pCh, errCh)
		select {
		case res := <-pCh:
			l.ReceivedPacket <- res
		case err := <-errCh:
			//TODO: should this count as loss?
			log.Error("Error on err channel: ", err)
		case <-l.Ctx.Done():
		}
	}

}
