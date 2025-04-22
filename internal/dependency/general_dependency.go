package dependency

import (
	"time"

	"github.com/Edu4rdoNeves/ingestor-magalu/application/service/rabbitmq"
	"github.com/Edu4rdoNeves/ingestor-magalu/application/service/redis"
	pulseUsecase "github.com/Edu4rdoNeves/ingestor-magalu/application/usecases/pulse"
	simulatorTask "github.com/Edu4rdoNeves/ingestor-magalu/cmd/simulator"
	pulsetask "github.com/Edu4rdoNeves/ingestor-magalu/cmd/worker/task/pulse_task"
	savePulseTask "github.com/Edu4rdoNeves/ingestor-magalu/cmd/worker/task/save_pulse_task"
	pulseRepo "github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/repository/pulse"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/sirupsen/logrus"
)

var (
	PulseTask     pulsetask.IPulseTask
	SavePulseTask savePulseTask.ISavePulseTask
	SimulatorTask simulatorTask.ISimulatorTask
)

func LoadGeneral() {
	logrus.Info("Loading dependencies...")

	//RABBITMQ
	pulseQueue := rabbitmq.NewDeclareQueue(rabbitmq.QueueDeclare{
		Name:       env.PulseQueueName,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	})

	pulseConsumer := rabbitmq.NewConsumerQueue(rabbitmq.QueueConsumer{
		QueueName: pulseQueue.Name,
		Consumer:  "",
		AutoAck:   true,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}, *pulseQueue)

	pulsePublisher := rabbitmq.NewPublish(rabbitmq.PublishQueue{
		Exchange:  "",
		Key:       pulseConsumer.QueueName,
		Mandatory: false,
		Immediate: false,
	})

	rabbitPulseInstance, err := rabbitmq.NewRabbitMQ(env.PulseQueueUrl, *pulseQueue, *pulseConsumer, *pulsePublisher)
	if err != nil {
		logrus.Fatalf("Error creating RabbitMQ connection: %v", err)
	}

	//REDIS
	redisClient := redis.NewRedisClient(redis.ClientConfig{
		Addr:         env.Addr,
		Password:     env.RedisPassword,
		Db:           env.RedisDb,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	//REPOSITORY
	pulseRepository := pulseRepo.NewPulseRepository(IngestorDb)

	//USECASES
	pulseUsecase := pulseUsecase.NewPulseUseCase(pulseRepository)

	//TASK
	PulseTask = pulsetask.NewPulseTask(redisClient, rabbitPulseInstance)
	if PulseTask == nil {
		logrus.Panic("Pulse Task is nil!")
	}
	logrus.Infof("Pulse Task initialized")

	SavePulseTask = savePulseTask.NewSavePulseTask(redisClient, pulseUsecase)
	if SavePulseTask == nil {
		logrus.Panic("Save Pulse Task is nil!")
	}
	logrus.Infof("Save Pulse Task initialized")

	SimulatorTask = simulatorTask.NewSimulatorTask(rabbitPulseInstance)
	if SimulatorTask == nil {
		logrus.Panic("Simulator Task is nil!")
	}
	logrus.Infof("Simulator Task initialized")
}
