package main

import (
	"flag"
	"time"
	"log"
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
				execCommand(conf.HAProxy.ReloadCommand)
			case _ = <-domainCh:
				log.Println("Domain mapping: Stated changed")
				handleHAPUpdate(conf, marathonConn)
				execCommand(conf.HAProxy.ReloadCommand)
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

func execCommand(command string) {
	out, err := exec.Command(command).Output()
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("executed %s \n", out)
}

func createAndListen(conf configuration.Zookeeper) (chan zk.Event, *zk.Conn) {
	conn, _, err := zk.Connect(conf.ConnectionString(), time.Second * 10)

	if err != nil { log.Panic(err) }

	ch, _ := qzk.ListenToConn(conn, conf.Path, true)

	return ch, conn
}
