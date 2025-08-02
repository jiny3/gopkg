package pulsarx

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/jiny3/gopkg/configx"
)

// toml struct detail in pulsarx/example.toml
func Default(pulsarPath string) (resource, error) {
	pulsarConfig := struct {
		Client   ClientModel   `toml:"client"`
		Producer ProducerModel `toml:"producer"`
		Consumer ConsumerModel `toml:"consumer"`
	}{
		Client:   ClientModel{},
		Producer: ProducerModel{},
		Consumer: ConsumerModel{},
	}

	err := configx.Read(pulsarPath, &pulsarConfig)
	if err != nil {
		return resource{}, err
	}

	return New(
		&pulsarConfig.Client,
		&pulsarConfig.Producer,
		&pulsarConfig.Consumer,
	)
}

type ClientModel struct {
	URL               string `toml:"url"`
	OperationTimeout  int64  `toml:"operation_timeout"`
	ConnectionTimeout int64  `toml:"connection_timeout"`
}

func (c *ClientModel) Options() pulsar.ClientOptions {
	cc := pulsar.ClientOptions{
		URL: c.URL,
	}
	if c.OperationTimeout > 0 {
		cc.OperationTimeout = time.Duration(c.OperationTimeout) * time.Millisecond
	}
	if c.ConnectionTimeout > 0 {
		cc.ConnectionTimeout = time.Duration(c.ConnectionTimeout) * time.Millisecond
	}
	return cc
}

type ProducerModel struct {
	Topic string `toml:"topic"`
}

func (p *ProducerModel) Options() pulsar.ProducerOptions {
	return pulsar.ProducerOptions{
		Topic: p.Topic,
	}
}

type ConsumerModel struct {
	Topic string `toml:"topic"`
}

func (c *ConsumerModel) Options() pulsar.ConsumerOptions {
	return pulsar.ConsumerOptions{
		Topic: c.Topic,
	}
}
