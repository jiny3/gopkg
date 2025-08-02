package pulsarx

import (
	"github.com/apache/pulsar-client-go/pulsar"
)

type resource struct {
	Client   pulsar.Client
	Producer pulsar.Producer
	Consumer pulsar.Consumer
}

func New(client client, producer producer, consumer consumer) (resource, error) {
	clientOptions := client.Options()
	c, err := pulsar.NewClient(clientOptions)
	if err != nil {
		return resource{}, err
	}

	producerOptions := producer.Options()
	p, err := c.CreateProducer(producerOptions)
	if err != nil {
		return resource{}, err
	}

	consumerOptions := consumer.Options()
	cc, err := c.Subscribe(consumerOptions)
	if err != nil {
		return resource{}, err
	}

	return resource{
		Client:   c,
		Producer: p,
		Consumer: cc,
	}, nil
}
