package simpleprotocol

import (
	"github.com/elastic/beats/packetbeat/config"
	"github.com/elastic/beats/packetbeat/protos"
)

type simpleprotocolConfig struct {
	config.ProtocolCommon `config:",inline"`
}

var (
	defaultConfig = simpleprotocolConfig{
		ProtocolCommon: config.ProtocolCommon{
			TransactionTimeout: protos.DefaultTransactionExpiration,
		},
	}
)

func (c *simpleprotocolConfig) Validate() error {
	return nil
}
