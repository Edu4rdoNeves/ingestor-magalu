package utils

import (
	"flag"

	"github.com/Edu4rdoNeves/ingestor-magalu/internal/constants"
)

type Flags struct {
	RunWorker bool
	RunAPI    bool
}

func ConfigFlags() *Flags {
	var flags Flags

	flag.BoolVar(&flags.RunWorker, constants.Worker, false, constants.RunWorker)
	flag.BoolVar(&flags.RunAPI, constants.Api, false, constants.RunApi)

	flag.Parse()

	return &flags
}
