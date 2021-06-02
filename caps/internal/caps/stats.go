package caps

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type Stats struct {
	Ctx      context.Context
	Inputch  chan *Packet
	Outputch chan Output
}

type ringBuffer struct {
	index int
	array []int64
}

type Output struct {
	DestName string  `json:"destination_name"`
	DestIp   string  `json:"destination_ip"`
	ConnType string  `json:"connection_type"`
	Protocol string  `json:"protocol"`
	SrcHost  string  `json:"source_name"`
	Latency  []int64 `json:"latencies"`
}

type key struct {
	destinationname string
	connectiontype string
}

func (s Stats) RunStats() {

	m := make(map[key]*ringBuffer)
	for {
		select {
		case myInput := <-s.Inputch:
			rb, ok := m[key{
				destinationname: myInput.Destinationname,
				connectiontype:  myInput.Connectiontype,
			}]
			mapInfo(m, myInput)

			if !ok {
				rb = newRingBuffer()
				m[key{
					destinationname: myInput.Destinationname,
					connectiontype:  myInput.Connectiontype,
				}] = rb
			}
			rb.set(myInput.Receivedtimestamp - myInput.Timestamp)
			if rb.BufferFull() {
				s.Outputch <- Output{
					DestName: myInput.Destinationname,
					DestIp:   myInput.Destinationip,
					ConnType: myInput.Connectiontype,
					Protocol: myInput.Protocol,
					SrcHost:  myInput.Sourcename,
					Latency:  rb.array,
				}
				delete(m, key{
					destinationname: myInput.Destinationname,
					connectiontype:  myInput.Connectiontype,
				})
			}

		}
	}

}


func newRingBuffer() *ringBuffer {
	return &ringBuffer{
		index: 0,
		array: make([]int64, 0, 20),
	}
}

func (b *ringBuffer) BufferFull() bool {
	return len(b.array) >= cap(b.array)
}

func (b *ringBuffer) set(latency int64) {
	if len(b.array) < cap(b.array) {
		b.array = append(b.array, latency)
	}
}

func mapInfo(m map[key]*ringBuffer, p *Packet) {
	for k, v := range m {
		log.Infof("Key: %+v\nValue: %+v\n", k, v)
	}
}
