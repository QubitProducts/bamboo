package marathon

import(
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
)

// Describes an app process running
type Task struct {
	Host string
	Port string
}

// An app may have multiple processes
type App struct {
	Id string
	Port string
	HealthCheckPath string
	Tasks []Task
}

type AppConfigResponse struct {
	Apps []AppConfiguration `json:apps`
}

type AppConfiguration struct {
	Id string `json:id`
	HealthChecks []HealthChecks `json:healthChecks`
}

type HealthChecks struct {
	Path string `json:path`
}

func fetchAppConfiguration(endpoint string) (map[string]AppConfiguration, error) {
	response, err := http.Get(endpoint + "/v2/apps")

	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		var appResponse AppConfigResponse

		contents, err := ioutil.ReadAll((response.Body))
		if (err != nil) {
			return nil, err
		}

		err = json.Unmarshal(contents, &appResponse)
		if err != nil {
			return nil, err
		}

		dataById := map[string]AppConfiguration{}

		for _, appConfig := range appResponse.Apps {
			dataById[appConfig.Id] = appConfig
		}

		return dataById, nil
	}
}

func fetchTasks(endpoint string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint + "/v2/tasks", nil)
	req.Header.Add("Accept", "text/plain")
	response, err := client.Do(req)

	if err != nil {
		return "", err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		return string(contents), nil
	}
}

func parseApps(contents string, appConfiguration map[string]AppConfiguration) []App {
	lines := strings.Split(contents, "\n")
	apps := []App{}
	for _, line := range lines {
		if len(line) > 0 {
			appId, appPort, tasks := parseTasks(line)

			app := App {
				Id: appId,
				Port: appPort,
				Tasks: tasks,
				HealthCheckPath: parseHealthCheckPath(appConfiguration[appId].HealthChecks),
			}
			apps = append(apps, app)
		}
	}
	return apps
}

func parseHealthCheckPath(checks []HealthChecks) string {
	if (len(checks) > 0) {
		return checks[0].Path
	}
	return ""
}

func parseTasks(line string) (appId string, appPort string, tasks []Task)  {
	columns := strings.Split(line, "\t")
	appId = columns[0]
	appPort = columns[1]
	tasks = []Task{}

	for _, process := range columns[2:] {
		if len(process) > 0 {
			values := strings.Split(process, ":")
			tasks = append(tasks, Task { Host: values[0], Port: values[1] })
		}
	}

	return appId, appPort, tasks
}

/*
	Apps returns a struct that describes Marathon current app and their
	sub tasks information.

	Parameters:
		endpoint: Marathon HTTP endpoint, e.g. http://localhost:8080
*/
func FetchApps(endpoint string) ([]App, error) {
	taskContents, err := fetchTasks(endpoint)
	if err != nil { return nil, err }

	appConfiguration, err := fetchAppConfiguration(endpoint)
	if err != nil { return nil, err }

	apps := parseApps(taskContents, appConfiguration)
	return apps, nil
}
