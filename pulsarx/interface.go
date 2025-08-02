package pulsarx

import "github.com/apache/pulsar-client-go/pulsar"

type client interface {
	Options() pulsar.ClientOptions
}

type producer interface {
	Options() pulsar.ProducerOptions
}

type consumer interface {
	Options() pulsar.ConsumerOptions
}
