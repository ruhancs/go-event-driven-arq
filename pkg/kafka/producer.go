package kafka

import (
	"encoding/json"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	ConfigMap *ckafka.ConfigMap//conecta com kafka
}

func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{ConfigMap: configMap}
}

//msg == msg que se quer enviar, key=key do kafka que mandar, topic=topico que sera enviado a msg
func (p *Producer) Publish(msg interface{}, key []byte, topic string) error {
	fmt.Println("ENVIANDO NO PUBLISH")
	producer, err := ckafka.NewProducer(p.ConfigMap)
	if err != nil {
		fmt.Println("ENTROU NUM ERRO NO PUBLISH newproducer")
		fmt.Println(err)
		return err
	}
	
	//converte a msg para json
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	//prepara a msg
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          msgJson,
		Key:            key,
	}
	//envia a msg
	err = producer.Produce(message, nil)
	fmt.Println("MSG ENVIADA")
	if err != nil {
		fmt.Println("ERROR ENVIANDO MSG")
		panic(err)
	}
	return nil
}