package marathon

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/QubitProducts/bamboo/configuration"
)

const (
	taskStateRunning                           = "TASK_RUNNING"
	readinessCheckDefaultTimeout time.Duration = 10 * time.Second
	readinessCheckSafetyMargin   time.Duration = 5 * time.Second
)

type readinessCalculator struct {
	checkDefaultTimeout time.Duration
	checkSafetyMargin   time.Duration
}

var readyCalculator = readinessCalculator{
	checkDefaultTimeout: readinessCheckDefaultTimeout,
	checkSafetyMargin:   readinessCheckSafetyMargin,
}

// Describes an app process running
type Task struct {
	Id    string
	Host  string
	Port  int
	Ports []int
	Alive bool
	State string
	Ready bool
}

// A health check on the application
type HealthCheck struct {
	// One of TCP, HTTP or COMMAND
	Protocol string
	// The path (if Protocol is HTTP)
	Path string
	// The position of the port targeted in the ports array
	PortIndex int
}

// An app may have multiple processes
type App struct {
	Id                  string
	MesosDnsId          string
	EscapedId           string
	HealthCheckPath     string
	HealthCheckProtocol string
	HealthChecks        []HealthCheck
	ReadinessCheckPath  string
	Tasks               []Task
	ServicePort         int
	ServicePorts        []int
	Env                 map[string]string
	Labels              map[string]string
	SplitId             []string
	IpAddress           AppIpAddress `json:"ipAddress"`
}

type AppIpAddress struct {
	Discovery Discovery `json:"discovery"`
}

type Discovery struct {
	Ports []Port `json:"ports"`
}

type Port struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
}

type TaskIpAddress struct {
	IpAddress string `json:"ipAddress"`
	Protocol  string `json:"protocol"`
}

type AppList []App

func (slice AppList) Len() int {
	return len(slice)
}

func (slice AppList) Less(i, j int) bool {
	return slice[i].Id < slice[j].Id
}

func (slice AppList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type HealthCheckResult struct {
	Alive bool
}

type marathonTask struct {
	AppId              string
	Id                 string
	Host               string
	Ports              []int
	ServicePorts       []int
	State              string
	StartedAt          string
	StagedAt           string
	Version            string
	IpAddresses        []TaskIpAddress `json:"IpAddresses"`
	HealthCheckResults []HealthCheckResult
}

type marathonTaskList []marathonTask

func (slice marathonTaskList) Len() int {
	return len(slice)
}

func (slice marathonTaskList) Less(i, j int) bool {
	return slice[i].Id < slice[j].Id
}

func (slice marathonTaskList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type marathonApps struct {
	Apps []marathonApp `json:"apps"`
}

type marathonApp struct {
	Id                    string                   `json:"id"`
	HealthChecks          []marathonHealthCheck    `json:"healthChecks"`
	Ports                 []int                    `json:"ports"`
	Env                   map[string]string        `json:"env"`
	Labels                map[string]string        `json:"labels"`
	Deployments           []deployment             `json:"deployments"`
	Tasks                 marathonTaskList         `json:"tasks"`
	ReadinessChecks       []marathonReadinessCheck `json:"readinessChecks"`
	ReadinessCheckResults []readinessCheckResult   `json:"readinessCheckResults"`
	IpAddress             AppIpAddress             `json:"ipAddress"`
}

type marathonHealthCheck struct {
	Path      string `json:"path"`
	Protocol  string `json:"protocol"`
	PortIndex int    `json:"portIndex"`
}

type marathonReadinessCheck struct {
	Path           string `json:"path"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
}

type deployment struct {
	ID string `json:"id"`
}

type readinessCheckResult struct {
	TaskID string `json:"taskId"`
	Ready  bool   `json:"ready"`
}

/*
	Apps returns a struct that describes Marathon current app and their
	sub tasks information.

	Parameters:
		endpoint: Marathon HTTP endpoint, e.g. http://localhost:8080
*/
func FetchApps(maraconf configuration.Marathon, conf *configuration.Configuration) (AppList, error) {
	var marathonApps []marathonApp
	var err error

	// Try all configured endpoints until one succeeds or we exhaust the list,
	// whichever comes first.
	for _, url := range maraconf.Endpoints() {
		marathonApps, err = fetchMarathonApps(url, conf)
		if err == nil {
			for _, marathonApp := range marathonApps {
				sort.Sort(marathonApp.Tasks)
			}
			apps := createApps(marathonApps)
			sort.Sort(apps)
			return apps, nil
		}
	}
	// return last error
	return nil, err
}

func fetchMarathonApps(endpoint string, conf *configuration.Configuration) ([]marathonApp, error) {
	var appResponse marathonApps
	if err := parseJSON(endpoint+"/v2/apps?embed=app.tasks&embed=app.deployments&embed=app.readiness", conf, &appResponse); err != nil {
		return nil, err
	}

	return appResponse.Apps, nil
}

func parseJSON(url string, conf *configuration.Configuration, out interface{}) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if len(conf.Marathon.User) > 0 && len(conf.Marathon.Password) > 0 {
		req.SetBasicAuth(conf.Marathon.User, conf.Marathon.Password)
	}

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &out)
	if err != nil {
		return err
	}

	return nil
}

func createApps(marathonApps []marathonApp) AppList {
	apps := AppList{}

	for _, mApp := range marathonApps {
		appId := mApp.Id
		// Try to handle old app id format without slashes
		appPath := "/" + strings.TrimPrefix(mApp.Id, "/")

		// build App from marathonApp
		app := App{
			Id:                  appPath,
			MesosDnsId:          getMesosDnsId(appPath),
			EscapedId:           strings.Replace(appId, "/", "::", -1),
			HealthCheckPath:     parseHealthCheckPath(mApp.HealthChecks),
			HealthCheckProtocol: parseHealthCheckProtocol(mApp.HealthChecks),
			ReadinessCheckPath:  parseReadinessCheckPath(mApp.ReadinessChecks),
			Env:                 mApp.Env,
			Labels:              mApp.Labels,
			SplitId:             strings.Split(appId, "/"),
			IpAddress:           mApp.IpAddress,
		}

		app.HealthChecks = make([]HealthCheck, 0, len(mApp.HealthChecks))
		for _, marathonCheck := range mApp.HealthChecks {
			check := HealthCheck{
				Protocol:  marathonCheck.Protocol,
				Path:      marathonCheck.Path,
				PortIndex: marathonCheck.PortIndex,
			}
			app.HealthChecks = append(app.HealthChecks, check)
		}

		if len(mApp.Ports) > 0 {
			app.ServicePort = mApp.Ports[0]
			app.ServicePorts = mApp.Ports
		}

		// build Tasks for this App
		tasks := []Task{}
		for _, mTask := range mApp.Tasks {
			var host string
			var port int
			if len(mTask.Ports) > 0 {
				host = mTask.Host
				port = mTask.Ports[0]
			}
			if len(app.IpAddress.Discovery.Ports) > 0 {
				if len(mTask.IpAddresses) > 0 {
					host = mTask.IpAddresses[0].IpAddress
					port = app.IpAddress.Discovery.Ports[0].Number
				}
			}

			if host != "" && port != 0 {
				t := Task{
					Id:    mTask.Id,
					Host:  host,
					Port:  port,
					Ports: mTask.Ports,
					Alive: calculateTaskHealth(mTask.HealthCheckResults, mApp.HealthChecks),
					Ready: readyCalculator.calculate(mTask, mApp),
					State: mTask.State,
				}
				tasks = append(tasks, t)
			}
		}
		app.Tasks = tasks

		apps = append(apps, app)
	}
	return apps
}

func getMesosDnsId(appPath string) string {
	// split up groups and recombine for how mesos-dns/consul/etc use service name
	//   "/nested/group/app" -> "app-group-nested"
	groups := strings.Split(appPath, "/")
	reverseGroups := []string{}
	for i := len(groups) - 1; i >= 0; i-- {
		if groups[i] != "" {
			reverseGroups = append(reverseGroups, groups[i])
		}
	}
	return strings.Join(reverseGroups, "-")
}

func parseHealthCheckPath(checks []marathonHealthCheck) string {
	for _, check := range checks {
		if check.Protocol != "HTTP" && check.Protocol != "HTTPS" {
			continue
		}
		return check.Path
	}
	return ""
}

/* maybe combine this with the above? */
func parseHealthCheckProtocol(checks []marathonHealthCheck) string {
	for _, check := range checks {
		if check.Protocol != "HTTP" && check.Protocol != "HTTPS" {
			continue
		}
		return check.Protocol
	}
	return ""
}

func parseReadinessCheckPath(checks []marathonReadinessCheck) string {
	if len(checks) > 0 {
		return checks[0].Path
	}

	return ""
}

func calculateTaskHealth(healthCheckResults []HealthCheckResult, healthChecks []marathonHealthCheck) bool {
	// If we don't even have health check results for every health check, don't
	// count the task as healthy.
	if len(healthChecks) > len(healthCheckResults) {
		return false
	}
	for _, healthCheck := range healthCheckResults {
		if !healthCheck.Alive {
			return false
		}
	}
	return true
}

func (rc *readinessCalculator) calculate(task marathonTask, maraApp marathonApp) bool {
	switch {
	case task.State != taskStateRunning:
		// By definition, a task not running cannot be ready.
		log.Printf("task %s app %s: ready = false [task state %s != required state %s]", task.Id, maraApp.Id, task.State, taskStateRunning)
		return false

	case len(maraApp.Deployments) == 0:
		// We only care about readiness during deployments; post-deployment readiness
		// should be covered by a separate HAProxy health check definition.
		log.Printf("task %s app %s: ready = true [no deployment ongoing]", task.Id, maraApp.Id)
		return true

	case len(maraApp.ReadinessChecks) == 0:
		// Applications without configured readiness checks are always considered
		// ready.
		log.Printf("task %s app %s: ready = true [no readiness checks on app]", task.Id, maraApp.Id)
		return true
	}

	// Loop through all readiness check results and return the results for
	// matching task IDs.
	for _, readinessCheckResult := range maraApp.ReadinessCheckResults {
		if readinessCheckResult.TaskID == task.Id {
			log.Printf("task %s app %s: ready = %t [evaluating readiness check ready state]", task.Id, maraApp.Id, readinessCheckResult.Ready)
			return readinessCheckResult.Ready
		}
	}

	// There's a corner case sometimes hit where the first new task of a
	// deployment goes from TASK_STAGING to TASK_RUNNING without a corresponding
	// health check result being included in the API response. This only happens
	// in a very short (yet unlucky) time frame and does not repeat for subsequent
	// tasks of the same deployment.
	// Complicating matters, the situation may occur for both initially deploying
	// applications as well as rolling-upgraded ones where one or more tasks from
	// a previous deployment exist already and are joined by new tasks from a
	// subsequent deployment. We must always make sure that pre-existing tasks
	// maintain their ready state while newly launched tasks must be considered
	// unready until a check result appears.
	// We distinguish the two cases by comparing the current time with the start
	// time of the task: It should take Marathon at most one readiness check timeout
	// interval (plus some safety margin to account for the delayed nature of
	// distributed systems) for readiness check results to be returned along the API
	// response. Once the task turns old enough, we assume it to be part of a
	// pre-existing deployment and mark it as ready. Note that it is okay to err
	// on the side of caution and consider a task unready until the safety time
	// window has elapsed because a newly created task should be readiness-checked
	// and be given a result fairly shortly after its creation (i.e., on the scale
	// of seconds).
	readinessCheckTimeoutSecs := maraApp.ReadinessChecks[0].TimeoutSeconds
	readinessCheckTimeout := time.Duration(readinessCheckTimeoutSecs) * time.Second
	if readinessCheckTimeout == 0 {
		log.Printf("task %s app %s: readiness check timeout not set, using default value %s", task.Id, maraApp.Id, rc.checkDefaultTimeout)
		readinessCheckTimeout = rc.checkDefaultTimeout
	} else {
		readinessCheckTimeout += rc.checkSafetyMargin
	}

	startTime, err := time.Parse(time.RFC3339, task.StartedAt)
	if err != nil {
		// An unparseable start time should never occur; if it does, we assume the
		// problem should be surfaced as quickly as possible, which is easiest if
		// we shun the task from rotation.
		log.Printf("task %s app %s: ready = false [task start-time %s not parseable]", task.Id, maraApp.Id, task.StartedAt)
		return false
	}

	since := time.Since(startTime)
	if since < readinessCheckTimeout {
		log.Printf("task %s app %s: ready = false [task with start-time %s not within assumed check timeout window of %s (elapsed time since task start: %s)]", task.Id, maraApp.Id, startTime.Format(time.RFC3339), readinessCheckTimeout, since)
		return false
	}

	// Finally, we can be certain this task is not part of the deployment (i.e.,
	// it's an old task that's going to transition into the TASK_KILLING and/or
	// TASK_KILLED state as new tasks' readiness checks gradually turn green.)
	log.Printf("task %s app %s: ready = true [task with start-time %s not involved in deployment (elapsed time since task start: %s)]", task.Id, maraApp.Id, startTime.Format(time.RFC3339), since)
	return true
}
