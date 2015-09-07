package configuration

type HAProxy struct {
	TemplatePath            string
	OutputPath              string
	ReloadCommand           string
	ReloadValidationCommand string
	ReloadCleanupCommand    string
}
