package savepulsetask

import (
	"context"
	"sync"
	"time"

	"github.com/Edu4rdoNeves/ingestor-magalu/application/service/redis"
	pulseUsecase "github.com/Edu4rdoNeves/ingestor-magalu/application/usecases/pulse"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/Edu4rdoNeves/ingestor-magalu/utils"
	"github.com/sirupsen/logrus"
)

type ISavePulseTask interface {
	Run(ctx context.Context)
}

type SavePulseTask struct {
	Redis        redis.IRedisClient
	PulseUseCase pulseUsecase.IPulseUseCase
}

func NewSavePulseTask(redis redis.IRedisClient, pulseUC pulseUsecase.IPulseUseCase) ISavePulseTask {
	return &SavePulseTask{
		Redis:        redis,
		PulseUseCase: pulseUC,
	}
}

func (t *SavePulseTask) Run(ctx context.Context) {
	logrus.Info("Save Pulse Task - Started")

	var wg sync.WaitGroup
	keysChan := make(chan string, env.SavePulseMessageBuffer)

	for i := 1; i <= env.SavePulseWorkersNumber; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			t.process(ctx, workerID, keysChan, env.SavePulseBatch)
		}(i)
	}

	keys, err := t.Redis.GetKeysByPattern("pulse:*")
	if err != nil {
		logrus.Errorf("Failed to get keys: %v", err)
		close(keysChan)
		wg.Wait()
		return
	}

	for _, key := range keys {
		keysChan <- key
	}

	close(keysChan)
	wg.Wait()

	logrus.Info("Save Pulse Task - Finished")
}

func (t *SavePulseTask) process(ctx context.Context, id int, keysChan <-chan string, batchSize int) {
	log := logrus.WithFields(logrus.Fields{
		"task":   "SavePulseTask",
		"worker": id,
	})

	log.Info("Worker started")

	var (
		pulseBatch []*dto.PulseData
		keyMap     = make(map[*dto.PulseData]string)
	)

	for {
		select {
		case <-ctx.Done():
			log.Info("Context canceled, flushing batch before exit...")
			t.saveAndDelete(ctx, pulseBatch, keyMap, log)
			pulseBatch = nil
			return

		case key, ok := <-keysChan:
			if !ok {
				log.Info("Channel closed, flushing remaining data...")
				t.saveAndDelete(ctx, pulseBatch, keyMap, log)
				pulseBatch = nil

				return
			}

			pulseData, err := t.buildPulseDataFromKey(key)
			if err != nil {
				log.Errorf("Failed to parse key %s: %v", key, err)
				continue
			}

			pulseBatch = append(pulseBatch, pulseData)
			keyMap[pulseData] = key

			if len(pulseBatch) >= batchSize {
				t.saveAndDelete(ctx, pulseBatch, keyMap, log)
				pulseBatch = nil
				keyMap = make(map[*dto.PulseData]string)
			}
		}
	}
}

func (t *SavePulseTask) saveAndDelete(ctx context.Context, batch []*dto.PulseData, keyMap map[*dto.PulseData]string, log *logrus.Entry) {
	if len(batch) == 0 {
		logrus.Debug("No data to flush")
		return
	}

	err := t.PulseUseCase.SavePulseBatch(ctx, batch)
	if err != nil {
		log.Errorf("Failed to save batch: %v", err)
		return
	}

	for _, pulse := range batch {
		key := keyMap[pulse]
		log.Infof("Deleting key: %s", key)
		err := utils.Retry(env.RedisMaxRetry, time.Duration(env.RedisTimeToSleep)*time.Millisecond, func() error {
			return t.Redis.DeleteKey(key)
		})

		if err != nil {
			logrus.Errorf("Failed to delete key %s: %v", key, err)
		}
	}

	logrus.Infof("Batch saved and deleted successfully. Total: %d", len(batch))
	batch = nil
}

func (t *SavePulseTask) buildPulseDataFromKey(key string) (*dto.PulseData, error) {
	count, err := t.Redis.GetValue(key)
	if err != nil {
		logrus.Errorf("Failed to get value for key %s: %v", key, err)
		return nil, err
	}

	convertedCount, err := utils.StringToFloat64(count)
	if err != nil {
		logrus.Errorf("error when try converting string to float64. Error: %v", err)
		return nil, err
	}

	pulseData, err := utils.ParsePulseKey(key, convertedCount)
	if err != nil {
		logrus.Errorf("Failed to parse key %s: %v", key, err)
		return nil, err
	}

	return pulseData, nil
}
