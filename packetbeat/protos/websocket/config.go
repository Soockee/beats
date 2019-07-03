package websocket

import (
	"github.com/elastic/beats/packetbeat/config"
	"github.com/elastic/beats/packetbeat/protos"
)

type websocketConfig struct {
	config.ProtocolCommon `config:",inline"`
}

var (
	defaultConfig = websocketConfig{
		ProtocolCommon: config.ProtocolCommon{
			TransactionTimeout: protos.DefaultTransactionExpiration,
		},
	}
)

func (c *websocketConfig) Validate() error {
	return nil
}
