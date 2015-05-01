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
	Host string
	Port int
	Alive bool
}

// An app may have multiple processes
type App struct {
	Id              string
	EscapedId       string
	HealthCheckPath string
	Tasks           []Task
	ServicePort     int
	Env             map[string]string
	Labels          map[string]string
	SplitId         []string
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

type MarathonTaskList []MarathonTask

type MarathonTasks struct {
	Tasks MarathonTaskList `json:tasks`
}

type HealthCheckResult struct {
	Alive bool
}

type MarathonTask struct {
	AppId        string
	Id           string
	Host         string
	Ports        []int
	ServicePorts []int
	StartedAt    string
	StagedAt     string
	Version      string
	HealthCheckResults []HealthCheckResult
}

func (slice MarathonTaskList) Len() int {
	return len(slice)
}

func (slice MarathonTaskList) Less(i, j int) bool {
	return slice[i].StagedAt < slice[j].StagedAt
}

func (slice MarathonTaskList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type MarathonApps struct {
	Apps []MarathonApp `json:apps`
}

type MarathonApp struct {
	Id           string            `json:id`
	HealthChecks []HealthChecks    `json:healthChecks`
	Ports        []int             `json:ports`
	Env          map[string]string `json:env`
	Labels       map[string]string `json:labels`
}

type HealthChecks struct {
	Path string `json:path`
}

func fetchMarathonApps(endpoint string) (map[string]MarathonApp, error) {
	response, err := http.Get(endpoint + "/v2/apps")

	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		var appResponse MarathonApps

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(contents, &appResponse)
		if err != nil {
			return nil, err
		}

		dataById := map[string]MarathonApp{}

		for _, appConfig := range appResponse.Apps {
			dataById[appConfig.Id] = appConfig
		}

		return dataById, nil
	}
}

func fetchTasks(endpoint string) (map[string][]MarathonTask, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint+"/v2/tasks", nil)
	req.Header.Add("Accept", "application/json")
	response, err := client.Do(req)

	var tasks MarathonTasks

	if err != nil {
		return nil, err
	} else {
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

		tasksById := map[string][]MarathonTask{}
		for _, task := range taskList {
			if tasksById[task.AppId] == nil {
				tasksById[task.AppId] = []MarathonTask{}
			}
			tasksById[task.AppId] = append(tasksById[task.AppId], task)
		}

		return tasksById, nil
	}
}

func calculateTaskHealth(healthCheckResults []HealthCheckResult, healthChecks []HealthChecks) bool {
	//If we don't even have health check results for every health check, don't count the task as healthy
	if len(healthChecks) > len(healthCheckResults) {
		return false;
	}
	for _, healthCheck := range healthCheckResults {
		if !healthCheck.Alive {
			return false;
		}
	}
	return true;
}

func createApps(tasksById map[string][]MarathonTask, marathonApps map[string]MarathonApp) AppList {

	apps := AppList{}

	for appId, tasks := range tasksById {
		simpleTasks := []Task{}

		for _, task := range tasks {
			if len(task.Ports) > 0 {
				simpleTasks = append(simpleTasks, Task{Host: task.Host, Port: task.Ports[0], Alive: calculateTaskHealth(task.HealthCheckResults, marathonApps[appId].HealthChecks)})
			}
		}

		// Try to handle old app id format without slashes
		appPath := appId
		if !strings.HasPrefix(appId, "/") {
			appPath = "/" + appId
		}

		app := App{
			// Since Marathon 0.7, apps are namespaced with path
			Id: appPath,
			// Used for template
			EscapedId:       strings.Replace(appId, "/", "::", -1),
			Tasks:           simpleTasks,
			HealthCheckPath: parseHealthCheckPath(marathonApps[appId].HealthChecks),
			Env:             marathonApps[appId].Env,
			Labels:          marathonApps[appId].Labels,
			SplitId:         strings.Split(appId, "/"),
		}

		if len(marathonApps[appId].Ports) > 0 {
			app.ServicePort = marathonApps[appId].Ports[0]
		}

		apps = append(apps, app)
	}
	return apps
}

func parseHealthCheckPath(checks []HealthChecks) string {
	if len(checks) > 0 {
		return checks[0].Path
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
