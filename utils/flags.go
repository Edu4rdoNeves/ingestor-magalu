package utils

import (
	"flag"

	"github.com/Edu4rdoNeves/ingestor-magalu/internal/constants"
)

func ConfigFlags() (*bool, *bool) {
	workerFlag := flag.Bool(constants.Worker, false, constants.RunWorker)
	scriptFlag := flag.Bool(constants.Script, false, constants.RunScript)

	return workerFlag, scriptFlag
}
