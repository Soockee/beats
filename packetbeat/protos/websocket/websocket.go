package websocket

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/elastic/beats/packetbeat/protos"
	"github.com/elastic/beats/packetbeat/protos/tcp"
)

// websocketPlugin application level protocol analyzer plugin
type websocketPlugin struct {
	ports        protos.PortsConfig
	parserConfig parserConfig
	transConfig  transactionConfig
	pub          transPub
}

// Application Layer tcp stream data to be stored on tcp connection context.
type connection struct {
	streams [2]*stream
	trans   transactions
}

// Uni-directional tcp stream state for parsing messages.
type stream struct {
	parser parser
}

var (
	debugf = logp.MakeDebug("websocket")

	// use isDebug/isDetailed to guard debugf/detailedf to minimize allocations
	// (garbage collection) when debug log is disabled.
	isDebug = false
)

func init() {
	protos.Register("websocket", New)
}

// New create and initializes a new websocket protocol analyzer instance.
func New(
	testMode bool,
	results protos.Reporter,
	cfg *common.Config,
) (protos.Plugin, error) {
	p := &websocketPlugin{}
	config := defaultConfig
	if !testMode {
		if err := cfg.Unpack(&config); err != nil {
			return nil, err
		}
	}

	if err := p.init(results, &config); err != nil {
		return nil, err
	}
	return p, nil
}

func (wp *websocketPlugin) init(results protos.Reporter, config *websocketConfig) error {
	if err := wp.setFromConfig(config); err != nil {
		return err
	}
	wp.pub.results = results

	isDebug = logp.IsDebug("http")
	return nil
}

func (wp *websocketPlugin) setFromConfig(config *websocketConfig) error {

	// set module configuration
	if err := wp.ports.Set(config.Ports); err != nil {
		return err
	}

	// set parser configuration
	parser := &wp.parserConfig
	parser.maxBytes = tcp.TCPMaxDataInStream

	// set transaction correlator configuration
	trans := &wp.transConfig
	trans.transactionTimeout = config.TransactionTimeout

	// set transaction publisher configuration
	pub := &wp.pub
	pub.sendRequest = config.SendRequest
	pub.sendResponse = config.SendResponse

	return nil
}

// ConnectionTimeout returns the per stream connection timeout.
// Return <=0 to set default tcp module transaction timeout.
func (wp *websocketPlugin) ConnectionTimeout() time.Duration {
	return wp.transConfig.transactionTimeout
}

// GetPorts returns the ports numbers packets shall be processed for.
func (wp *websocketPlugin) GetPorts() []int {
	return wp.ports.Ports
}

// Parse processes a TCP packet. Return nil if connection
// state shall be dropped (e.g. parser not in sync with tcp stream)
func (wp *websocketPlugin) Parse(
	pkt *protos.Packet,
	tcptuple *common.TCPTuple, dir uint8,
	private protos.ProtocolData,
) protos.ProtocolData {
	defer logp.Recover("Parse websocketPlugin exception")

	conn := wp.ensureConnection(private)
	st := conn.streams[dir]
	if st == nil {
		st = &stream{}
		st.parser.init(&wp.parserConfig, func(msg *message) error {
			return conn.trans.onMessage(tcptuple.IPPort(), dir, msg)
		})
		conn.streams[dir] = st
	}

	if err := st.parser.feed(pkt.Ts, pkt.Payload); err != nil {
		debugf("%v, dropping TCP stream for error in direction %v.", err, dir)
		wp.onDropConnection(conn)
		return nil
	}
	return conn
}

// ReceivedFin handles TCP-FIN packet.
func (wp *websocketPlugin) ReceivedFin(
	tcptuple *common.TCPTuple, dir uint8,
	private protos.ProtocolData,
) protos.ProtocolData {
	return private
}

// GapInStream handles lost packets in tcp-stream.
func (wp *websocketPlugin) GapInStream(tcptuple *common.TCPTuple, dir uint8,
	nbytes int,
	private protos.ProtocolData,
) (protos.ProtocolData, bool) {
	conn := getConnection(private)
	if conn != nil {
		wp.onDropConnection(conn)
	}

	return nil, true
}

// onDropConnection processes and optionally sends incomplete
// transaction in case of connection being dropped due to error
func (wp *websocketPlugin) onDropConnection(conn *connection) {
}

func (wp *websocketPlugin) ensureConnection(private protos.ProtocolData) *connection {
	conn := getConnection(private)
	if conn == nil {
		conn = &connection{}
		conn.trans.init(&wp.transConfig, wp.pub.onTransaction)
	}
	return conn
}

func (conn *connection) dropStreams() {
	conn.streams[0] = nil
	conn.streams[1] = nil
}

func getConnection(private protos.ProtocolData) *connection {
	if private == nil {
		return nil
	}

	priv, ok := private.(*connection)
	if !ok {
		logp.Warn("websocket connection type error")
		return nil
	}
	if priv == nil {
		logp.Warn("Unexpected: websocket connection data not set")
		return nil
	}
	return priv
}
