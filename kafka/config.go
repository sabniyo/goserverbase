package kafka

type KafkaCred struct {
	Brokers       interface{} `json:"bootstrap.servers,omitempty"`
	Username      interface{} `json:"sasl.username,omitempty"`
	Password      interface{} `json:"sasl.password,omitempty"`
	SASLMechanism interface{} `json:"sasl.mechanisms,omitempty"`
	SASLProtocol  interface{} `json:"security.protocol,omitempty"`
	ClientID      interface{} `json:"client.id,omitempty"`
}

type KafkaConsumerConfig struct {
	*KafkaCred
	GroupID        interface{} `json:"group.id,omitempty"`
	GoEventChannel bool        `json:"go.events.channel.enable,omitempty"`
	OffsetReset    interface{} `json:"auto.offset.reset,omitempty"`
}

type KafkaProducerConfig struct {
	*KafkaCred
	Acknowledge interface{} `json:"acks,omitempty"`
}
