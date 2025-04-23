package savepulsetask_test

import (
	"context"
	"errors"
	"testing"

	redisMock "github.com/Edu4rdoNeves/ingestor-magalu/application/service/redis/mocks"
	usecaseMock "github.com/Edu4rdoNeves/ingestor-magalu/application/usecases/pulse/mocks"
	savePulsetask "github.com/Edu4rdoNeves/ingestor-magalu/cmd/worker/task/save_pulse_task"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/golang/mock/gomock"
)

func TestSavePulseTask_Run(t *testing.T) {
	t.Run("SavePulseTask/Run", func(t *testing.T) {
		tests := []struct {
			name        string
			setMocks    func(mockRedis *redisMock.MockIRedisClient, mockUsecase *usecaseMock.MockIPulseUseCase)
			expectedLog string
		}{
			{
				name: "✅ Deve processar todas as chaves com sucesso",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockUsecase *usecaseMock.MockIPulseUseCase) {
					mockRedis.EXPECT().GetKeysByPattern("pulse:*").Return([]string{"pulse:tenant:sku:unity"}, nil)

					mockRedis.EXPECT().GetValue("pulse:tenant:sku:unity").Return("3.5", nil)

					mockUsecase.EXPECT().SavePulseBatch(gomock.Any(), gomock.Any()).Return(nil)
					mockRedis.EXPECT().DeleteKey("pulse:tenant:sku:unity").Return(nil)
				},
				expectedLog: "Batch saved and deleted successfully",
			},
			{
				name: "❌ Deve retornar erro ao pegar as chaves pelo indentificador",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockUsecase *usecaseMock.MockIPulseUseCase) {
					mockRedis.EXPECT().GetKeysByPattern(gomock.Any()).Return(nil, errors.New("erro simulado no Redis"))
				},
				expectedLog: "Failed to get keys",
			},
			{
				name: "❌ Deve retornar erro ao pegar as pegar o valor da chave no redis",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockUsecase *usecaseMock.MockIPulseUseCase) {
					mockRedis.EXPECT().GetKeysByPattern("pulse:*").Return([]string{"pulse:tenant:sku:unity"}, nil)

					mockRedis.EXPECT().GetValue(gomock.Any()).Return("nil", errors.New("erro simulado no GetValue"))
				},
				expectedLog: "Failed to get value for keys",
			},
			{
				name: "❌ Deve retornar erro ao converter de string para float64",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockUsecase *usecaseMock.MockIPulseUseCase) {
					mockRedis.EXPECT().GetKeysByPattern("pulse:*").Return([]string{"pulse:tenant:sku:unity"}, nil)

					mockRedis.EXPECT().GetValue(gomock.Any()).Return("não é número", nil)
				},
				expectedLog: "error when try converting string to float64",
			},
			{
				name: "❌ Deve retornar erro ao tentar fazer o parse da chave",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockUsecase *usecaseMock.MockIPulseUseCase) {
					mockRedis.EXPECT().GetKeysByPattern("pulse:*").Return([]string{"pulse:tenant:sku"}, nil)

					mockRedis.EXPECT().GetValue("pulse:tenant:sku").Return("1", nil)
				},
				expectedLog: "Failed to parse key",
			},
			{
				name: "❌ Deve retornar erro ao tentar salvar o lote no banco",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockUsecase *usecaseMock.MockIPulseUseCase) {
					mockRedis.EXPECT().GetKeysByPattern("pulse:*").Return([]string{"pulse:tenant:sku:unity"}, nil)

					mockRedis.EXPECT().GetValue("pulse:tenant:sku:unity").Return("3.5", nil)

					mockUsecase.EXPECT().SavePulseBatch(gomock.Any(), gomock.Any()).Return(errors.New("erro simulado no SavePulseBatch"))
				},
				expectedLog: "Failed to save batch",
			},
			{
				name: "❌ Deve retornar erro ao tentar deletar os dados do  Redis",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockUsecase *usecaseMock.MockIPulseUseCase) {
					mockRedis.EXPECT().GetKeysByPattern("pulse:*").Return([]string{"pulse:tenant:sku:unity"}, nil)

					mockRedis.EXPECT().GetValue("pulse:tenant:sku:unity").Return("3.5", nil)

					mockUsecase.EXPECT().SavePulseBatch(gomock.Any(), gomock.Any()).Return(nil)

					mockRedis.EXPECT().DeleteKey("pulse:tenant:sku:unity").Return(errors.New("erro simulado no DeleteKey"))
				},
				expectedLog: "Failed to delete key",
			},
		}

		env.SavePulseWorkersNumber = 1
		env.SavePulseBatch = 0
		env.SavePulseMessageBuffer = 10
		env.RedisMaxRetry = 1
		env.RedisTimeToSleep = 1

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockRedis := redisMock.NewMockIRedisClient(ctrl)
				mockUsecase := usecaseMock.NewMockIPulseUseCase(ctrl)

				tc.setMocks(mockRedis, mockUsecase)

				task := savePulsetask.NewSavePulseTask(mockRedis, mockUsecase)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				task.Run(ctx)

				// Aqui você pode validar logs ou usar hooks de logrus pra interceptar mensagens se necessário
			})
		}
	})
}
