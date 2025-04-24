package populatequeuetask

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Edu4rdoNeves/ingestor-magalu/application/service/rabbitmq"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/sirupsen/logrus"
)

type IPopulateQueueTask interface {
	Run(populateQueueParams *dto.PopulateQueueParams)
}

type PopulateQueueTask struct {
	RabbitMq rabbitmq.IRabbitMQ
}

func NewPopulateQueueTask(rabbitmq rabbitmq.IRabbitMQ) IPopulateQueueTask {
	return &PopulateQueueTask{
		RabbitMq: rabbitmq,
	}
}

func (p *PopulateQueueTask) Run(populateQueueParams *dto.PopulateQueueParams) {
	logrus.Info("Pulse Task - Started")

	var wg sync.WaitGroup
	msgChan := make(chan dto.PulseData, populateQueueParams.BufferSize)

	for i := 1; i <= populateQueueParams.WorkersNumber; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					logrus.Errorf("[Worker %d] panic: %v", workerID, r)
				}
				wg.Done()
			}()
			p.publishMessages(workerID, msgChan)
		}(i)
	}

	for i := 0; i < populateQueueParams.TotalMessages; i++ {
		msg := dto.PulseData{
			Tenant:     "magalu",
			ProductSku: fmt.Sprintf("SKU-%d", i%env.QntdProductSku),
			UseUnity:   fmt.Sprintf("loja-%d", i%env.QntdUseUnity),
			UsedAmount: float64(i%10) + 1,
		}
		msgChan <- msg
	}
	close(msgChan)

	wg.Wait()
	logrus.Info("✅ All messages have been published.")
}

func (p *PopulateQueueTask) publishMessages(workerID int, messages <-chan dto.PulseData) {
	for msg := range messages {
		body, err := json.Marshal(msg)
		if err != nil {
			logrus.Errorf("[Worker %d] ❌ Failed to marshal message: %v", workerID, err)
			continue
		}

		err = p.RabbitMq.PublishWithNewChannel(body)
		if err != nil {
			logrus.Errorf("[Worker %d] ❌ Failed to publish message: %v", workerID, err)
		} else {
			logrus.Infof("[Worker %d] ✅ Message published: SKU=%s, Unity=%s", workerID, msg.ProductSku, msg.UseUnity)
		}
	}
}
