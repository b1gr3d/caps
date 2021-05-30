package caps

import (
	"context"
	"net"
	"os"
	"sync"

	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
)

var config Config
var configPath = os.Getenv("CONFIG_PATH")

func Run() {
	//config
	var err error
	config, err = SetupConfig(configPath)
	if err != nil {
		logrus.Fatalf("Config Setup Error:%s", err)
	}
	logrus.Infof("CONFIG:%+v", config)

	//create connection for send and listen
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		log.Error("Error dialing UDP: ", err)
	}
	logrus.Infof("Connection Created:%v", conn.LocalAddr())

	// create gochannels for stats processing
	inputChan := make(chan *Packet)
	outputChan := make(chan Output)
	sentpacketChan := make(chan *Packet)
	receivepacketChan := make(chan *Packet)

	// setup targets
	targets = Targets{
		targets: []Target{},
		lock:    &sync.Mutex{},
	}

	//setup register and initial get targets
	var apimonitor = Apimonitor{context.Background(), outputChan}
	Register()
	apimonitor.GetTargets()

	//var setup
	var runreflector = Reflector{context.Background()}
	var stats = Stats{
		Ctx:      context.Background(),
		Inputch:  inputChan,
		Outputch: outputChan,
	}
	var runsender = Sender{
		Ctx:        context.Background(),
		SentPacket: sentpacketChan,
	}

	var runlistener = Listener{
		Ctx:            context.Background(),
		ReceivedPacket: receivepacketChan,
	}

	var runpackettracker = Tracker{
		Ctx:            context.Background(),
		SentPacket:     sentpacketChan,
		ReceivedPacket: receivepacketChan,
		TrackerMap:     make(map[string]*Packet),
		InputChan:      inputChan,
	}

	var runmetric = Metric{
		Ctx:      context.Background(),
		Metricch: outputChan,
	}

	//go routines
	go runpackettracker.PacketTracker()
	go runsender.RunSender(conn)
	go runreflector.RunReflector()
	go runlistener.RunListener(conn)
	go stats.RunStats()
	go apimonitor.RunApimonitor()
	go runmetric.MetricSender()

}
