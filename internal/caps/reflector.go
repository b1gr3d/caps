package caps

import (
	"context"
	"encoding/json"
	"net"

	"github.com/sirupsen/logrus"
)

type Reflector struct {
	Ctx context.Context
}

type Result struct {
	rtt        int64
	targetName string
}

func (r Reflector) RunReflector() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   nil, //listen on all available IPs
		Port: config.ReflectorUdpPort,
	})
	if err != nil {
		logrus.Errorf("Reflector ListenUDP %s ", err)
	} else {
		logrus.Infof("REFLECTOR LISTENING ON:%v", conn.LocalAddr())
	}
	defer conn.Close()

	// byte array for packet
	b := make([]byte, 2048)
	for {

		//read the packet that was received
		i, n, err := conn.ReadFromUDP(b)
		if err != nil {
			logrus.Errorf("REFLECTOR READUDP %s", err)
			continue
		}

		//convert byte packet to string
		res := Packet{}
		err = json.Unmarshal(b[:i], &res)
		if err != nil {
			logrus.Errorf("Unmarshal on Reflector:%s %s", err, string(b[:i]))
			continue
		}
		logrus.Infof("REFLECTOR READ:%v", res)

		if res.Reflected == false {
			res.Reflected = true
			rp, err := json.Marshal(res)
			if err != nil {
				logrus.Errorf("Marshaling packet in Reflector: %s ", err)
			}
			conn.WriteTo(rp, n)
			logrus.Infof("REFLECTOR SENT:%v TO:%s FROM HOST:%s TO HOST:%s", res, n, res.Destinationname, res.Sourcename)
			continue
		}
	}
}
