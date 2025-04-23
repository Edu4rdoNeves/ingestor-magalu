package dependency

import "github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"

func Load() {
	env.LoadEnv()
	LoadApiDependency()
	LoadDataBases()
	LoadGeneral()
	LoadWorkerDependencies()
}
