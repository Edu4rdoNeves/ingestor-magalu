package pulse

import (
	"context"
	"errors"
	"fmt"

	"github.com/Edu4rdoNeves/ingestor-magalu/domain/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IPulseRepository interface {
	SavePulseBatch(ctx context.Context, pulses []*entity.PulseData) error
	GetPulses(offset, limit int) ([]*entity.PulseData, error)
}

type PulseRepository struct {
	DB *gorm.DB
}

func NewPulseRepository(db *gorm.DB) *PulseRepository {
	return &PulseRepository{
		DB: db,
	}
}

func (r *PulseRepository) SavePulseBatch(ctx context.Context, pulseEntities []*entity.PulseData) error {
	if len(pulseEntities) == 0 {
		return errors.New("no pulse reported")
	}

	return r.DB.Debug().WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "tenant"},
			{Name: "product_sku"},
			{Name: "use_unity"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"used_amount": gorm.Expr("pulses.used_amount + EXCLUDED.used_amount"),
			"updated_at":  gorm.Expr("NOW()"),
		}),
	}).Create(&pulseEntities).Error
}

func (r *PulseRepository) GetPulses(offset, limit int) ([]*entity.PulseData, error) {
	pulses := []*entity.PulseData{}

	err := r.DB.Limit(limit).Offset(offset).Find(&pulses).Error
	if err != nil {
		logrus.Error("fail to get pulses")
		return nil, fmt.Errorf("fail to get pulses. Error: %v", err)

	}

	return pulses, nil
}
