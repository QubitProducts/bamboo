package marathon

import (
	"encoding/json"
	"github.com/QubitProducts/bamboo/configuration"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// Describes an app process running
type Task struct {
	Host  string
	Port  int
	Ports []int
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
	Id              string
	EscapedId       string
	HealthCheckPath string
	HealthChecks    []HealthCheck
	Tasks           []Task
	ServicePort     int
	ServicePorts    []int
	Env             map[string]string
	Labels          map[string]string
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

type marathonTaskList []marathonTask

type marathonTasks struct {
	Tasks marathonTaskList `json:"tasks"`
}

type marathonTask struct {
	AppId        string
	Id           string
	Host         string
	Ports        []int
	ServicePorts []int
	StartedAt    string
	StagedAt     string
	Version      string
}

func (slice marathonTaskList) Len() int {
	return len(slice)
}

func (slice marathonTaskList) Less(i, j int) bool {
	return slice[i].StagedAt < slice[j].StagedAt
}

func (slice marathonTaskList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type marathonApps struct {
	Apps []marathonApp `json:"apps"`
}

type marathonApp struct {
	Id           string                `json:"id"`
	HealthChecks []marathonHealthCheck `json:"healthChecks"`
	Ports        []int                 `json:"ports"`
	Env          map[string]string     `json:"env"`
	Labels       map[string]string     `json:"labels"`
}

type marathonHealthCheck struct {
	Path      string `json:"path"`
	Protocol  string `json:"protocol"`
	PortIndex int    `json:"portIndex"`
}

func fetchMarathonApps(endpoint string) (map[string]marathonApp, error) {
	response, err := http.Get(endpoint + "/v2/apps")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	var appResponse marathonApps

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &appResponse)
	if err != nil {
		return nil, err
	}

	dataById := map[string]marathonApp{}

	for _, appConfig := range appResponse.Apps {
		dataById[appConfig.Id] = appConfig
	}

	return dataById, nil
}

func fetchTasks(endpoint string) (map[string][]marathonTask, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint+"/v2/tasks", nil)
	req.Header.Add("Accept", "application/json")
	response, err := client.Do(req)

	var tasks marathonTasks

	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &tasks)
	if err != nil {
		return nil, err
	}

	taskList := tasks.Tasks
	sort.Sort(taskList)

	tasksById := map[string][]marathonTask{}
	for _, task := range taskList {
		if tasksById[task.AppId] == nil {
			tasksById[task.AppId] = []marathonTask{}
		}
		tasksById[task.AppId] = append(tasksById[task.AppId], task)
	}

	return tasksById, nil
}

func createApps(tasksById map[string][]marathonTask, marathonApps map[string]marathonApp) AppList {

	apps := AppList{}

	for appId, mApp := range marathonApps {

		// Try to handle old app id format without slashes
		appPath := appId
		if !strings.HasPrefix(appId, "/") {
			appPath = "/" + appId
		}

		// build App from marathonApp
		app := App{
			Id:              appPath,
			EscapedId:       strings.Replace(appId, "/", "::", -1),
			HealthCheckPath: parseHealthCheckPath(mApp.HealthChecks),
			Env:             mApp.Env,
			Labels:          mApp.Labels,
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
		for _, mTask := range tasksById[appId] {
			if len(mTask.Ports) > 0 {
				t := Task{
					Host:  mTask.Host,
					Port:  mTask.Ports[0],
					Ports: mTask.Ports,
				}
				tasks = append(tasks, t)
			}
		}
		app.Tasks = tasks

		apps = append(apps, app)
	}
	return apps
}

func parseHealthCheckPath(checks []marathonHealthCheck) string {
	for _, check := range checks {
		if check.Protocol != "HTTP" {
			continue
		}
		return check.Path
	}
	return ""
}

/*
	Apps returns a struct that describes Marathon current app and their
	sub tasks information.

	Parameters:
		endpoint: Marathon HTTP endpoint, e.g. http://localhost:8080
*/
func FetchApps(maraconf configuration.Marathon) (AppList, error) {

	var applist AppList
	var err error

	// try all configured endpoints until one succeeds
	for _, url := range maraconf.Endpoints() {
		applist, err = _fetchApps(url)
		if err == nil {
			return applist, err
		}
	}
	// return last error
	return nil, err
}

func _fetchApps(url string) (AppList, error) {
	tasks, err := fetchTasks(url)
	if err != nil {
		return nil, err
	}

	marathonApps, err := fetchMarathonApps(url)
	if err != nil {
		return nil, err
	}

	apps := createApps(tasks, marathonApps)
	sort.Sort(apps)
	return apps, nil
}
