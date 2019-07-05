package simpleprotocol

import (
	"errors"
	"log"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/streambuf"
	"github.com/elastic/beats/packetbeat/protos/applayer"
	//Debug
)

type parser struct {
	buf     streambuf.Buffer
	config  *parserConfig
	message *message

	onMessage func(m *message) error
}

type parserConfig struct {
	maxBytes int
}

type message struct {
	applayer.Message

	// indicator for parsed message being complete or requires more messages
	// (if false) to be merged to generate full message.
	isComplete bool
	// content stores the number of characters of the message
	content common.NetString
	// indicator for successful parsing
	failed bool
	// list element use by 'transactions' for correlation
	next *message
}

// Error code if stream exceeds max allowed size on append.
var (
	ErrStreamTooLarge = errors.New("Stream data too large")
)

func newParser(config *parserConfig) *parser {
	return &parser{config: config}
}

func (p *parser) init(
	cfg *parserConfig,
	onMessage func(*message) error,
) {
	*p = parser{
		buf:       streambuf.Buffer{},
		config:    cfg,
		onMessage: onMessage,
	}
}

func (p *parser) append(data []byte) error {
	_, err := p.buf.Write(data)
	if err != nil {
		return err
	}

	if p.config.maxBytes > 0 && p.buf.Total() > p.config.maxBytes {
		return ErrStreamTooLarge
	}
	return nil
}

func (p *parser) feed(ts time.Time, data []byte) error {
	if err := p.append(data); err != nil {
		return err
	}

	for p.buf.Total() > 0 {
		if p.message == nil {
			// allocate new message object to be used by parser with current timestamp
			p.message = p.newMessage(ts)
		}

		msg, err := p.parse()
		if err != nil {
			return err
		}
		if msg == nil {
			break // wait for more data
		}

		// reset buffer and message -> handle next message in buffer
		p.buf.Reset()
		p.message = nil

		// call message handler callback
		if err := p.onMessage(msg); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) newMessage(ts time.Time) *message {
	return &message{
		Message: applayer.Message{
			Ts: ts,
		},
	}
}

func (p *parser) parse() (*message, error) {
	// Wait for message being complete
	buf, err := p.buf.CollectUntil([]byte{'\n'})
	if err == streambuf.ErrNoMoreBytes {
		return nil, nil
	}
	msg := p.message
	msg.Size = uint64(p.buf.BufferConsumed())

	isRequest := true
	dir := applayer.NetOriginalDirection

	// check if buffer contains data
	if len(buf) > 0 {
		// get protocol token
		c := buf[0]
		// check if data is a invalid request
		isRequest = !(c == '<') || (c == '!')
		if !isRequest {
			msg.failed = true
			dir = applayer.NetReverseDirection
		}
		buf = buf[1:]
	}
	log.Printf("MSG isRequest: %t\n", isRequest)
	// Set message fields
	msg.content = common.NetString(buf)
	msg.IsRequest = isRequest
	msg.Direction = dir
	return msg, nil
}

/*
func (*parser) parseHTTPLine(m *message) (cont, ok, complete bool) {
	i := bytes.Index(s.data[s.parseOffset:], []byte("\r\n"))
	if i == -1 {
		return false, true, false
	}

	// Very basic tests on the first line. Just to check that
	// we have what looks as an HTTP message
	var version []byte
	var err error
	fline := s.data[s.parseOffset:i]
	if len(fline) < 9 {
		if isDebug {
			debugf("First line too small")
		}
		return false, false, false
	}
	return true, true, true
}*/
