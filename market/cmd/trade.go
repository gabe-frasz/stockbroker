package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gabe-frasz/stockbroker/market/internal/app/dto"
	"github.com/gabe-frasz/stockbroker/market/internal/app/entity"
	"github.com/gabe-frasz/stockbroker/market/internal/app/transformer"
	"github.com/gabe-frasz/stockbroker/market/internal/infra/messaging"
)

func main() {
	ordersChan := make(chan *entity.Order)
	ordersChanOut := make(chan *entity.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	}
	producer := messaging.NewKafkaProducer(configMap)
	consumer := messaging.NewKafkaConsumer(configMap, []string{"input"})

	go consumer.Consume(kafkaMsgChan)

	book := entity.NewBook(ordersChan, ordersChanOut, wg)
	go book.Trade()

	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			orderInput := &dto.OrderInput{}
			err := json.Unmarshal(msg.Value, &orderInput)
			if err != nil {
				panic(err)
			}
			order := transformer.ToDomainOrder(orderInput)
			ordersChan <- order
		}
	}()

	for res := range ordersChanOut {
		output := transformer.ToDtoOrder(res)
		outputJson, err := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
