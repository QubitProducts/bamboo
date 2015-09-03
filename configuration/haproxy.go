package configuration

type HAProxy struct {
	TemplatePath    string
	OutputPath      string
	ReloadCommand   string
	ShutdownCommand string
	GraceSeconds int
}
