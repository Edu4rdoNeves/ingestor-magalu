package pulsetask

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Edu4rdoNeves/ingestor-magalu/application/service/rabbitmq"
	"github.com/Edu4rdoNeves/ingestor-magalu/application/service/redis"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/Edu4rdoNeves/ingestor-magalu/utils"
	"github.com/sirupsen/logrus"
)

type IPulseTask interface {
	Run(ctx context.Context)
}

type PulseTask struct {
	Redis    redis.IRedisClient
	RabbitMq rabbitmq.IRabbitMQ
}

func NewPulseTask(redis redis.IRedisClient, rabbitmq rabbitmq.IRabbitMQ) IPulseTask {
	return &PulseTask{
		Redis:    redis,
		RabbitMq: rabbitmq,
	}
}

func (t *PulseTask) Run(ctx context.Context) {
	logrus.Info("Pulse Task - Started")

	numWorkers := env.PulseWorkersNumber
	var wg sync.WaitGroup
	messages := make(chan []byte, env.PulseMessageBuffer)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					logrus.Errorf("Worker %d panicked: %v", workerID, r)
				}
				wg.Done()
			}()
			t.process(ctx, workerID, messages)
		}(i)
	}

	err := t.RabbitMq.Consumer(
		func(body []byte) {
			logrus.Infof("Message received: %s", string(body))
			messages <- body
		},
	)
	if err != nil {
		logrus.Errorf("Pulse Task - Failed to consume message: %v", err)
		close(messages)
		wg.Wait()
		return
	}

	<-ctx.Done()

	logrus.Info("Pulse Task - Context canceled, shutting down...")
	close(messages)
	wg.Wait()
}

func (t *PulseTask) process(ctx context.Context, id int, messages <-chan []byte) {
	logger := logrus.WithField("worker_id", id)

	for {
		select {
		case <-ctx.Done():
			logger.Infof("Closing due to context canceled")
			return
		case msg, ok := <-messages:
			if !ok {
				logger.Infof("Closing: channel closed")
				return
			}

			var pulseData *dto.PulseData
			if err := json.Unmarshal(msg, &pulseData); err != nil {
				logger.Errorf("JSON unmarshal error: %v", err)
				continue
			}

			key := fmt.Sprintf("pulse:%s:%s:%s", pulseData.Tenant, pulseData.ProductSku, pulseData.UseUnity)

			err := utils.Retry(env.RedisMaxRetry, time.Duration(env.RedisTimeToSleep)*time.Millisecond, func() error {
				return t.Redis.IncrementCounter(key, pulseData.UsedAmount)
			})
			if err != nil {
				logger.WithFields(logrus.Fields{
					"worker_id": id,
					"key":       key,
				}).Errorf("Failed to update Redis after retries: %v", err)
				continue
			}

			logger.Infof("Processed pulse for key: pulse:%s:%s:%s - Incremented by: %.2f",
				pulseData.Tenant, pulseData.ProductSku, pulseData.UseUnity, pulseData.UsedAmount)
		}
	}
}
