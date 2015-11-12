package event_bus

import (
	"github.com/QubitProducts/bamboo/configuration"
	"github.com/QubitProducts/bamboo/services/haproxy"
	"github.com/QubitProducts/bamboo/services/service"
	"github.com/QubitProducts/bamboo/services/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

type MarathonEvent struct {
	// EventType can be
	// api_post_event, status_update_event, subscribe_event
	EventType string
	Timestamp string
}

type ZookeeperEvent struct {
	Source    string
	EventType string
}

type ServiceEvent struct {
	EventType string
}

type Handlers struct {
	Conf    *configuration.Configuration
	Storage service.Storage
}

func (h *Handlers) MarathonEventHandler(event MarathonEvent) {
	log.Printf("%s => %s\n", event.EventType, event.Timestamp)
	queueUpdate(h)
	h.Conf.StatsD.Increment(1.0, "callback.marathon", 1)
}

func (h *Handlers) ServiceEventHandler(event ServiceEvent) {
	log.Println("Domain mapping: Stated changed")
	queueUpdate(h)
	h.Conf.StatsD.Increment(1.0, "reload.domain", 1)
}

var updateChan = make(chan *Handlers, 1)

func init() {
	go func() {
		log.Println("Starting update loop")
		for {
			h := <-updateChan
			handleHAPUpdate(h.Conf, h.Storage)
		}
	}()
}

var queueUpdateSem = make(chan int, 1)

func queueUpdate(h *Handlers) {
	queueUpdateSem <- 1

	select {
	case _ = <-updateChan:
		log.Println("Found pending update request. Don't start another one.")
	default:
		log.Println("Queuing an haproxy update.")
	}
	updateChan <- h

	<-queueUpdateSem
}

func handleHAPUpdate(conf *configuration.Configuration, storage service.Storage) {
	reloadStart := time.Now()
	reloaded, err := ensureLatestConfig(conf, storage)

	if err != nil {
		conf.StatsD.Increment(1.0, "haproxy.reload.error", 1)
		log.Println("Failed to update HAProxy configuration:", err)
	} else if reloaded {
		conf.StatsD.Timing(1.0, "haproxy.reload.marathon.duration", time.Since(reloadStart))
		conf.StatsD.Increment(1.0, "haproxy.reload.marathon.reloaded", 1)
		log.Println("Reloaded HAProxy configuration")
	} else {
		conf.StatsD.Increment(1.0, "haproxy.reload.skipped", 1)
		log.Println("Skipped HAProxy configuration reload due to lack of changes")
	}
}

// For values of 'latest' conforming to general relativity.
func ensureLatestConfig(conf *configuration.Configuration, storage service.Storage) (reloaded bool, err error) {
	content, err := generateConfig(conf.HAProxy.TemplatePath, conf, storage)
	if err != nil {
		return
	}

	req, err := isReloadRequired(conf.HAProxy.OutputPath, content)
	if err != nil || !req {
		return
	}

	err = validateConfig(conf.HAProxy.ReloadValidationCommand, content)
	if err != nil {
		return
	}

	defer cleanupConfig(conf.HAProxy.ReloadCleanupCommand)

	reloaded, err = changeConfig(conf, content)
	if err != nil {
		return
	}

	return
}

// Generates the new config to be written
func generateConfig(templatePath string, conf *configuration.Configuration, storage service.Storage) (config string, err error) {
	templateContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Println("Failed to read template contents")
		return
	}

	templateData, err := haproxy.GetTemplateData(conf, storage)
	if err != nil {
		log.Println("Failed to retrieve template data")
		return
	}

	config, err = template.RenderTemplate(templatePath, string(templateContent), templateData)
	if err != nil {
		log.Println("Template syntax error")
		return
	}

	return
}

// Loads the existing config and decides if a reload is required
func isReloadRequired(configPath string, newContent string) (bool, error) {
	// An error here means that the template may not exist, in which case we simply continue
	currentContent, err := ioutil.ReadFile(configPath)

	if err == nil {
		return newContent != string(currentContent), nil
	} else if os.IsNotExist(err) {
		return true, nil
	}

	return false, err // Returning false here as is default value for bool
}

// Takes the ReloadValidateCommand and returns nil if the command succeeded
func validateConfig(validateTemplate string, newContent string) (err error) {
	if validateTemplate == "" {
		return nil
	}

	tmpFile, err := ioutil.TempFile("/tmp", "bamboo")
	if err != nil {
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	log.Println("Generating validation command")
	_, err = tmpFile.WriteString(newContent)
	if err != nil {
		return
	}

	validateCommand, err := template.RenderTemplate(
		"validate",
		validateTemplate,
		tmpFile.Name())
	if err != nil {
		return
	}

	log.Println("Validating config")
	err = execCommand(validateCommand)

	return
}

func changeConfig(conf *configuration.Configuration, newContent string) (reloaded bool, err error) {
	// This failing scares me a lot, as could end up with very invalid config
	// content. I'd suggest restoring the original config, but that adds all
	// kinds of new and interesting failure cases
	err = ioutil.WriteFile(conf.HAProxy.OutputPath, []byte(newContent), 0666)
	if err != nil {
		log.Println("Failed to write template on path", conf.HAProxy.OutputPath)
		return
	}

	err = execCommand(conf.HAProxy.ReloadCommand)
	if err != nil {
		return
	}

	reloaded = true
	return
}

// This will be executed in a deferred, so is rather self contained
func cleanupConfig(command string) {
	log.Println("Cleaning up config")
	execCommand(command)
}

func execCommand(cmd string) error {
	log.Printf("Exec cmd: %s \n", cmd)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Println(err.Error())
		log.Println("Output:\n" + string(output[:]))
	}
	return err
}
