package mapper

import (
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/entity"
)

func PulseEntityToPulseDto(entities []*entity.PulseData) []*dto.PulseData {
	dtoPulse := make([]*dto.PulseData, 0, len(entities))
	for _, d := range entities {
		dtoPulse = append(dtoPulse, &dto.PulseData{
			Tenant:     d.Tenant,
			ProductSku: d.ProductSku,
			UseUnity:   d.UseUnity,
			UsedAmount: d.UsedAmount,
		})
	}
	return dtoPulse
}
