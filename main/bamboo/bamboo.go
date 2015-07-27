package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/kardianos/osext"
	"github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/natefinch/lumberjack"
	"github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/samuel/go-zookeeper/zk"
	"github.com/QubitProducts/bamboo/api"
	"github.com/QubitProducts/bamboo/configuration"
	"github.com/QubitProducts/bamboo/qzk"
	"github.com/QubitProducts/bamboo/services/event_bus"
)

/*
	Commandline arguments
*/
var haproxyCheck bool
var configFromFlags bool
var configFilePath string
var logPath string
var serverBindPort string

func init() {
	flag.BoolVar(&haproxyCheck, "haproxy_check", false, "Check the process of HAProxy ")
	flag.BoolVar(&configFromFlags, "config_from_flags", false, "Read configuration from flags")
	flag.StringVar(&configFilePath, "config", "config/development.json", "Full path of the configuration JSON file")
	flag.StringVar(&logPath, "log", "", "Log path to a file. Default logs to stdout")
	flag.StringVar(&serverBindPort, "bind", ":8000", "Bind HTTP server to a specific port")
}

func main() {
	flag.Parse()
	configureLog()

	// Check HAProxy
	if haproxyCheck {
		validHAProxy, err := checkHAProxy()
		if err != nil {
			log.Fatal("A error occur when HAProxy checking: ", err)
		}
		if !validHAProxy {
			var hintMsg = `Not found the process of HAProxy. 
Please install & run HAProxy before running Bamboo. Look at the following tips:
- Install HAProxy:
    apt-get install -yqq software-properties-common && \
    apt-add-repository ppa:vbernat/haproxy-1.5 && \
    apt-get update -yqq && \
    apt-get install -yqq haproxy
- Append the following content to /etc/haproxy/haproxy.cfg (optional):
  # Begin
  listen stats :1936
      mode http
      stats enable
      stats hide-version
      stats realm Haproxy\ Statistics
      stats uri /
      stats auth Username:Password
  # End
- Run HAProxy:
    haproxy -f /etc/haproxy/haproxy.cfg -p /var/run/haproxy.pid -D
		`
			log.Fatal("Error: ", hintMsg)
		}
	}

	// Load configuration
	var conf configuration.Configuration
	var err error
	if configFromFlags {
		conf, err = configuration.FromFlags()
	} else {
		conf, err = configuration.FromFile(configFilePath)
	}
	if err != nil {
		log.Fatal("A error occur when load configuration: ", err)
	}

	eventBus := event_bus.New()

	// Wait for died children to avoid zombies
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGCHLD)
	go func() {
		for {
			sig := <-signalChannel
			if sig == syscall.SIGCHLD {
				r := syscall.Rusage{}
				syscall.Wait4(-1, nil, 0, &r)
			}
		}
	}()

	// Create StatsD client
	conf.StatsD.CreateClient()

	// Create Zookeeper connection
	zkConn := listenToZookeeper(conf, eventBus)

	// Register handlers
	handlers := event_bus.Handlers{Conf: &conf, Zookeeper: zkConn}
	eventBus.Register(handlers.MarathonEventHandler)
	eventBus.Register(handlers.ServiceEventHandler)
	eventBus.Publish(event_bus.MarathonEvent{EventType: "bamboo_startup", Timestamp: time.Now().Format(time.RFC3339)})
<<<<<<< HEAD

	// Handle gracefully exit
	registerOSSignals()
=======
>>>>>>> Improve the method 'Plaintext of ''MarathonEvent'.

	// Start server
	initServer(&conf, zkConn, eventBus)
}

func checkHAProxy() (bool, error) {
	var result bool
	output, err := exec.Command("pidof", "haproxy").Output()
	if err != nil {
		return result, err
	}
	pid := string(output)
	if len(pid) > 0 {
		result = true
	}
	return result, err
}

func initServer(conf *configuration.Configuration, conn *zk.Conn, eventBus *event_bus.EventBus) {
	stateAPI := api.StateAPI{Config: conf, Zookeeper: conn}
	serviceAPI := api.ServiceAPI{Config: conf, Zookeeper: conn}
	eventSubAPI := api.EventSubscriptionAPI{Conf: conf, EventBus: eventBus}

	conf.StatsD.Increment(1.0, "restart", 1)
	// Status live information
	router := martini.Classic()
	router.Get("/status", api.HandleStatus)

	// API
	router.Group("/api", func(api martini.Router) {
		// State API
		api.Get("/state", stateAPI.Get)
		// Service API
		api.Get("/services", serviceAPI.All)
		api.Post("/services", serviceAPI.Create)
		api.Put("/services/**", serviceAPI.Put)
		api.Delete("/services/**", serviceAPI.Delete)
		api.Post("/marathon/event_callback", eventSubAPI.Callback)
	})

	// Static pages
	router.Use(martini.Static(path.Join(executableFolder(), "webapp")))

	registerMarathonEvent(conf)
	router.RunOnAddr(serverBindPort)
}

// Get current executable folder path
func executableFolder() string {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	return folderPath
}

func registerMarathonEvent(conf *configuration.Configuration) {

	client := &http.Client{}
	// it's safe to register with multiple marathon nodes
	for _, marathon := range conf.Marathon.Endpoints() {
		url := marathon + "/v2/eventSubscriptions?callbackUrl=" + conf.Bamboo.Endpoint + "/api/marathon/event_callback"
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			errorMsg := "An error occurred while accessing Marathon callback system: %s\n"
			log.Printf(errorMsg, err)
			return
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
		body := string(bodyBytes)
		if strings.HasPrefix(body, "{\"message") {
			warningMsg := "Access to the callback system of Marathon seems to be failed, response: %s\n"
			log.Printf(warningMsg, body)
		}
	}
}

func createAndListen(conf configuration.Zookeeper) (chan zk.Event, *zk.Conn) {
	conn, _, err := zk.Connect(conf.ConnectionString(), time.Second*10)

	if err != nil {
		log.Panic(err)
	}

	ch, _ := qzk.ListenToConn(conn, conf.Path, true, conf.Delay())
	return ch, conn
}

func listenToZookeeper(conf configuration.Configuration, eventBus *event_bus.EventBus) *zk.Conn {
	serviceCh, serviceConn := createAndListen(conf.Bamboo.Zookeeper)

	go func() {
		for {
			select {
			case _ = <-serviceCh:
				eventBus.Publish(event_bus.ServiceEvent{EventType: "change"})
			}
		}
	}()
	return serviceConn
}

func configureLog() {
	if len(logPath) > 0 {
		log.SetOutput(io.MultiWriter(&lumberjack.Logger{
			Filename: logPath,
			// megabytes
			MaxSize:    100,
			MaxBackups: 3,
			//days
			MaxAge: 28,
		}, os.Stdout))
	}
}

func registerOSSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Println("Server Stopped")
			os.Exit(0)
		}
	}()
}
