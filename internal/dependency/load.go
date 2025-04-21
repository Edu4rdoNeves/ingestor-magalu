package dependency

func Load() {
	LoadDataBases()
	LoadGeneral()
	LoadWorkerDependencies()
}
