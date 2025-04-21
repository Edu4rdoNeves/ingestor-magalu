package utils

import (
	"fmt"
	"strings"

	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
)

func ParsePulseKey(key string, count float64) (*dto.PulseData, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid key format: %s", key)
	}

	return &dto.PulseData{
		Tenant:     parts[1],
		ProductSku: parts[2],
		UseUnity:   parts[3],
		UsedAmount: count,
	}, nil
}
