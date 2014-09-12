package event_bus

import (
	"log"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/QubitProducts/bamboo/configuration"
	"github.com/QubitProducts/bamboo/services/haproxy"
	"os/exec"
)

type MarathonEvent struct {
	// EventType can be
	// api_post_event, status_update_event, subscribe_event
	EventType string `param:"id"`
	Timestamp string `param:"timestamp"`
}

type ZookeeperEvent struct {
	Source string
	EventType string
}

type DomainEvent struct {
	EventType string
}

type Handlers struct {
	Conf *configuration.Configuration
	Zookeeper *zk.Conn
}

func (h *Handlers) MarathonEventHandler(event MarathonEvent) {
	log.Printf("%s => %s\n", event.EventType, event.Timestamp)
	handleHAPUpdate(h.Conf, h.Zookeeper)
	execCommand(h.Conf.HAProxy.ReloadCommand)
	h.Conf.StatsD.Increment(1.0, "reload.marathon", 1)
}

func (h *Handlers) DomainEventHandler(event DomainEvent) {
	log.Println("Domain mapping: Stated changed")
	handleHAPUpdate(h.Conf, h.Zookeeper)
	execCommand(h.Conf.HAProxy.ReloadCommand)
	h.Conf.StatsD.Increment(1.0, "reload.domain", 1)
}

func handleHAPUpdate(conf *configuration.Configuration, conn *zk.Conn) {
	err := haproxy.WriteHAProxyConfig(conf.HAProxy, haproxy.GetTemplateData(conf, conn))
	if err != nil {
		log.Panic(err)
	}
	log.Println("HAProxy: Configuration updated")
}

func execCommand(cmd string) {
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Exec cmd: %s \n", cmd)
}
