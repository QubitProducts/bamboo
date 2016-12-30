package stats

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type HAProxyStats struct {
	ProxyName                   string `json:"pxname"`
	ServiceName                 string `json:"svname"`
	CurrentQueuedRequests       string `json:"qcur"`
	MaxCurrentQueuedRequests    string `json:"smax"`
	CurrentSessions             string `json:"scur"`
	MaxSessions                 string `json:"smax"`
	ConfiguredSessionsLimit     string `json:"slim"`
	CumulativeConnectionsNumber string `json:"stot"`
	BytesIn                     string `json:"bin"`
	BytesOut                    string `json:"bout"`
	DeniedReuqests              string `json:"dreq"`
	DeniedResponses             string `json:"dresp"`
	ErrorsRequests              string `json:"ereq"`
	ErrorConnect                string `json:"econ"`
	ErrowResponse               string `json:"eresp"`
	RetriedConnectionTimes      string `json:"wretr"`
	RedispatchedAnotherServer   string `json:"wredis"`
	Status                      string `json:"status"`
	Weight                      string `json:"weight"`
	ActiveServerNumber          string `json:"act"`
	BackupServerNumber          string `json:"bck"`
	ChecksFailedNumber          string `json:"chekfail"`
	ChecksDownNumber            string `json:"chkdown"`
	LastChecksSeconeds          string `json:"lastchg"`
	DownTime                    string `json:"downtime"`
	ConfiguredMaxQueue          string `json:"qlimit"`
	ProcessID                   string `json:"pid"`
	UniqueProxyID               string `json:"iid"`
	UniqueInsideProxyID         string `json:"sid"`
	Throttle                    string `json:"throttle"`
	ServerSelectedTimes         string `json:"lbtot"`
	TrackedServerOrProxyID      string `json:"tracked"`
	Type                        string `json:"type"`
	SessionsRate                string `json:"rate"`
	SessionRateLimit            string `json:"rate_lim"`
	SessionsRateMax             string `json:"rate_max"`
	CheckStatus                 string `json:"check_status"`
	CheckCode                   string `json:"check_code"`
	CheckDuration               string `json:"check_duration"`
	HTTPResponses_1xx           string `json:"hrsp_1xx"`
	HTTPResponses_2xx           string `json:"hrsp_2xx"`
	HTTPResponses_3xx           string `json:"hrsp_3xx"`
	HTTPResponses_4xx           string `json:"hrsp_4xx"`
	HTTPResponses_5xx           string `json:"hrsp_5xx"`
	HTTPResponses_other         string `json:"hrsp_other"`
	HealthFailed                string `json:"hanafail"`
	RequestsRate                string `json:"req_rate"`
	RequestsRateMax             string `json:"req_rate_max"`
	RequestsTotal               string `json:"req_tot"`
	CliendAborted               string `json:"cli_abrt"`
	ServerAborted               string `json:"srv_abrt"`
	CompressorResponseBytesIn   string `json:"com_in"`
	CompressorResponseBytesOut  string `json:"comp_out"`
	CompressorBytesBypassed     string `json:"comp_byp"`
	CompressorResponse          string `json:"comp_rsp"`
	LastSession                 string `json:"lastsess"`
	LastCheckHealth             string `json:"last_chk"`
	LastAgent                   string `json:"last_agt"`
	QueueTime                   string `json:"qtime"`
	ConnectTime                 string `json:"ctime"`
	ResponseTime                string `json:"rtime"`
	TotalTime                   string `json:"ttime"`
}

type Config struct {
	HAProxyStatsPort     int
	HAProxyStatsHost     string
	HAProxyStatsUsername string
	HAProxyStatsPassword string
	HAProxyStatsEndpoint string
	HAProxyHost          string
	HAProxyPort          int
	Interval             int
	FalconAgentPort      int
	FalconAgentHost      string
	HAProxyStatsSock     string
	ShouldUploadMetric   []string
}

var (
	config *Config
)

func getEnv(name, defaultValue string) string {
	res := os.Getenv(name)
	if len(res) == 0 {
		res = defaultValue
	}
	return res
}

func initConfig() {
	HAProxyStatsPortEnv := getEnv("HAPROXY_STATS_PORT", "9000")
	HAProxyStatsPort, err := strconv.Atoi(HAProxyStatsPortEnv)
	if err != nil {
		log.Fatalf("get wrong HAProxy stats port: %s\n", HAProxyStatsPortEnv)
	}
	IntervalEnv := getEnv("INTERVAL", "10")
	Interval, err := strconv.Atoi(IntervalEnv)
	if err != nil {
		log.Fatalf("get wrong Interval: %s\n", IntervalEnv)
	}
	FalconAgentPortEnv := getEnv("FALCON_AGENT_PORT", "1988")
	FalconAgentPort, err := strconv.Atoi(FalconAgentPortEnv)
	if err != nil {
		log.Fatalf("get wrong FalconAgentPort: %s\n", FalconAgentPortEnv)
	}
	ShouldUploadMetricEnv := getEnv("METRIC", "hrsp_1xx,hrsp_2xx,hrsp_3xx,hrsp_4xx,hrsp_5xx,bin,bout")
	ShouldUploadMetric := strings.Split(ShouldUploadMetricEnv, ",")

	HAProxyPortEnv := getEnv("HAPROXY_PORT", "20010")
	HAProxyPort, err := strconv.Atoi(HAProxyPortEnv)
	if err != nil {
		log.Fatalf("get wrong Haproxy port: %s\n", HAProxyPort)
	}
	config = &Config{
		HAProxyStatsPort:     HAProxyStatsPort,
		HAProxyStatsHost:     getEnv("HAPROXY_STATS_HOST", "127.0.0.1"),
		HAProxyStatsUsername: getEnv("HAPROXY_STATS_USERNAME", ""),
		HAProxyStatsPassword: getEnv("HAPROXY_STATS_PASSWORD", ""),
		HAProxyStatsEndpoint: getEnv("HAPROXY_STATS_ENDPOINT", ""),
		Interval:             Interval,
		FalconAgentHost:      getEnv("FALCON_AGENT_HOST", "127.0.0.1"),
		FalconAgentPort:      FalconAgentPort,
		ShouldUploadMetric:   ShouldUploadMetric,
		HAProxyHost:          getEnv("HAPROXY_HOST", "127.0.0.1"),
		HAProxyPort:          HAProxyPort,
	}
	log.Println("config: ", config)
}

func main() {
	Stats()

}

func Stats() {
	initConfig()
	Interval := time.Duration(config.Interval) * time.Second
	t := time.NewTicker(Interval)
	for {
		now := <-t.C
		timestamp := now.Unix()
		go func() {
			haproxy_stats, err := getStatsByHTTP()
			if err != nil {
				log.Printf("get haproxy stats failed. %v", err)
			}
			for _, v := range haproxy_stats {
				if _, ok := v["pxname"]; !ok {
					continue
				}
				if _, ok := v["svname"]; !ok {
					continue
				}
				for _, metric := range config.ShouldUploadMetric {
					value, ok := v[metric]
					if !ok {
						log.Printf("metric: %s not exist.\n", metric)
					}
					tags := fmt.Sprintf("pxname=%s,svname=%s", v["pxname"], v["svname"])
					endpoint := fmt.Sprintf("%s-%d", config.HAProxyHost, config.HAProxyPort)
					go func(v, ts, m, t string) {
						err := pushToFalconAgent(
							v,
							ts,
							m,
							t,
							"GAUGE",
							endpoint,
						)
						if err != nil {
							log.Println("push to agent failed.", err)
						}
					}(value, fmt.Sprintf("%d", timestamp), metric, tags)

				}
			}
		}()

	}

}

func getStatsByHTTP() ([]map[string]string, error) {
	uri := fmt.Sprintf("%s:%d", config.HAProxyStatsHost, config.HAProxyStatsPort)
	if !strings.HasPrefix(uri, "http") {
		uri = "http://" + uri
	}
	endpoint := config.HAProxyStatsEndpoint
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	uri += endpoint
	req, err := http.NewRequest("GET", uri+";csv", nil)
	if len(config.HAProxyStatsUsername) > 0 && len(config.HAProxyStatsPassword) > 0 {
		req.SetBasicAuth(config.HAProxyStatsUsername, config.HAProxyStatsPassword)
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) <= 1 {
		return nil, nil
	}
	names := records[0]
	var haproxy_stats_slice []map[string]string
	for _, v := range records[1:] {
		temp := make(map[string]string)
		for index, i := range v {
			name := names[index]
			if index == 0 {
				name = name[2:]
			}
			temp[name] = i
		}
		haproxy_stats_slice = append(haproxy_stats_slice, temp)

	}
	return haproxy_stats_slice, nil
}

func pushToFalconAgent(value, timestamp, metric, tags, counterType,
	endpoint string) error {
	postThing := `[{"metric": "` + metric + `", "endpoint": "haproxy-` +
		endpoint + `", "timestamp": ` + timestamp + `,"step": ` + fmt.Sprintf("%d", config.Interval) + `,"value": ` + value + `,"counterType": "` + counterType + `","tags": "` + tags + `"}]`
	//push data to falcon-agent
	url := fmt.Sprintf("http://%s:%d/v1/push", config.FalconAgentHost, config.FalconAgentPort)
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(postThing))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
