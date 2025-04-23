package dependency

import (
	"github.com/Edu4rdoNeves/ingestor-magalu/cmd/api/controller/pulse"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/dependency/api"
)

var (
	PulseDependency pulse.IPulseController
)

func LoadApiDependency() {
	PulseDependency = api.PulseDependency()
}
