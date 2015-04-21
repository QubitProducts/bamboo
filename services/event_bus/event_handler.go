package event_bus

import (
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/samuel/go-zookeeper/zk"
	"github.com/QubitProducts/bamboo/configuration"
	"github.com/QubitProducts/bamboo/services/haproxy"
	"github.com/QubitProducts/bamboo/services/template"
)

type Handlers struct {
	Conf      *configuration.Configuration
	Zookeeper *zk.Conn
}

func (h *Handlers) MarathonEventHandler(event MarathonEvent) {
	log.Printf("Marathon event(type=%s, timestamp=%s): %s\n",
		event.EventType, event.Timestamp, event.Plaintext())
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
		loop := true
		monitor := func() (abnormal bool) {
			defer func() {
				if err := recover(); err != nil {
					log.Fatalf("An fatal error occurred while updating haproxy config file: %s\n", err)
					abnormal = true
				}
			}()
			for {
				h := <-updateChan
				handleHAPUpdate(h.Conf, h.Zookeeper)
			}
		}
		for loop {
			loop = monitor()
			if loop {
				log.Println("Restarting update loop")
			}
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

func handleHAPUpdate(conf *configuration.Configuration, conn *zk.Conn) bool {
	log.Println("Updating HAProxy configuration file...")
	currentContent, _ := ioutil.ReadFile(conf.HAProxy.OutputPath)

	templateContent, err := ioutil.ReadFile(conf.HAProxy.TemplatePath)
	if err != nil {
		log.Fatalf("Ignore to updating because cannot read template file: %s", err)
		return false
	}

	templateData := haproxy.GetTemplateData(conf, conn)

	newContent, err := template.RenderTemplate(conf.HAProxy.TemplatePath, string(templateContent), templateData)

	if err != nil {
		log.Fatalf("Ignore to updating because template syntax error: \n %s", err)
		return false
	}

	if currentContent == nil || string(currentContent) != newContent {
		err := ioutil.WriteFile(conf.HAProxy.OutputPath, []byte(newContent), 0666)
		if err != nil {
			log.Fatalf("Failed to write template on path: %s", err)
		}

		err = execCommand(conf.HAProxy.ReloadCommand)
		if err != nil {
			log.Fatalf("HAProxy: update failed\n")
		} else {
			conf.StatsD.Increment(1.0, "reload.marathon", 1)
			log.Println("HAProxy: Configuration updated")
		}
		return true
	} else {
		log.Println("HAProxy: Same content, no need to reload")
		return false
	}
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
