package comm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	
	zmq "github.com/pebbe/zmq4"
	
	"noname001/dilemma/comm/zmqdep/kvmsg"
)

// temporary implementation of clonesrv5
// TODO: CHECK FOR ALL that can return ERRROR

type clonesrv_t struct {
	kvmap     map[string]*kvmsg.Kvmsg
	sequence  int64 // TODO: THIS CANNOT RUN FOREVER WITHOUT RESETTING SEQUENCE
	snapshot  *zmq.Socket
	publisher *zmq.Socket
	collector *zmq.Socket
}

type PubSubServerConfig struct {
	Verbose       bool

	SnapshotHost  string
	PublisherHost string
	CollectorHost string

}
type PubSubServer struct {
	cfg        PubSubServerConfig
	cloneSrv5  *clonesrv_t
	zmqReactor *zmq.Reactor
}

func NewPubsubServer(cfg PubSubServerConfig) (*PubSubServer, error) {
	var err error
	
	srv := &clonesrv_t {
		kvmap: make(map[string]*kvmsg.Kvmsg),
	}

	snapshotAddr := fmt.Sprintf("tcp://%s", cfg.SnapshotHost)
	publisherAddr := fmt.Sprintf("tcp://%s", cfg.PublisherHost)
	collectorAddr := fmt.Sprintf("tcp://%s", cfg.CollectorHost)

	srv.snapshot, err = zmq.NewSocket(zmq.ROUTER)
	if err != nil { return nil, err }

	err = srv.snapshot.Bind(snapshotAddr)
	if err != nil { return nil, err }

	srv.publisher, err = zmq.NewSocket(zmq.PUB)
	if err != nil { return nil, err }

	err = srv.publisher.Bind(publisherAddr)
	if err != nil { return nil, err }

	srv.collector, err = zmq.NewSocket(zmq.PULL)
	if err != nil { return nil, err }

	err = srv.collector.Bind(collectorAddr)
	if err != nil { return nil, err }

	reactor := zmq.NewReactor()
	reactor.SetVerbose(cfg.Verbose) // the only way to see if the reactor is indeed runnin ?
	reactor.AddSocket(srv.snapshot, zmq.POLLIN,
		func(e zmq.State) error { return snapshots(srv) })
	reactor.AddSocket(srv.collector, zmq.POLLIN,
		func(e zmq.State) error { return collector(srv) })
	reactor.AddChannelTime(time.Tick(1000*time.Millisecond), 1,
		func(v any) error { return flush_ttl(srv) })

	logger.Debugf("COMM::PSS: snapshot  is active at %s", snapshotAddr)
	logger.Debugf("COMM::PSS: publisher is active at %s", publisherAddr)
	logger.Debugf("COMM::PSS: collector is active at %s", collectorAddr)

	pubsubServer := &PubSubServer{
		cfg: cfg,
		cloneSrv5: srv,
		zmqReactor: reactor,
	}

	return pubsubServer, nil
}
// TODO: err handling
func (pss *PubSubServer) Start() {
	logger.Debugf("COMM::PSS: pubsubserver started...")
	go func() {
		err := pss.zmqReactor.Run(100 * time.Millisecond) // ??? precision: .1 secs 
		if err != nil {
			logger.Errorf("COMM::PSS: pubsubserver start err %s", err)
			return
		}
	}()
}

func (pss *PubSubServer) Stop() (err error) {
	// TODO: remove socket form reactor, close all socket

	logger.Debugf("COMM::PSS: pubsubserver stopped...")
	return
}

//  This is the reactor handler for the snapshot socket; it accepts
//  just the ICANHAZ? request and replies with a state snapshot ending
//  with a KTHXBAI message:
func snapshots(srv *clonesrv_t) (err error) {
	msg, err := srv.snapshot.RecvMessage(0)
	if err != nil {
		// TODO
		return
	}
	identity := msg[0]

	request := msg[1]
	if request != "ICANHAZ?" {
		// TODO
		err = errors.New("E: bad request, aborting")
		return
	}
	subtree := msg[2]

	for _, kvmsg := range srv.kvmap {
		if key, _ := kvmsg.GetKey(); strings.HasPrefix(key, subtree) {
			srv.snapshot.Send(identity, zmq.SNDMORE)
			kvmsg.Send(srv.snapshot)
		}
	}

	// now send END message with sequence number
	logger.Debug("I: sending snapshot = ", srv.sequence)
	srv.snapshot.Send(identity, zmq.SNDMORE)
	kvmsg := kvmsg.NewKvmsg(srv.sequence)
	kvmsg.SetKey("KTHXBAI")
	kvmsg.SetBody(subtree)
	kvmsg.Send(srv.snapshot)

	return
}

//  We store each update with a new sequence number, and if necessary, a
//  time-to-live. We publish updates immediately on our publisher socket:
func collector(srv *clonesrv_t) (err error) {
	kvmsg, err := kvmsg.RecvKvmsg(srv.collector)
	if err != nil {
		// TODO:
		return
	}

	srv.sequence++
	kvmsg.SetSequence(srv.sequence)
	kvmsg.Send(srv.publisher)
	if ttls, e := kvmsg.GetProp("ttl"); e == nil {
		// change duration into specific time, using the same property: ugly!
		ttl, e := strconv.ParseInt(ttls, 10, 64)
		if e != nil {
			// TODO
			err = e
			return
		}
		kvmsg.SetProp("ttl", fmt.Sprint(time.Now().Add(time.Duration(ttl)*time.Second).Unix()))
	}
	kvmsg.Store(srv.kvmap)
	logger.Debug("I: publishing update = ", srv.sequence)

	return
}

//  At regular intervals we flush ephemeral values that have expired. This
//  could be slow on very large data sets:
func flush_ttl(srv *clonesrv_t) (err error) {
	for _, kvmsg := range srv.kvmap {
		//  If key-value pair has expired, delete it and publish the
		//  fact to listening clients.
		if ttls, e := kvmsg.GetProp("ttl"); e == nil {
			ttl, e := strconv.ParseInt(ttls, 10, 64)
			if e != nil {
				// TODO:
				err = e
				continue
			}
			if time.Now().After(time.Unix(ttl, 0)) {
				srv.sequence++
				kvmsg.SetSequence(srv.sequence)
				kvmsg.SetBody("")
				e = kvmsg.Send(srv.publisher)
				if e != nil {
					err = e
				}
				kvmsg.Store(srv.kvmap)
				logger.Debug("I: publishing delete = ", srv.sequence)
			}
		}
	}

	return
}


// === ### === ### === CLIENT === ### === ### ===
// const (
// 	SUBTREE = "/client/"
// )

type PublisherClientConfig struct {
	Context context.Context

	CollectorServerHost string // hostname:port
}
type PublisherClient struct {
	context context.Context
	cancel  context.CancelFunc

	cfg     PublisherClientConfig
	pubSock *zmq.Socket

	dataChannel chan *kvmsg.Kvmsg
}

func NewPublisherClient(cfg PublisherClientConfig) (pc *PublisherClient, err error) {
	pc = &PublisherClient{cfg: cfg}
	pc.context, pc.cancel = context.WithCancel(cfg.Context)

	pc.pubSock, err = zmq.NewSocket(zmq.PUSH)
	if err != nil {
		logger.Errorf("COMM::PUB: pub new err %s", err.Error())
		return
	}

	pc.dataChannel = make(chan *kvmsg.Kvmsg)

	return
}
func (pc *PublisherClient) Connect() (err error) {
	err = pc.pubSock.Connect(fmt.Sprintf("tcp://%s", pc.cfg.CollectorServerHost))
	if err != nil {
		logger.Errorf("COMM::PUB: pub conn err %s", err.Error())
		return
	}

	go func() {
		SendLoop:
		for {
			select {
			case <- pc.context.Done():
				break SendLoop
			case _kvmsg, _ := <- pc.dataChannel:
				err = _kvmsg.Send(pc.pubSock)
				if err != nil {
					logger.Errorf("publisher send err ? %s", err)
					continue
				}
			}
		}
	}()

	return
}
func (pc *PublisherClient) Disconnect() (err error) {

	return
}

// https://pkg.go.dev/github.com/pebbe/zmq4#Context.NewSocket
// WARNING: The Socket is not thread safe. This means that you cannot access the same Socket from different goroutines without using something like a mutex. 
func (pc *PublisherClient) Publish(kvmsg *kvmsg.Kvmsg) (err error) {
	pc.dataChannel <- kvmsg
	return
}

type SubscriberClientConfig struct {
	Context             context.Context

	SnapshotServerHost  string
	PublisherServerHost string
	Subtree             string

	DataChannel         chan string
}
type SubscriberClient struct {
	ctx context.Context
	ccl context.CancelFunc
	
	cfg SubscriberClientConfig

	snapSock *zmq.Socket
	subSock  *zmq.Socket
}

func NewSubscriber(cfg SubscriberClientConfig) (sc *SubscriberClient, err error) {
	sc = &SubscriberClient{cfg: cfg}
	sc.ctx, sc.ccl = context.WithCancel(cfg.Context)

	sc.snapSock, err = zmq.NewSocket(zmq.DEALER)
	if err != nil {
		logger.Errorf("COMM::SUB: snap new sock err %s", err.Error())
		return
	}

	sc.subSock, err = zmq.NewSocket(zmq.SUB)
	if err != nil {
		logger.Errorf("COMM::SUB: sub new sock err %s", err.Error())
		return
	}

	return
}

func (sc *SubscriberClient) Connect() (err error) {
	err = sc.snapSock.Connect(fmt.Sprintf("tcp://%s", sc.cfg.SnapshotServerHost))
	if err != nil {
		logger.Errorf("COMM::SUB: snap sock conn err %s", err.Error())
		return
	}

	err = sc.subSock.Connect(fmt.Sprintf("tcp://%s", sc.cfg.PublisherServerHost))
	if err != nil {
		logger.Errorf("COMM::SUB: sub sock conn err %s", err.Error())
		return
	}
	err = sc.subSock.SetSubscribe(sc.cfg.Subtree)
	if err != nil {
		logger.Errorf("COMM::SUB: sub sock subtree err %s", err.Error())
		return
	}

	go func() {
		kvmap := make(map[string]*kvmsg.Kvmsg)

		// request snapshot state
		lastSeq := int64(0)
		sc.snapSock.SendMessage("ICANHAZ?", sc.cfg.Subtree)
		for {
			kvmsg, err := kvmsg.RecvKvmsg(sc.snapSock)
			if err != nil {
				// TODO: channel
				logger.Errorf("COMM::SUB: snapshot state err %s", err.Error())
				break
			}
			if key, _ := kvmsg.GetKey(); key == "KTHXBAI" {
				lastSeq, _ := kvmsg.GetSequence()
				logger.Debugf("COMM::SUB: snapshot received = %d", lastSeq)
				break // done
			}
			kvmsg.Store(kvmap)
		}
		sc.snapSock.Close()
		sc.snapSock = nil


		// start receiving msg
		poller := zmq.NewPoller()
		poller.Add(sc.subSock, zmq.POLLIN)
		RecvLoop:
		for {
			polled, err := poller.Poll(1000 * time.Millisecond)
			if err != nil {
				// TODO: channel
				logger.Errorf("COMM::SUB: poll err %s", err.Error())
				time.Sleep(1*time.Second)
				continue
			}
			if len(polled) > 0 {
				kvmsg, err := kvmsg.RecvKvmsg(sc.subSock)
				if err != nil {
					// TODO: channel
					logger.Errorf("COMM::SUB: recv err %s", err.Error())
					time.Sleep(1*time.Second)
					continue
				}

				if seq, _ := kvmsg.GetSequence(); seq > lastSeq {
					// TODO: channel to pass message
					lastSeq = seq
					kvmsg.Store(kvmap)
					logger.Debugf("COMM::SUB: update received = %d", lastSeq)
				}

				if sc.cfg.DataChannel != nil {
					bod, err := kvmsg.GetBody()
					if err != nil {
						// TODO
						continue
					}
					
					sc.cfg.DataChannel <- bod
				}
			}

			select {
			case <- sc.ctx.Done(): break RecvLoop
			default: // let it pass
			}
		}
	}()

	return
}
func (sc *SubscriberClient) Disconnect() (err error) {
	return
}
