package marathon

import (
	"net/http"
	"io/ioutil"
	"strings"
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
	Tasks []Task
}

func fetchTasks(endpoint string) (string, error) {
	
	response, err := http.Get(endpoint + "/v2/tasks")
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

func parseApps(contents string) []App {
	lines := strings.Split(contents, "\n")
	apps := []App{}
	for _, line := range lines {
		if len(line) > 0 {
			appId, appPort, tasks := parseTasks(line)
			app := App{ Id: appId, Port: appPort, Tasks: tasks }
			apps = append(apps, app)
		}
	}
	return apps
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
func Apps(endpoint string) ([]App, error) {
	contents, err := fetchTasks(endpoint)

	if err != nil {
		return nil, err
	}

	apps := parseApps(contents)
	return apps, nil
}
