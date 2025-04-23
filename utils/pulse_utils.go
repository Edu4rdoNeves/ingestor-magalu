package utils

import (
	"fmt"
	"strings"

	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
)

func ParsePulseKey(key string, count float64) (*dto.PulseData, error) {
	parts := strings.Split(key, ":")
	const expectedParts = 4

	if len(parts) != expectedParts {
		return nil, fmt.Errorf("invalid key format: expected %d parts, got %d - key: %s", expectedParts, len(parts), key)
	}

	return &dto.PulseData{
		Tenant:     parts[1],
		ProductSku: parts[2],
		UseUnity:   parts[3],
		UsedAmount: count,
	}, nil
}
