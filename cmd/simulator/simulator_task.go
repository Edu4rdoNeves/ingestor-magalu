package simulator

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Edu4rdoNeves/ingestor-magalu/application/service/rabbitmq"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/sirupsen/logrus"
)

type ISimulatorTask interface {
	Run()
}

type SimulatorTask struct {
	RabbitMq rabbitmq.IRabbitMQ
}

func NewSimulatorTask(rabbitmq rabbitmq.IRabbitMQ) ISimulatorTask {
	return &SimulatorTask{
		RabbitMq: rabbitmq,
	}
}

func (s *SimulatorTask) Run() {
	logrus.Info("Pulse Task - Started")

	const (
		totalMessages = 100000
		numWorkers    = 10
		bufferSize    = 1000
	)

	var wg sync.WaitGroup
	msgChan := make(chan dto.PulseData, bufferSize)

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					logrus.Errorf("[Worker %d] panic: %v", workerID, r)
				}
				wg.Done()
			}()
			s.publishMessages(workerID, msgChan)
		}(i)
	}

	// Geração de mensagens
	for i := 0; i < totalMessages; i++ {
		msg := dto.PulseData{
			Tenant:     "magalu",
			ProductSku: fmt.Sprintf("SKU-%d", i%150),
			UseUnity:   fmt.Sprintf("loja-%d", i%10),
			UsedAmount: float64(i%10) + 1,
		}
		msgChan <- msg
	}
	close(msgChan)

	wg.Wait()
	logrus.Info("✅ All messages have been published.")
}

func (s *SimulatorTask) publishMessages(workerID int, messages <-chan dto.PulseData) {
	for msg := range messages {
		body, err := json.Marshal(msg)
		if err != nil {
			logrus.Errorf("[Worker %d] ❌ Failed to marshal message: %v", workerID, err)
			continue
		}

		err = s.RabbitMq.PublishWithNewChannel(body)
		if err != nil {
			logrus.Errorf("[Worker %d] ❌ Failed to publish message: %v", workerID, err)
		} else {
			logrus.Infof("[Worker %d] ✅ Message published: SKU=%s, Unity=%s", workerID, msg.ProductSku, msg.UseUnity)
		}
	}
}

//fazer assincrono

// func (s *SimulatorTask) Run() {
// 	const (
// 		totalMessages = 100000
// 		numWorkers    = 10 // quantidade de workers publicadores
// 		bufferSize    = 1000
// 	)

// 	msgChan := make(chan int, bufferSize)
// 	var wg sync.WaitGroup

// 	for i := 0; i < numWorkers; i++ {
// 		wg.Add(1)
// 		go func(workerID int) {
// 			defer wg.Done()

// 			for i := range msgChan {
// 				defer s.RabbitMq.Close()

// 				rabbitmq.NewRabbitMQ()

// 				msg := dto.PulseData{
// 					Tenant:     "magalu",
// 					ProductSku: fmt.Sprintf("SKU-%d", i%50),
// 					UseUnity:   fmt.Sprintf("loja-%d", i%5),
// 					UsedAmount: float64(i%10) + 1,
// 				}

// 				body, err := json.Marshal(msg)
// 				if err != nil {
// 					log.Printf("[Worker %d] ❌ failed to marshal message %d: %v", workerID, i, err)
// 					continue
// 				}

// 				err = s.RabbitMq.Publish(body)
// 				if err != nil {
// 					log.Printf("[Worker %d] ❌ failed to publish message %d: %v", workerID, i, err)
// 				} else {
// 					log.Printf("[Worker %d] ✅ published message %d", workerID, i)
// 				}

// 			}
// 		}(i)
// 	}

// 	// Feed the workers
// 	for i := 0; i < totalMessages; i++ {
// 		msgChan <- i
// 	}
// 	close(msgChan)

// 	wg.Wait()
// 	log.Println("🚀 All messages published.")
// }
