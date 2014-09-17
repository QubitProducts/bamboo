package marathon

import(
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
)

// Describes an app process running
type Task struct {
	Host string
	Port int
}

// An app may have multiple processes
type App struct {
	Id string
	EscapedId string
	HealthCheckPath string
	Tasks []Task
}


type MarathonTasks struct {
	Tasks []MarathonTask `json:tasks`
}

type MarathonTask struct {
	AppId string
	Id    string
	Host  string
	Ports []int
	startedAt string
	stagedAt  string
	version   string
}

type MarathonApps struct {
	Apps []MarathonApp `json:apps`
}

type MarathonApp struct {
	Id string `json:id`
	HealthChecks []HealthChecks `json:healthChecks`
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
		if (err != nil) {
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
	req, err := http.NewRequest("GET", endpoint + "/v2/tasks", nil)
	req.Header.Add("Accept", "application/json")
	response, err := client.Do(req)

	var tasks MarathonTasks

	if err != nil {
		return nil, err
	} else {
		contents, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil { return nil, err }

		err = json.Unmarshal(contents, &tasks)
		if err != nil { return nil, err }

		tasksById := map[string][]MarathonTask{}

		for _, task := range tasks.Tasks {
			if tasksById[task.AppId] == nil {
				tasksById[task.AppId] = []MarathonTask{}
			}
			tasksById[task.AppId] = append(tasksById[task.AppId], task)
		}

		return tasksById, nil
	}
}

func createApps(tasksById map[string][]MarathonTask, marathonApps map[string]MarathonApp) []App {

	apps := []App{}

	for appId, tasks  := range tasksById {
			simpleTasks := []Task{}

			for _, task := range tasks {
				simpleTasks = append(simpleTasks, Task{ Host: task.Host, Port: task.Ports[0] })
			}

			app := App {
				// Since Marathon 0.7, apps are namespaced with path
				Id: appId,
				// Used for template
				EscapedId: strings.Replace(appId, "/", "::", -1),
				Tasks: simpleTasks,
				HealthCheckPath: parseHealthCheckPath(marathonApps[appId].HealthChecks),
			}
			apps = append(apps, app)
	}
	return apps
}

func parseHealthCheckPath(checks []HealthChecks) string {
	if (len(checks) > 0) {
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
func FetchApps(endpoint string) ([]App, error) {
	tasks, err := fetchTasks(endpoint)
	if err != nil { return nil, err }

	marathonApps, err := fetchMarathonApps(endpoint)
	if err != nil { return nil, err }

	apps := createApps(tasks, marathonApps)
	return apps, nil
}
