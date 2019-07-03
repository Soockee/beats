package simpleprotocol

import (
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"

	"github.com/elastic/beats/packetbeat/protos"
)

// Transaction Publisher.
type transPub struct {
	sendRequest  bool
	sendResponse bool

	results protos.Reporter
}

func (pub *transPub) onTransaction(requ, resp *message) error {
	if pub.results == nil {
		return nil
	}

	pub.results(pub.createEvent(requ, resp))
	return nil
}

func (pub *transPub) createEvent(requ, resp *message) beat.Event {
	status := common.OK_STATUS

	// resp_time in milliseconds
	responseTime := int32(resp.Ts.Sub(requ.Ts).Nanoseconds() / 1e6)

	src := &common.Endpoint{
		IP:      requ.Tuple.SrcIP.String(),
		Port:    requ.Tuple.SrcPort,
		Process: requ.CmdlineTuple.Src,
	}
	dst := &common.Endpoint{
		IP:      requ.Tuple.DstIP.String(),
		Port:    requ.Tuple.DstPort,
		Process: requ.CmdlineTuple.Dst,
	}

	fields := common.MapStr{
		"type":         "simpleprotocol",
		"status":       status,
		"responsetime": responseTime,
		"bytes_in":     requ.Size,
		"bytes_out":    resp.Size,
		"src":          src,
		"dst":          dst,
	}

	// add processing notes/errors to event
	if len(requ.Notes)+len(resp.Notes) > 0 {
		fields["notes"] = append(requ.Notes, resp.Notes...)
	}

	if pub.sendRequest {
		fields["request"] = requ.content
	}
	if pub.sendResponse {
		fields["response"] = requ.content
	}

	return beat.Event{
		Timestamp: requ.Ts,
		Fields:    fields,
	}
}
