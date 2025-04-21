package mapper

import (
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/entity"
)

func PulseDtosToEntities(dtos []*dto.PulseData) []*entity.PulseData {
	entities := make([]*entity.PulseData, 0, len(dtos))
	for _, d := range dtos {
		entities = append(entities, &entity.PulseData{
			Tenant:     d.Tenant,
			ProductSku: d.ProductSku,
			UseUnity:   d.UseUnity,
			UsedAmount: d.UsedAmount,
		})
	}
	return entities
}
