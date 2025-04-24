package pulse

import (
	"context"

	populatequeuetask "github.com/Edu4rdoNeves/ingestor-magalu/cmd/api/task/populate_queue_task"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/mapper"
	pulseRepo "github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/repository/pulse"
	"github.com/Edu4rdoNeves/ingestor-magalu/utils"
	"github.com/sirupsen/logrus"
)

type IPulseUseCase interface {
	SavePulseBatch(ctx context.Context, pulseDto []*dto.PulseData) error
	GetPulses(page, limit int) ([]*dto.PulseData, error)
	GetPulseByID(id int) (*dto.PulseData, error)
	PopulateQueueWithPulses(populateQueue *dto.PopulateQueueParams) error
}

type PulseUseCase struct {
	pulseRepo pulseRepo.IPulseRepository
	task      populatequeuetask.IPopulateQueueTask
}

func NewPulseUseCase(pulseRepo pulseRepo.IPulseRepository, task populatequeuetask.IPopulateQueueTask) IPulseUseCase {
	return &PulseUseCase{
		pulseRepo: pulseRepo,
		task:      task,
	}
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

	convertedPulses := mapper.PulsesEntityToPulsesDto(pulseResponse)

	return convertedPulses, nil
}

func (uc *PulseUseCase) GetPulseByID(id int) (*dto.PulseData, error) {
	pulseResponse, err := uc.pulseRepo.GetPulseByID(id)
	if err != nil {
		return nil, err
	}

	convertedPulse := mapper.PulseEntityToPulseDto(pulseResponse)

	return convertedPulse, nil
}

func (uc *PulseUseCase) PopulateQueueWithPulses(populateQueue *dto.PopulateQueueParams) error {
	if err := utils.ValidateAndSetDefaults(populateQueue); err != nil {
		return err
	}

	go uc.task.Run(populateQueue)
	return nil
}
