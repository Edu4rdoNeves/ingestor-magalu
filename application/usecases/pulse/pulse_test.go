package pulse_test

import (
	"context"
	"errors"
	"testing"

	pulseUsecase "github.com/Edu4rdoNeves/ingestor-magalu/application/usecases/pulse"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	mockpulseRepo "github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/repository/pulse/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUseCase_SavePulseBatch(t *testing.T) {
	t.Run("PulseUseCase/SavePulseBatch", func(t *testing.T) {
		tests := []struct {
			name        string
			input       []*dto.PulseData
			mockSetup   func(repo *mockpulseRepo.MockIPulseRepository)
			expectedErr string
		}{
			{
				name: "✅ Deve salvar com sucesso",
				input: []*dto.PulseData{
					{
						Tenant:     "tenant-1",
						ProductSku: "sku-123",
						UseUnity:   "kg",
						UsedAmount: 5,
					},
				},
				mockSetup: func(repo *mockpulseRepo.MockIPulseRepository) {
					repo.EXPECT().
						SavePulseBatch(gomock.Any(), gomock.Any()).
						Return(nil)
				},
			},
			{
				name:  "❌ Deve retornar erro do repositório",
				input: []*dto.PulseData{},
				mockSetup: func(repo *mockpulseRepo.MockIPulseRepository) {
					repo.EXPECT().
						SavePulseBatch(gomock.Any(), gomock.Any()).
						Return(errors.New("erro simulado no repositório"))
				},
				expectedErr: "erro simulado no repositório",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockRepo := mockpulseRepo.NewMockIPulseRepository(ctrl)

				if tt.mockSetup != nil {
					tt.mockSetup(mockRepo)
				}

				usecase := pulseUsecase.NewPulseUseCase(mockRepo)

				err := usecase.SavePulseBatch(context.Background(), tt.input)

				if tt.expectedErr != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
