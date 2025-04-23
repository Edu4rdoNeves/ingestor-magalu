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

	SimulatorTotalMessages int
	SimulatorWorkersNumber int
	SimulatorBufferSize    int

	ScheduleSavePulse string
)

// API
var (
	GinMode    string
	CORSEnable bool
	ServerPort string
)

// DATABASE
var (
	IngestorDb database.DbConfig
)

// OTHERS
var (
	QntdProductSku int
	QntdUseUnity   int
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

	//API
	GinMode = os.Getenv("GIN_MODE")
	ServerPort = os.Getenv("SERVER_PORT")

	//DATABASE
	IngestorDb.Host = os.Getenv("INGESTOR_HOST")
	IngestorDb.User = os.Getenv("INGESTOR_USER")
	IngestorDb.Password = os.Getenv("INGESTOR_PASSWORD")
	IngestorDb.DbName = os.Getenv("INGESTOR_DB_NAME")
	IngestorDb.Port, err = utils.StringToInt64(os.Getenv("INGESTOR_PORT"))
	if err != nil {
		IngestorDb.Port = 5432
		logrus.Error("Fail to convert Ingestor DB to int. Erro:", err)
	}

	//CRON
	ScheduleSavePulse = os.Getenv("SCHEDULE_SAVE_PULSE")

	//OTHERS
	QntdProductSku, err = utils.StringToInt(os.Getenv("QNTD_PRODUCT_SKU"))
	if err != nil {
		logrus.Error("Fail to convert QntdProductSku to int. Erro:", err)
	}

	QntdUseUnity, err = utils.StringToInt(os.Getenv("QNTD_USE_UNITY"))
	if err != nil {
		logrus.Error("Fail to convert QntdUseUnity to int. Erro:", err)
	}

	SimulatorTotalMessages, err = utils.StringToInt(os.Getenv("SIMULATOR_TOTAL_MESSAGES"))
	if err != nil {
		logrus.Error("Fail to convert SimulatorTotalMessages to int. Erro:", err)
	}

	SimulatorWorkersNumber, err = utils.StringToInt(os.Getenv("SIMULATOR_WORKERS_NUMBER"))
	if err != nil {
		logrus.Error("Fail to convert SimulatorWorkersNumber to int. Erro:", err)
	}

	SimulatorBufferSize, err = utils.StringToInt(os.Getenv("SIMULATOR_BUFFER_SIZE"))
	if err != nil {
		logrus.Error("Fail to convert SimulatorBufferSize to int. Erro:", err)
	}

}
