package pulsetask_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	rabbitMock "github.com/Edu4rdoNeves/ingestor-magalu/application/service/rabbitmq/mocks"
	redisMock "github.com/Edu4rdoNeves/ingestor-magalu/application/service/redis/mocks"
	pulsetask "github.com/Edu4rdoNeves/ingestor-magalu/cmd/worker/task/pulse_task"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/dto"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/golang/mock/gomock"
)

func TestPulseTask_Run(t *testing.T) {
	t.Run("PulseTask/Run", func(t *testing.T) {
		tests := []struct {
			name            string
			setMocks        func(mockRedis *redisMock.MockIRedisClient, mockRabbit *rabbitMock.MockIRabbitMQ, pulseBody []byte, wg *sync.WaitGroup)
			expectErrLog    bool
			expectRedisCall bool
		}{
			{
				name: "✅ Deve processar mensagem e chamar Redis com sucesso",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockRabbit *rabbitMock.MockIRabbitMQ, pulseBody []byte, wg *sync.WaitGroup) {
					mockRabbit.EXPECT().
						Consumer(gomock.Any()).
						DoAndReturn(func(handler func([]byte)) error {
							go func() {
								time.Sleep(10 * time.Millisecond)
								handler(pulseBody)
							}()
							return nil
						})

					mockRedis.EXPECT().
						IncrementCounter(gomock.Any(), gomock.Any()).
						DoAndReturn(func(key string, val float64) error {
							t.Logf("✅ Redis chamado com sucesso - key=%s val=%.2f", key, val)
							wg.Done()
							return nil
						})
				},
				expectRedisCall: true,
			},
			{
				name: "❌ Deve logar erro se der falha ao consumir a mensagem no rabbitMq",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockRabbit *rabbitMock.MockIRabbitMQ, pulseBody []byte, wg *sync.WaitGroup) {
					mockRabbit.EXPECT().
						Consumer(gomock.Any()).
						DoAndReturn(func(handler func([]byte)) error {
							t.Logf("🚨 RabbitMq retornou erro simulado")
							return errors.New("erro simulado ao consumir mensagens do RabbitMq")
						})
				},
				expectErrLog:    true,
				expectRedisCall: false,
			},
			{
				name: "❌ Deve logar erro se Redis falhar",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockRabbit *rabbitMock.MockIRabbitMQ, pulseBody []byte, wg *sync.WaitGroup) {
					mockRabbit.EXPECT().
						Consumer(gomock.Any()).
						DoAndReturn(func(handler func([]byte)) error {
							go func() {
								time.Sleep(10 * time.Millisecond)
								handler(pulseBody)
							}()
							return nil
						})

					mockRedis.EXPECT().
						IncrementCounter(gomock.Any(), gomock.Any()).
						DoAndReturn(func(key string, val float64) error {
							t.Logf("🚨 Redis retornou erro simulado")
							wg.Done()
							return errors.New("erro simulado no Redis")
						})
				},
				expectErrLog:    true,
				expectRedisCall: true,
			},
			{
				name: "❌ Deve logar erro se JSON for inválido",
				setMocks: func(mockRedis *redisMock.MockIRedisClient, mockRabbit *rabbitMock.MockIRabbitMQ, pulseBody []byte, wg *sync.WaitGroup) {
					pulseBody = []byte("isso não é um json válido")

					mockRabbit.EXPECT().
						Consumer(gomock.Any()).
						DoAndReturn(func(handler func([]byte)) error {
							go func() {
								time.Sleep(10 * time.Millisecond)
								handler(pulseBody)
							}()
							return nil
						})
				},
				expectRedisCall: false,
			},
		}

		env.PulseWorkersNumber = 1
		env.PulseMessageBuffer = 1
		env.RedisMaxRetry = 1
		env.RedisTimeToSleep = 1

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockRedis := redisMock.NewMockIRedisClient(ctrl)
				mockRabbit := rabbitMock.NewMockIRabbitMQ(ctrl)

				pulse := &dto.PulseData{
					Tenant:     "magalu",
					ProductSku: "SKU-1",
					UseUnity:   "loja-1",
					UsedAmount: 5.0,
				}
				body, _ := json.Marshal(pulse)

				var wg sync.WaitGroup

				tc.setMocks(mockRedis, mockRabbit, body, &wg)

				task := pulsetask.NewPulseTask(mockRedis, mockRabbit)

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				go task.Run(ctx)

				done := make(chan struct{})
				go func() {
					wg.Wait()
					close(done)
				}()

				if tc.expectRedisCall {
					wg.Add(1)
				}

				if tc.expectRedisCall {
					select {
					case <-done:
						t.Log("✅ Worker executado")
					case <-time.After(2 * time.Second):
						t.Fatal("❌ Timeout: Redis não foi chamado")
					}
				} else {
					time.Sleep(200 * time.Millisecond)
					t.Log("✅ Handler executado")
				}
			})
		}
	})
}
