package env

func IsScript() bool {
	return Flag == "script"
}

func IsWorker() bool {
	return Flag == "worker"
}
