package utils

import (
	"fmt"
	"strings"

	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/constants"
	"github.com/sirupsen/logrus"
)

func ParsePulseKey(key string, amount float64) (*dto.PulseData, error) {
	parts := strings.Split(key, ":")
	const expectedParts = 4

	if len(parts) != expectedParts {
		return nil, fmt.Errorf("invalid key format: expected %d parts, got %d - key: %s", expectedParts, len(parts), key)
	}

	return &dto.PulseData{
		Tenant:     parts[1],
		ProductSku: parts[2],
		UseUnity:   parts[3],
		UsedAmount: amount,
	}, nil
}

func ValidateAndSetDefaults(populateQueue *dto.PopulateQueueParams) error {
	if populateQueue.TotalMessages <= 0 {
		logrus.Warn("TotalMessages não especificado ou inválido. Usando valor padrão.")
		populateQueue.TotalMessages = constants.DefaultTotalMessages
	}

	if populateQueue.WorkersNumber <= 0 {
		logrus.Warn("WorkersNumber não especificado ou inválido. Usando valor padrão.")
		populateQueue.WorkersNumber = constants.DefaultWorkersNumber
	}

	if populateQueue.BufferSize <= 0 {
		logrus.Warn("BufferSize não especificado ou inválido. Usando valor padrão.")
		populateQueue.BufferSize = constants.DefaultBufferSize
	}

	return nil
}
