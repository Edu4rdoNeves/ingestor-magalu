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
}

type PulseUseCase struct {
	PulseRepo pulseRepo.IPulseRepository
}

func NewPulseUseCase(pulseRepo pulseRepo.IPulseRepository) IPulseUseCase {
	return &PulseUseCase{PulseRepo: pulseRepo}
}

func (uc *PulseUseCase) SavePulseBatch(ctx context.Context, pulseDto []*dto.PulseData) error {
	entityPulse := mapper.PulseDtosToEntities(pulseDto)

	err := uc.PulseRepo.SavePulseBatch(ctx, entityPulse)
	if err != nil {
		logrus.Errorf("Failed to save pulse to database.Error: %v", err)
		return err
	}

	return nil
}
