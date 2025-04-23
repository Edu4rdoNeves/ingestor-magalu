package pulse_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/entity"
	pulseRepository "github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/repository/pulse"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestRepository_SavePulseBatch(t *testing.T) {

	mockDb, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDb.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(t, err)

	gormDB = gormDB.Debug()

	repo := pulseRepository.NewPulseRepository(gormDB)
	t.Run("✅ Deve Retornar sucesso ao gravar dados no banco", func(t *testing.T) {

		pulses := []*entity.PulseData{
			{
				Tenant:     "tenant-1",
				ProductSku: "sku-123",
				UseUnity:   "kg",
				UsedAmount: 5.0,
				CreatedAt:  time.Now(),
			},
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "pulses"`).
			WithArgs(
				pulses[0].Tenant,
				pulses[0].ProductSku,
				pulses[0].UseUnity,
				pulses[0].UsedAmount,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err = repo.SavePulseBatch(context.Background(), pulses)
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
	t.Run("❌ Deve Retornar erro ao falhar no INSERT", func(t *testing.T) {

		pulses := []*entity.PulseData{
			{
				Tenant:     "tenant-1",
				ProductSku: "sku-123",
				UseUnity:   "kg",
				UsedAmount: 5.0,
			},
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "pulses"`).
			WithArgs(
				pulses[0].Tenant,
				pulses[0].ProductSku,
				pulses[0].UseUnity,
				pulses[0].UsedAmount,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("erro no banco"))
		mock.ExpectRollback()

		err = repo.SavePulseBatch(context.Background(), pulses)
		require.Error(t, err)
		require.Contains(t, err.Error(), "erro no banco")

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
	t.Run("❌ Deve Retornar erro ao quando receber um parametro nulo", func(t *testing.T) {
		err = repo.SavePulseBatch(context.Background(), nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no pulse reported")

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}
