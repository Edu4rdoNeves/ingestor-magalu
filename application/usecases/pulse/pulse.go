package pulse

import (
	"context"

	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/mapper"
	pulseRepo "github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/repository/pulse"
	"github.com/sirupsen/logrus"
)

type IPulseUseCase interface {
	SavePulseBatch(ctx context.Context, pulseDto []*dto.PulseData) error
	GetPulses(page, limit int) ([]*dto.PulseData, error)
}

type PulseUseCase struct {
	pulseRepo pulseRepo.IPulseRepository
}

func NewPulseUseCase(pulseRepo pulseRepo.IPulseRepository) IPulseUseCase {
	return &PulseUseCase{pulseRepo: pulseRepo}
}

func (uc *PulseUseCase) SavePulseBatch(ctx context.Context, pulseDto []*dto.PulseData) error {
	entityPulse := mapper.PulseDtosToEntities(pulseDto)

	err := uc.pulseRepo.SavePulseBatch(ctx, entityPulse)
	if err != nil {
		logrus.Errorf("Failed to save pulse to database.Error: %v", err)
		return err
	}

	return nil
}

func (uc *PulseUseCase) GetPulses(page, limit int) ([]*dto.PulseData, error) {
	offset := (page - 1) * limit

	pulseResponse, err := uc.pulseRepo.GetPulses(offset, limit)
	if err != nil {
		return nil, err
	}

	convertedPulses := mapper.PulseEntityToPulseDto(pulseResponse)

	return convertedPulses, nil
}
