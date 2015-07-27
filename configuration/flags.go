package configuration

import "flag"

var (
	marathonEndpoint       string
	bambooEndpoint         string
	bambooZkHost           string
	bambooZkPath           string
	bambooZkReportingDelay int64
	haproxyTemplatePath    string
	haproxyOutputPath      string
	haproxyReloadCommand   string
	statsdEnabled          bool
	statsdHost             string
	statsdPrefix           string
)

const (
	defaultMarathonEndpoint       = "http://192.168.3.4:8080"
	defaultBambooEndpoint         = ""
	defaultBambooZkHost           = "192.168.101.2:2181,192.168.101.3:2181,192.168.101.4:2181"
	defaultBambooZkPath           = "/bamboo"
	defaultBambooZkReportingDelay = int64(5)
	defaultHaproxyTemplatePath    = "config/haproxy_template.cfg"
	defaultHaproxyOutputPath      = "/etc/haproxy/haproxy.cfg"
	defaultHaproxyReloadCommand   = "PIDS=`pidof haproxy`; haproxy -f /etc/haproxy/haproxy.cfg -p /var/run/haproxy.pid -sf $PIDS && while ps -p $PIDS; do sleep 0.2; done"
	defaultStatsdEnabled          = false
	defaultStatsdHost             = "localhost:8125"
	defaultStatsdPrefix           = "bamboo-server.production."
)

func init() {
	flag.StringVar(
		&marathonEndpoint,
		"marathon_endpoint",
		defaultMarathonEndpoint,
		"The endpoint of Marathon.")
	flag.StringVar(
		&bambooEndpoint,
		"bamboo_endpoint",
		defaultBambooEndpoint,
		"The endpoint of Bamboo.")
	flag.StringVar(
		&bambooZkHost,
		"bamboo_zk_host",
		defaultBambooZkHost,
		"The host of Zookeeper for Bamboo.")
	flag.StringVar(
		&bambooZkPath,
		"bamboo_zk_path",
		defaultBambooZkPath,
		"The path of Zookeeper for Bamboo.")
	flag.Int64Var(
		&bambooZkReportingDelay,
		"bamboo_zk_reporting_delay",
		defaultBambooZkReportingDelay,
		"The reporting delay of Zookeeper for Bamboo.")
	flag.StringVar(
		&haproxyTemplatePath,
		"haproxy_template_path",
		defaultHaproxyTemplatePath,
		"The template path of HAProxy configuration file.")
	flag.StringVar(
		&haproxyOutputPath,
		"haproxy_output_path",
		defaultHaproxyOutputPath,
		"The output path of HAProxy configuration file.")
	flag.StringVar(
		&haproxyReloadCommand,
		"haproxy_reload_command",
		defaultHaproxyReloadCommand,
		"The reload command for HAProxy.")
	flag.BoolVar(
		&statsdEnabled,
		"statsd_enabled",
		defaultStatsdEnabled,
		"Start StatsD or not.")
	flag.StringVar(
		&statsdHost,
		"statsd_host",
		defaultStatsdHost,
		"The host of StatsD")
	flag.StringVar(
		&statsdPrefix,
		"statsd_prefix",
		defaultStatsdPrefix,
		"The prefix in StatsD")
}

func setValueFromFlag(field *string, flagVar string) {
	if len(flagVar) > 0 {
		*field = flagVar
	}
}

func setInt64ValueFromFlag(field *int64, flagVar int64) {
	if flagVar > 0 {
		*field = flagVar
	}
}

func setBoolValueFromFlag(field *bool, flagVar bool) {
	if flagVar {
		*field = flagVar
	}
}
