package main

import (
	"flag"
	"time"
	"log"
	"sync"
	"strings"
	"os/exec"
	"net/http"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/zenazn/goji"

	"bamboo/api"
	"bamboo/configuration"
	"bamboo/qzk"
	"bamboo/services/haproxy"
)




/*
	Commandline arguments
 */
var configFilePath string
func init() {
	flag.StringVar(&configFilePath, "config", "config/development.json", "Full path of the configuration JSON file")
}


func main() {
	flag.Parse()
	conf, err := configuration.FromFile(configFilePath)
	if err != nil { log.Fatal(err) }

	conns := listenToZookeeper(conf)

	initServer(conf, conns)
}

func initServer(conf configuration.Configuration, conns Conns) {
	stateAPI := api.State{Config: conf, Zookeeper: conns.DomainMapping}
	domainAPI := api.Domain{Config: conf, Zookeeper: conns.DomainMapping}

	// Status live information
	goji.Get("/status", api.HandleStatus)

	// State API
	goji.Get("/api/state", stateAPI.Get)

	// Domains API
	goji.Get("/api/state/domains", domainAPI.All)
	goji.Post("/api/state/domains", domainAPI.Create)
	goji.Delete("/api/state/domains/:id", domainAPI.Delete)
	goji.Put("/api/state/domains/:id", domainAPI.Put)

	// Static pages
	goji.Get("/*", http.FileServer(http.Dir("./webapp")))

	goji.Serve()
}

type Conns struct {
	Marathon      *zk.Conn
	DomainMapping *zk.Conn
}

func listenToZookeeper(conf configuration.Configuration) Conns {
	marathonCh, marathonConn := createAndListen(conf.Marathon.Zookeeper)
	domainCh, domainConn := createAndListen(conf.DomainMapping.Zookeeper)

	go func() {
		for {
			select {
			case _ = <-marathonCh:
				log.Println("Marathon: State changed")
				handleHAPUpdate(conf, marathonConn)
				wg := new(sync.WaitGroup)
				wg.Add(1)
				execCommand(conf.HAProxy.ReloadCommand, wg)
			case _ = <-domainCh:
				log.Println("Domain mapping: Stated changed")
				handleHAPUpdate(conf, marathonConn)
				wg := new(sync.WaitGroup)
				wg.Add(1)
				execCommand(conf.HAProxy.ReloadCommand, wg)
			}
		}
	}()

	return Conns{marathonConn, domainConn}
}

func handleHAPUpdate(conf configuration.Configuration, conn * zk.Conn) {
	err := haproxy.WriteHAProxyConfig(conf.HAProxy, haproxy.GetTemplateData(conf, conn))
	if err != nil {
		log.Panic(err)
	}
	log.Println("HAProxy: Configuration updated")
}

func execCommand(cmd string, wg *sync.WaitGroup) {
	log.Println("Exec command: ", cmd)
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head,parts...).Output()
	if err != nil {
		log.Printf("%s \n", err)
	}
	log.Printf("%s \n", out)
	wg.Done()
}

func createAndListen(conf configuration.Zookeeper) (chan zk.Event, *zk.Conn) {
	conn, _, err := zk.Connect(conf.ConnectionString(), time.Second * 10)

	if err != nil { log.Panic(err) }

	ch, _ := qzk.ListenToConn(conn, conf.Path, true)

	return ch, conn
}
