package env

import (
	"os"

	"github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/database"
	"github.com/Edu4rdoNeves/ingestor-magalu/utils"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Queue
var (
	PulseQueueUrl  string
	PulseQueueName string
)

// REDIS
var (
	Addr             string
	RedisPassword    string
	RedisDb          int
	RedisMaxRetry    int
	RedisTimeToSleep int
)

// WORKERS
var (
	PulseWorkersNumber int
	PulseMessageBuffer int

	SavePulseWorkersNumber int
	SavePulseMessageBuffer int
	SavePulseBatch         int

	ScheduleSavePulse string
)

// DATABASE
var (
	IngesterDb database.DbConfig
)

// OTHERS
var (
	Flag string
)

func LoadEnv() {
	var err error
	if err = godotenv.Load(); err != nil {
		logrus.Info("runnning the application without a .env file")
	}

	//RABBITMQ
	PulseQueueUrl = os.Getenv("PULSE_QUEUE_URL")
	PulseQueueName = os.Getenv("PULSE_QUEUE_NAME")

	//REDIS
	Addr = os.Getenv("ADDR")
	RedisPassword = os.Getenv("PASSWORD")
	RedisDb, err = utils.StringToInt(os.Getenv("REDIS_DB"))
	if err != nil {
		logrus.Error("Fail to convert RedisDb to int.Erro:", err)
	}

	RedisMaxRetry, err = utils.StringToInt(os.Getenv("REDIS_MAX_RETRY"))
	if err != nil {
		logrus.Error("Fail to convert RedisMaxRetry to int. Erro:", err)
	}

	RedisTimeToSleep, err = utils.StringToInt(os.Getenv("REDIS_TIME_TO_SLEEP"))
	if err != nil {
		logrus.Error("Fail to convert RedisTimeToSleep to int. Erro:", err)
	}

	//WORKERS
	PulseWorkersNumber, err = utils.StringToInt(os.Getenv("PULSE_WORKERS_NUMBER"))
	if err != nil {
		logrus.Error("Fail to convert PulseWorkersNumber to int. Erro:", err)
	}

	PulseMessageBuffer, err = utils.StringToInt(os.Getenv("PULSE_MESSAGE_BUFFER"))
	if err != nil {
		logrus.Error("Fail to convert PulseMessageBuffer to int. Erro:", err)
	}

	SavePulseWorkersNumber, err = utils.StringToInt(os.Getenv("SAVE_PULSE_WORKERS_NUMBER"))
	if err != nil {
		logrus.Error("Fail to convert SavePulseWorkersNumber to int. Erro:", err)
	}

	SavePulseMessageBuffer, err = utils.StringToInt(os.Getenv("SAVE_PULSE_MESSAGE_BUFFER"))
	if err != nil {
		logrus.Error("Fail to convert SavePulseMessageBuffer to int. Erro:", err)
	}

	SavePulseBatch, err = utils.StringToInt(os.Getenv("SAVE_PULSE_BATCH"))
	if err != nil {
		logrus.Error("Fail to convert SavePulseBatch to int. Erro:", err)
	}

	//DATABASE
	IngesterDb.Host = os.Getenv("INGESTER_HOST")
	IngesterDb.User = os.Getenv("INGESTER_USER")
	IngesterDb.Password = os.Getenv("INGESTER_PASSWORD")
	IngesterDb.DbName = os.Getenv("INGESTER_DB_NAME")
	IngesterDb.Port, err = utils.StringToInt64(os.Getenv("INGESTER_PORT"))
	if err != nil {
		IngesterDb.Port = 5432
		logrus.Error("Fail to convert Ingester DB to int. Erro:", err)
	}

	//CRON
	ScheduleSavePulse = os.Getenv("SCHEDULE_SAVE_PULSE")

	//OTHERS
	Flag = os.Getenv("FLAG")
}
