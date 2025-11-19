package definition

type ModuleIntface interface {
	Start()     (error)
	PostStart()

	PreStop()
	Stop()      (error) // TODO: no error ?

	State() (string) // TODO

	// Debug()
}
