package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

type Consumer struct {
	ConfigMap *ckafka.ConfigMap
	Topics    []string
}

//topicos é onde o consumidor se inscrevera para receber os resultados
func NewConsumer(configMap *ckafka.ConfigMap, topics []string) *Consumer {
	return &Consumer{
		ConfigMap: configMap,
		Topics:    topics,
	}
}

func (c *Consumer) Consume(msgChan chan *ckafka.Message) error {
	consumer, err := ckafka.NewConsumer(c.ConfigMap)
	if err != nil {
		panic(err)
	}
	err = consumer.SubscribeTopics(c.Topics, nil)
	if err != nil {
		panic(err)
	}
	//todos resultados recebidos sao enviados para os inscritos no topico
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			//envia para um canal de thread para ler uma thread consome os dados do kafka e outra sobe o webserver
			msgChan <- msg 
		}
	}
}