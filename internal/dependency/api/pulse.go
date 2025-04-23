package api

import (
	pulseUC "github.com/Edu4rdoNeves/ingestor-magalu/application/usecases/pulse"
	pulseController "github.com/Edu4rdoNeves/ingestor-magalu/cmd/api/controller/pulse"
	"github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/database"
	pulseRepository "github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/repository/pulse"
)

func PulseDependency() pulseController.IPulseController {
	pulseRepository := pulseRepository.NewPulseRepository(database.Get())
	pulseUseCases := pulseUC.NewPulseUseCase(pulseRepository)
	pulseController := pulseController.NewPulseController(pulseUseCases)

	return pulseController
}
