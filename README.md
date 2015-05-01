# Bamboo  [![Build Status](https://travis-ci.org/QubitProducts/bamboo.svg?branch=master)](https://travis-ci.org/QubitProducts/bamboo)

![bamboo-logo](https://cloud.githubusercontent.com/assets/37033/4110258/a8cc58bc-31ef-11e4-87c9-dd20bd2468c2.png)

Bamboo is a web daemon that automatically configures
HAProxy for web services deployed on [Apache Mesos](http://mesos.apache.org) and [Marathon](https://mesosphere.github.io/marathon/).

It features:

* User interface for configuring HAProxy ACL rules for each Marathon application
* Rest API for configuring proxy ACL rules
* Auto configure HAProxy configuration file based your template; you can provision your own template in production to enable SSL and HAProxy stats interface, or configuring different load balance strategy
* Optionally handles health check endpoint if Marathon app is configured with [Healthchecks](https://mesosphere.github.io/marathon/docs/health-checks.html)
* Daemon itself is stateless; enables horizontal replication and scalability
* Developed in Golang, deployment on HAProxy instance has no additional dependency
* Optionally integrates with StatsD to monitor configuration reload event


### Compatibility

v0.1.1 supports Marathon 0.6 and Mesos 0.19.x

v0.2.2 supports both DNS and non-DNS proxy ACL rules

v0.2.8 supports both HTTP & TCP via custom Marathon enviroment variables (read below for details)

v0.2.9 supports Marathon 0.7.* (with [http_callback
enabled](https://mesosphere.github.io/marathon/docs/rest-api.html#event-subscriptions)) and Mesos 0.21.x

v0.2.11 improves API, deprecate previous API endpoint


### Releases and changelog

Since Marathon API and behaviour may change over time, especially in this early days. You should expect we aim to catch up those changes, improve design and adding new features. We aim to maintain backwards compatibility when possible. Releases and changelog are maintained in the [releases page](https://github.com/QubitProducts/bamboo/releases). Please read them when upgrading.

## Deployment Guide

You can deploy Bamboo with HAProxy on each Mesos slave. Each web service being allocated on Mesos Slave can discover services via localhost or domain you assigned by ACL rules. Alternatively, you can deploy Bamboo and HAProxy on separate instances, which means you need to loadbalance HAProxy cluster.

![bamboo-setup-guide](https://cloud.githubusercontent.com/assets/37033/4110199/a6226b8e-31ee-11e4-9734-68e0da00767c.png)

## User Interface

UI is useful to manage and visualize current state of proxy rules. Of course, you can configure HAProxy template to load balance Bamboo.

![user-interface-list](https://cloud.githubusercontent.com/assets/37033/4320901/527988dc-3f3b-11e4-8672-666605eb1ddf.png)

![user interface](https://cloud.githubusercontent.com/assets/37033/4320873/2ac65c48-3f3b-11e4-969e-52381dd33aae.png)

## StatsD Monitoring

![bamboo-graphite](https://cloud.githubusercontent.com/assets/37033/4117219/cef5cea2-328e-11e4-8346-ecc4e4e6046b.png)

## Configuration and Template

Bamboo binary accepts `-config` option to specify application configuration JSON file location. Type `-help` to get current available options.

Example configuration and HAProxy template can be found under [config/production.example.json](config/production.example.json) and  [config/haproxy_template.cfg](config/haproxy_template.cfg)
This section tries to explain usage in code comment style:

```JavaScript

{
  // Marathon instance configuration
  "Marathon": {
    // Marathon service HTTP endpoint
    "Endpoint": "http://localhost:8080"
  },

  "Bamboo": {

    // Bamboo's HTTP address can be accessed by Marathon
    // This is used for Marathon HTTP callback; must be reachable by Marathon
    "Host": "http://localhost:8000",

    // Proxy setting information is stored in Zookeeper
    // Bamboo will create this path if it does not already exist
    "Zookeeper": {
      // Use the same ZK setting if you run on the same ZK cluster
      "Host": "zk01.example.com:2812,zk02.example.com:2812",
      "Path": "/marathon-haproxy/state",
      "ReportingDelay": 5
    }
  }
  
  
  // Make sure using absolute path on production
  "HAProxy": {
    "TemplatePath": "/var/bamboo/haproxy_template.cfg",
    "OutputPath": "/etc/haproxy/haproxy.cfg",
    "ReloadCommand": "haproxy -f /etc/haproxy/haproxy.cfg -p /var/run/haproxy.pid -D -sf $(cat /var/run/haproxy.pid)"
  },
  
  // Enable or disable StatsD event tracking
  "StatsD": {
    "Enabled": false,
    // StatsD or Graphite server host
    "Host": "localhost:8125",
    // StatsD namespace prefix
    // If you have multiple Bamboo instances, you might want to label each node
    // by bamboo-server.production.n1.
    "Prefix": "bamboo-server.production."
  }
}
```

### Customize HAProxy Template with Marathon App Environment Variables

Marathon app env variables are available to be called in the template.
The default template shipped with Bamboo is aware of `BAMBOO_TCP_PORT`. When this variable is specified in Marathon app creation, the application will be configured with TCP mode. For example:

```JavaScript
{
  "id": "FileServer",
  "cmd": "python -m SimpleHTTPServer $PORT0",
  "cpus": 0.1,
  "mem": 90,
  "ports": [0],
  "instances": 2,
  "env": {
    "BAMBOO_TCP_PORT": "1080",
    "MY_CUSTOM_ENV": "hello"
  }
}
```

In this example, both `BAMBOO_TCP_PORT` and `MY_CUSTOM_ENV` can be accessed in HAProxy template. This enables flexible template customization depending on your preferences.

### Environment Variables

Configuration in the `production.json` file can be overridden with environment variables below. This is generally useful when you are building a Docker image for Bamboo and HAProxy. If they are not specified then the values from the configuration file will be used.

Environment Variable | Corresponds To
---------------------|---------------
`MARATHON_ENDPOINT` | Marathon.Endpoint
`BAMBOO_ENDPOINT` | Bamboo.Endpoint
`BAMBOO_ZK_HOST` | Bamboo.Zookeeper.Host
`BAMBOO_ZK_PATH` | Bamboo.Zookeeper.Path
`HAPROXY_TEMPLATE_PATH` | HAProxy.TemplatePath
`HAPROXY_OUTPUT_PATH` | HAProxy.OutputPath
`HAPROXY_RELOAD_CMD` | HAProxy.ReloadCommand
`BAMBOO_DOCKER_AUTO_HOST` | Sets `BAMBOO_ENDPOINT=$HOST` when Bamboo container starts. Can be any value.
`STATSD_ENABLED` | StatsD.Enabled
`STATSD_PREFIX` | StatsD.Prefix
`STATSD_HOST` | StatsD.Host


## REST APIs


#### GET /api/state

Shows the data structure used for rendering template

```bash
curl -i http://localhost:8000/api/state
```

#### GET /api/services

Shows all service configurations

```bash
curl -i http://localhost:8000/api/services
```

Example result:

```json
{
    "/authentication-service": {
        "Id": "/authentication-service",
        "Acl": "path_beg -i /authentication-service"
    },
    "/payment-service": {
        "Id": "/payment-service",
        "Acl": "path_beg -i /payment-service"
    }
}
```

#### POST /api/services

Creates a service configuration for a Marathon Application ID

```bash
curl -i -X POST -d '{"id":"/ExampleAppGroup/app1","acl":"hdr(host) -i app-1.example.com"}' http://localhost:8000/api/services
```

#### PUT /api/services/:id

Updates an existing service configuration for a Marathon application. `:id` is Marathon Application ID

```bash
curl -i -X PUT -d '{"id":"/ExampleAppGroup/app1", "acl":"path_beg -i /group/app-1"}' http://localhost:8000/api/services//ExampleAppGroup/app1
```


#### DELETE /api/services/:id

Deletes an existing service configuration. `:id` Marathon Application ID

```bash
curl -i -X DELETE http://localhost:8000/api/services//ExampleAppGroup/app1
```

#### GET /status

Bamboo webapp's healthcheck point

```
curl -i http://localhost:8000/status
```


## Deployment

We recommend installing binary with deb or rpm package. 
The repository includes examples of [a Jenkins build script](https://github.com/QubitProducts/bamboo/blob/master/builder/ci-jenkins.sh)
and [a deb packages build script](https://github.com/QubitProducts/bamboo/blob/master/builder/build.sh).
Read comments in the script to customize your build distribution workflow.

In short, [install fpm](https://github.com/jordansissel/fpm) and run the following command:

```
go build bamboo.go
./builder/build.sh
```

A deb package will be generated in `./builder` directory. You can copy to a server or publish to your own apt repository.

The example deb package deploys:

* Upstart job [`bamboo-server`](https://github.com/QubitProducts/bamboo/blob/master/builder/bamboo-server), e.g. upstart assumes `/var/bamboo/production.json` is configured correctly.
* Application directory is under `/opt/bamboo/`
* Configuration and logs is under `/var/bamboo/`
* Log file is rotated automatically

In case you're not using upstart, a template init.d service is provided in [`init.d-bamboo-server`](https://github.com/QubitProducts/bamboo/blob/master/builder/init.d-bamboo-server). Install it with
```
sudo cp builder/init.d-bamboo-server /etc/init.d/bamboo-server
sudo chown root:root /etc/init.d/bamboo-server
sudo chmod 755 /etc/init.d/bamboo-server
sudo update-rc.d "bamboo-server" defaults
```

You can then start the server with ```sudo service bamboo-server start```. Other commands: status, restart, stop

### As a Docker container

There is a `Dockerfile` that will allow Bamboo to be built and run from within a Docker container.

#### Building the image

The Docker image can be built and added to your local repository with the following command from within the project root directory:

````
docker build -t bamboo .
````

#### Running Bamboo as a Docker container

Once the image has been built, running as a container is straightforward - you do however still need to provide the configuration to the image as environment variables. Docker allows two options for this - using the `-e` option  or by putting them in a file and using the `--env-file` option. For this example we will use the former and we will map through ports 8000 and 80 to the docker host (obviously the hosts configured here will need to be reachable from this container):

````
docker run -t -i --rm -p 8000:8000 -p 80:80 \
    -e MARATHON_ENDPOINT=http://marathon:8080 \
    -e BAMBOO_ENDPOINT=http://bamboo:8000 \
    -e BAMBOO_ZK_HOST=zk:2181 \
    -e BAMBOO_ZK_PATH=/bamboo \
    -e BIND=":8000"
    -e CONFIG_PATH="config/production.example.json"
    -e BAMBOO_DOCKER_AUTO_HOST=true
    bamboo
````

Bamboo is started by supervisord in this Docker image. The [default Supervisord configuration](https://github.com/QubitProducts/bamboo/blob/master/builder/supervisord.conf) redirects stderr/stdout logs to the terminal. If you wish to turn the debug information off in production, you can use an [alternative configuration](https://github.com/QubitProducts/bamboo/blob/master/builder/supervisord.conf.prod).

## Development and Contribution

We use [godep](https://github.com/tools/godep) managing Go package dependencies; Goconvey for unit testing; CommonJS and SASS for frontend development and build distribution.
 
* Golang 1.3
* Node.js 0.10.x+

Golang:

```bash
# Pakcage manager
go get github.com/tools/godep
# Testing Toolkit
go get -t github.com/smartystreets/goconvey

cd $GOPATH/src/github.com/QubitProducts/bamboo

# Build your binary
go build

# Run test
goconvey
```

Node.js UI dependencies:

```bash
# Global 
npm install -g grunt-cli napa browserify node-static foreman karma-cli
# Local
npm install && napa

# Start a foreman configured with Procfile for building SASS and JavaScript 
nf start
```



## License

Bamboo is released under Apache License 2.0
