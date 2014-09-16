package event_bus

import (
	"log"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/QubitProducts/bamboo/configuration"
	"github.com/QubitProducts/bamboo/services/haproxy"
	"os/exec"
	"github.com/QubitProducts/bamboo/writer"
	"io/ioutil"
)

type MarathonEvent struct {
	// EventType can be
	// api_post_event, status_update_event, subscribe_event
	EventType string
	Timestamp string
}

type ZookeeperEvent struct {
	Source string
	EventType string
}

type ServiceEvent struct {
	EventType string
}

type Handlers struct {
	Conf *configuration.Configuration
	Zookeeper *zk.Conn
}

func (h *Handlers) MarathonEventHandler(event MarathonEvent) {
	log.Printf("%s => %s\n", event.EventType, event.Timestamp)
	handleHAPUpdate(h.Conf, h.Zookeeper)
	h.Conf.StatsD.Increment(1.0, "reload.marathon", 1)
}

func (h *Handlers) ServiceEventHandler(event ServiceEvent) {
	log.Println("Domain mapping: Stated changed")
	handleHAPUpdate(h.Conf, h.Zookeeper)
	h.Conf.StatsD.Increment(1.0, "reload.domain", 1)
}

func handleHAPUpdate(conf *configuration.Configuration, conn *zk.Conn) bool {
	currentContent, _ := ioutil.ReadFile(conf.HAProxy.OutputPath)

	templateContent, err := ioutil.ReadFile(conf.HAProxy.TemplatePath)
	if err != nil { log.Panicf("Cannot read template file: %s", err) }

	templateData := haproxy.GetTemplateData(conf, conn)

	newContent, err := writer.RenderTemplate(conf.HAProxy.TemplatePath, string(templateContent), templateData)

	if err != nil { log.Fatalf("Template syntax error: \n %s", err ) }

	if (currentContent == nil || string(currentContent) != newContent) {
		err := ioutil.WriteFile(conf.HAProxy.OutputPath, []byte(newContent), 0666)
		if err != nil { log.Fatalf("Failed to write template on path: %s", err) }

		execCommand(conf.HAProxy.ReloadCommand)
		log.Println("HAProxy: Configuration updated")
		return true
	} else {
		log.Println("HAProxy: Same content, no need to reload")
		return false
	}
}

func execCommand(cmd string) {
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Exec cmd: %s \n", cmd)
}
