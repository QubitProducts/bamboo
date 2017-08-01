package marathon

import (
	"fmt"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/QubitProducts/bamboo/configuration"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetMesosDnsId_Simple(t *testing.T) {
	Convey("#getMesosDnsId", t, func() {
		Convey("should return simple appname", func() {
			So(getMesosDnsId("appname"), ShouldEqual, "appname")
		})

		Convey("should return simple appname if slash prefixed", func() {
			So(getMesosDnsId("/appname"), ShouldEqual, "appname")
		})

		Convey("should return groups reverse-added to appname", func() {
			So(getMesosDnsId("/group/appname"), ShouldEqual, "appname-group")
		})

		Convey("should return groups reverse-added to appname but no blanks", func() {
			So(getMesosDnsId("//group/again//appname/"), ShouldEqual, "appname-again-group")
		})
	})
}

func TestParseHealthCheckPathTCP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []marathonHealthCheck{
			marathonHealthCheck{"/", "TCP", 0},
			marathonHealthCheck{"/foobar", "TCP", 0},
			marathonHealthCheck{"", "TCP", 0},
		}
		Convey("should return no path if all checks are TCP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "")
		})
	})
}

func TestParseHealthCheckPathHTTP(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []marathonHealthCheck{
			marathonHealthCheck{"/first", "HTTP", 0},
			marathonHealthCheck{"/", "HTTP", 0},
			marathonHealthCheck{"", "HTTP", 0},
		}
		Convey("should return the first path if all checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/first")
		})
	})
}

func TestParseHealthCheckPathMixed(t *testing.T) {
	Convey("#parseHealthCheckPath", t, func() {
		checks := []marathonHealthCheck{
			marathonHealthCheck{"", "TCP", 0},
			marathonHealthCheck{"/path", "HTTP", 0},
			marathonHealthCheck{"/", "HTTP", 0},
		}
		Convey("should return the first path if some checks are HTTP", func() {
			So(parseHealthCheckPath(checks), ShouldEqual, "/path")
		})
	})
}

func TestParseJSONRequest(t *testing.T) {
	tests := []struct {
		user          string
		password      string
		wantBasicAuth bool
	}{
		{
			wantBasicAuth: false,
		},
		{
			user:          "user",
			wantBasicAuth: false,
		},
		{
			password:      "password",
			wantBasicAuth: false,
		},
		{
			user:          "user",
			password:      "password",
			wantBasicAuth: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("user='%s' password='%s'", test.user, test.password), func(t *testing.T) {
			t.Parallel()
			conf := configuration.Configuration{
				Marathon: configuration.Marathon{
					User:     test.user,
					Password: test.password,
				},
			}

			var req *http.Request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				req = r
				fmt.Fprint(w, "{}")
			}))
			defer ts.Close()

			var res interface{}
			err := parseJSON(ts.URL, &conf, &res)
			if err != nil {
				t.Fatalf("parseJSON returned error: %s", err)
			}

			if req.Method != http.MethodGet {
				t.Errorf("got method '%s', want '%s'", req.Method, http.MethodGet)
			}

			for _, hdrKey := range []string{"Accept", "Content-Type"} {
				hdrValue := req.Header.Get(hdrKey)
				switch {
				case hdrValue == "":
					t.Errorf("%s header missing", hdrKey)
				case hdrValue != "application/json":
					t.Errorf("got %s header value '%s', want 'application/json'", hdrKey, hdrValue)
				}
			}

			authHdrValue := req.Header.Get("Authorization")
			if test.wantBasicAuth != (authHdrValue != "") {
				t.Errorf("got Authorization header value '%s', wanted header: %t", authHdrValue, test.wantBasicAuth)
			}
		})
	}
}

func TestParseJSONHandling(t *testing.T) {
	tests := []struct {
		desc          string
		handler       http.Handler
		shouldSucceed bool
	}{
		{
			desc:          "request failed",
			shouldSucceed: false,
		},
		{
			desc: "invalid JSON",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "{")
			}),
			shouldSucceed: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			var endpoint string
			if test.handler != nil {
				ts := httptest.NewServer(test.handler)
				defer ts.Close()
				endpoint = ts.URL
			}

			conf := configuration.Configuration{}
			var res interface{}
			err := parseJSON(endpoint, &conf, res)

			if test.shouldSucceed != (err == nil) {
				t.Errorf("got error '%s', wanted error: %t", err, !test.shouldSucceed)
			}
		})
	}
}

func TestFetchApps(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, `{
	"apps": [
		{
			"id": "/app2WithSlash",
			"tasks": [
				{
					"host":"localhost",
					"id": "task2",
					"ports": [8002]
				},
				{
					"host":"localhost",
					"id": "task1",
					"ports": [8001]
				}
			]
		},
		{
			"id": "app1WithoutSlash",
			"tasks": [
				{
					"host":"localhost",
					"id": "task1",
					"ports": [8001]
				},
				{
					"host":"localhost",
					"id": "task2",
					"ports": [8002]
				}
			]
		},
		{
			"id": "app3WithIpAddress",
			"ipAddress": {
				"discovery":{
					"ports": [
						{
							"number": 80	
						}
					]
				}	
			},
			"tasks": [
				{
					"host":"localhost",
					"ipAddresses":[
						{
							"ipAddress": "127.0.0.1"
						}
					],
					"id": "task1"
				},
				{
					"host":"localhost",
					"ipAddresses":[
						{
							"ipAddress": "127.0.0.1"
						}
					],
					"id": "task2"
				}
			]
		}
	]
}`)
	}))
	defer ts.Close()

	// First Marathon URL is invalid to verify failover behavior.
	maraConf := configuration.Marathon{
		Endpoint: fmt.Sprintf("http://127.0.0.1:4242,%s", ts.URL),
	}

	apps, err := FetchApps(maraConf, &configuration.Configuration{})

	if err != nil {
		t.Fatalf("FetchApps returned error: %s", err)
	}

	if len(apps) < 1 {
		t.Fatal("no apps fetched")
	}
	assertFetchedApp(t, 1, "/app1WithoutSlash", apps[0])

	if len(apps) < 2 {
		t.Fatal("missing second app")
	}
	assertFetchedApp(t, 2, "/app2WithSlash", apps[1])

	if len(apps) < 3 {
		t.Fatal("missing third app")
	}
	assertFetchedApp(t, 3, "/app3WithIpAddress", apps[2])

	if len(apps) > 3 {
		t.Fatalf("got %d apps, want 3", len(apps))
	}
}

func TestCalculateReadiness(t *testing.T) {
	tests := []struct {
		desc      string
		task      marathonTask
		app       marathonApp
		rc        readinessCalculator
		wantReady bool
	}{
		{
			desc: "non-running task",
			task: marathonTask{
				State: "TASK_STAGED",
			},
			wantReady: false,
		},
		{
			desc: "no deployment running for app",
			task: marathonTask{
				State: taskStateRunning,
			},
			app: marathonApp{
				Deployments: []deployment{},
			},
			wantReady: true,
		},
		{
			desc: "no readiness checks defined for app",
			task: marathonTask{
				State: taskStateRunning,
			},
			app: marathonApp{
				Deployments: []deployment{
					deployment{ID: "deploymentId"},
				},
				ReadinessChecks: []marathonReadinessCheck{},
			},
			wantReady: true,
		},
		{
			desc: "readiness check result negative",
			task: marathonTask{
				Id:    "taskId",
				State: taskStateRunning,
			},
			app: marathonApp{
				Deployments: []deployment{
					deployment{ID: "deploymentId"},
				},
				ReadinessChecks: []marathonReadinessCheck{
					marathonReadinessCheck{
						Path: "/ready",
					},
				},
				ReadinessCheckResults: []readinessCheckResult{
					readinessCheckResult{
						Ready:  false,
						TaskID: "taskId",
					},
				},
			},
			wantReady: false,
		},
		{
			desc: "readiness check result positive",
			task: marathonTask{
				Id:    "taskId",
				State: taskStateRunning,
			},
			app: marathonApp{
				Deployments: []deployment{
					deployment{ID: "deploymentId"},
				},
				ReadinessChecks: []marathonReadinessCheck{
					marathonReadinessCheck{
						Path: "/ready",
					},
				},
				ReadinessCheckResults: []readinessCheckResult{
					readinessCheckResult{
						Ready:  false,
						TaskID: "otherTaskId",
					},
					readinessCheckResult{
						Ready:  true,
						TaskID: "taskId",
					},
				},
			},
			wantReady: true,
		},
		{
			desc: "no readiness check result with default timeout",
			task: marathonTask{
				Id:        "newTaskId",
				State:     taskStateRunning,
				StartedAt: timestampOffsetFromNow(3 * time.Minute),
			},
			app: marathonApp{
				Deployments: []deployment{
					deployment{ID: "deploymentId"},
				},
				ReadinessChecks: []marathonReadinessCheck{
					marathonReadinessCheck{
						Path: "/ready",
					},
				},
				ReadinessCheckResults: []readinessCheckResult{},
			},
			rc: readinessCalculator{
				checkDefaultTimeout: 5 * time.Minute,
			},
			wantReady: false,
		},
		{
			desc: "no readiness check result with readiness check timeout",
			task: marathonTask{
				Id:        "newTaskId",
				State:     taskStateRunning,
				StartedAt: timestampOffsetFromNow(4 * time.Minute),
			},
			app: marathonApp{
				Deployments: []deployment{
					deployment{ID: "deploymentId"},
				},
				ReadinessChecks: []marathonReadinessCheck{
					marathonReadinessCheck{
						Path:           "/ready",
						TimeoutSeconds: int((3 * time.Minute).Seconds()),
					},
				},
				ReadinessCheckResults: []readinessCheckResult{},
			},
			rc: readinessCalculator{
				checkSafetyMargin: 3 * time.Minute,
			},
			wantReady: false,
		},
		{
			desc: "invalid task start time",
			task: marathonTask{
				Id:        "newTaskId",
				State:     taskStateRunning,
				StartedAt: "invalid",
			},
			app: marathonApp{
				Deployments: []deployment{
					deployment{ID: "deploymentId"},
				},
				ReadinessChecks: []marathonReadinessCheck{
					marathonReadinessCheck{
						Path: "/ready",
					},
				},
				ReadinessCheckResults: []readinessCheckResult{},
			},
			wantReady: false,
		},
		{
			desc: "task not involved in deployment",
			task: marathonTask{
				Id:        "oldTaskId",
				State:     taskStateRunning,
				StartedAt: timestampOffsetFromNow(1 * time.Hour),
			},
			app: marathonApp{
				Deployments: []deployment{
					deployment{ID: "deploymentId"},
				},
				ReadinessChecks: []marathonReadinessCheck{
					marathonReadinessCheck{
						Path: "/ready",
					},
				},
				ReadinessCheckResults: []readinessCheckResult{},
			},
			rc: readinessCalculator{
				checkDefaultTimeout: 10 * time.Second,
			},
			wantReady: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			rc := readyCalculator
			if test.rc.checkDefaultTimeout > 0 {
				rc.checkDefaultTimeout = test.rc.checkDefaultTimeout
			}
			if test.rc.checkSafetyMargin > 0 {
				rc.checkSafetyMargin = test.rc.checkSafetyMargin
			}
			gotReady := test.rc.calculate(test.task, test.app)
			if gotReady != test.wantReady {
				t.Errorf("got ready = %t, want ready = %t", gotReady, test.wantReady)
			}
		})
	}
}

func assertFetchedApp(t *testing.T, index int, id string, app App) {
	if app.Id != id {
		t.Errorf("app #%d: got app ID '%s', want '%s'", index, app.Id, id)
	}
	switch {
	case len(app.Tasks) != 2:
		t.Errorf("app #%d: got %d tasks, want 2", index, len(app.Tasks))
	case app.Tasks[0].Id != "task1":
		t.Errorf("app #%d: got ID '%s' for task #1, want 'task1", index, app.Tasks[0].Id)
	case app.Tasks[1].Id != "task2":
		t.Errorf("app #%d: got ID '%s' for task #2, want 'task2", index, app.Tasks[1].Id)
	}
}

func timestampOffsetFromNow(offset time.Duration) string {
	return time.Now().Add(-offset).Format(time.RFC3339)
}
