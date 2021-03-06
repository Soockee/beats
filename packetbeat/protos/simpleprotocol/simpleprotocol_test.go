// +build !integration

package simpleprotocol

import (
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/stretchr/testify/assert"
)

type testParser struct {
	payloads   []string
	simpleprot *simpleprotocolPlugin
	stream     *stream
}

var testParserConfig = parserConfig{}

func newTestParser(simpleprot *simpleprotocolPlugin, payloads ...string) *testParser {
	if simpleprot == nil {
		simpleprot = simpleModForTests()
	}
	tp := &testParser{
		simpleprot: simpleprot,
		payloads:   payloads,
		stream:     &stream{},
	}
	return tp
}
func simpleModForTests() *simpleprotocolPlugin {
	callback := func(beat.Event) {}
	simpleprot, err := New(false, callback, common.NewConfig())
	if err != nil {
		panic(err)
	}
	return simpleprot.(*simpleprotocolPlugin)
}
func (tp *testParser) parse() bool {
	st := tp.stream
	if len(tp.payloads) > 0 {
		// append data from payload to stream
		// currently only one data-pack
		st.parser.feed(time.Now(), []byte(tp.payloads[0]))
	}

	st.parser = *newParser(&tp.simpleprot.parserConfig)
	st.parser.parse()
	return true
}

func TestSimpleprotocolParser_ReceivingData(t *testing.T) {
	/* data := "HTTP/1.1 200 OK\r\n" +
	"Date: Tue, 14 Aug 2012 22:31:45 GMT\r\n" +
	"Expires: -1\r\n" +
	"Cache-Control: private, max-age=0\r\n" +
	"Content-Type: text/html; charset=UTF-8\r\n" +
	"Content-Encoding: gzip\r\n" +
	"Server: gws\r\n" +
	"Content-Length: 0\r\n" +
	"X-XSS-Protection: 1; mode=block\r\n" +
	"X-Frame-Options: SAMEORIGIN\r\n" +
	"\r\n"
	*/

	//message, ok, complete := &message{}, true, true //testParse(nil, data)

	assert.True(t, true)
	assert.True(t, true)
	//assert.False(t, message.isRequest)
	//assert.Equal(t, 200, int(message.statusCode))
	//assert.Equal(t, "OK", string(message.statusPhrase))
	//assert.True(t, isVersion(message.version, 1, 1))
	//assert.Equal(t, 262, int(message.size))
	//assert.Equal(t, 0, message.contentLength)
}

/*
func TestHttpParser_simpleResponse(t *testing.T) {
	data := "HTTP/1.1 200 OK\r\n" +
		"Date: Tue, 14 Aug 2012 22:31:45 GMT\r\n" +
		"Expires: -1\r\n" +
		"Cache-Control: private, max-age=0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Content-Encoding: gzip\r\n" +
		"Server: gws\r\n" +
		"Content-Length: 0\r\n" +
		"X-XSS-Protection: 1; mode=block\r\n" +
		"X-Frame-Options: SAMEORIGIN\r\n" +
		"\r\n"
	message, ok, complete := testParse(nil, data)

	assert.True(t, ok)
	assert.True(t, complete)
	assert.False(t, message.isRequest)
	assert.Equal(t, 200, int(message.statusCode))
	assert.Equal(t, "OK", string(message.statusPhrase))
	assert.True(t, isVersion(message.version, 1, 1))
	assert.Equal(t, 262, int(message.size))
	assert.Equal(t, 0, message.contentLength)
}
**/

/**
func TestHttpParser_simpleRequest(t *testing.T) {
	http := httpModForTests(nil)
	http.parserConfig.sendHeaders = true
	http.parserConfig.sendAllHeaders = true

	data := "GET / HTTP/1.1\r\n" +
		"Host: www.google.ro\r\n" +
		"Connection: keep-alive\r\n" +
		"User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_4) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.75 Safari/537.1\r\n" +
		"Accept: */ //*\r\n" +
/*		"X-Chrome-Variations: CLa1yQEIj7bJAQiftskBCKS2yQEIp7bJAQiptskBCLSDygE=\r\n" +
		"Referer: http://www.google.ro/\r\n" +
		"Accept-Encoding: gzip,deflate,sdch\r\n" +
		"Accept-Language: en-US,en;q=0.8\r\n" +
		"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.3\r\n" +
		"Cookie: PREF=ID=6b67d166417efec4:U=69097d4080ae0e15:FF=0:TM=1340891937:LM=1340891938:S=8t97UBiUwKbESvVX; NID=61=sf10OV-t02wu5PXrc09AhGagFrhSAB2C_98ZaI53-uH4jGiVG_yz9WmE3vjEBcmJyWUogB1ZF5puyDIIiB-UIdLd4OEgPR3x1LHNyuGmEDaNbQ_XaxWQqqQ59mX1qgLQ\r\n" +
		"\r\n" +
		"garbage"

	message, ok, complete := testParse(http, data)

	assert.True(t, ok)
	assert.True(t, complete)
	assert.True(t, message.isRequest)
	assert.True(t, isVersion(message.version, 1, 1))
	assert.Equal(t, 669, int(message.size))
	assert.Equal(t, "GET", string(message.method))
	assert.Equal(t, "/", string(message.requestURI))
	assert.Equal(t, "www.google.ro", string(message.headers["host"]))
}
*/
