# Bamboo  [![Build Status](https://travis-ci.org/QubitProducts/bamboo.svg?branch=master)](https://travis-ci.org/QubitProducts/bamboo)

![bamboo-logo](https://cloud.githubusercontent.com/assets/37033/4110258/a8cc58bc-31ef-11e4-87c9-dd20bd2468c2.png)

Bamboo is a web daemon that automatically configures
HAProxy for web services deployed on [Apache Mesos](http://mesos.apache.org) and [Marathon](https://mesosphere.github.io/marathon/).

It features:

* User interface for configuring DNS mapping to Marathon ID
* Rest API for configuring DNS mapping; when application deployment in Marathon is automated and work with other DNS managing services, API is handy
* Auto configure HAProxy configuration file based your template; you can provision your own template in production to enable ssl and HAProxy stats interface, or configuring different load balance strategy
* Optionally hanldes health check endpoint if Marathon app is configured with [Healthchecks](https://mesosphere.github.io/marathon/docs/health-checks.html)
* Daemon itself is stateless; enables horizontal replication and scalability
* Developed in Golang, deployment on HAProxy instance has no additional dependency
* Optionally integrates with StatsD to monitor configuration reload event  

![user-interface](https://cloud.githubusercontent.com/assets/37033/4110199/a6226b8e-31ee-11e4-9734-68e0da00767c.png)

## User Interface

If you have very small scale web services to manage, UI is useful to manage and visualize current state of DNS mapping.
You can find out if DNS is assigned or missing from the interface. Of course, you can configure HAProxy template to load balance Bamboo. 

![user-interface](https://cloud.githubusercontent.com/assets/37033/4109769/318f2ad2-31e9-11e4-8f5f-b6a3368412b9.png)

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
    "Endpoint": "http://localhost:8080",
    // Same configuration as Marathon Zookeeper
    "Zookeeper": {
      "Host": "zk01.example.com:2812,zk02.example.com:2812",
      // Marathon Zookeeper state  
      // Marathon default set to /marathon/state
      "Path": "/marathon/state",
      // Number of seconds to delay the reload event
      "ReportingDelay": 5
    }
  },
   
  // DNS mapping information is stored in Zookeeper
  // Make sure the Path is pre-created; Bamboo does not create the missing Path.
  "DomainMapping": {
    "Zookeeper": {
      // Use the same ZK setting if you run on the same ZK cluster
      "Host": "zk01.example.com:2812,zk02.example.com:2812",
      "Path": "/marathon-haproxy/state",
      "ReportingDelay": 5
    }
  },
  
  // Make sure using absolute path on production
  "HAProxy": {
    "TemplatePath": "/var/bamboo/haproxy_template.cfg",
    "OutputPath": "/etc/haproxy/haproxy.cfg",
    "ReloadCommand": "haproxy -f /etc/haproxy/haproxy.cfg -p /var/run/haproxy.pid -sf $(cat /var/run/haproxy.pid)"
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

### Environment Variables

Configuration in the `production.json` file can be overridden with environment variables below. This is generally useful when you are building a Docker image for Bamboo and HAProxy. If they are not specified then the values from the configuration file will be used.

Environment Variable | Corresponds To
---------------------|---------------
`MARATHON_ZK_HOST` | Marathon.Zookeeper.Host
`MARATHON_ZK_PATH` | Marathon.Zookeeper.Path
`MARATHON_ENDPOINT` | Marathon.Endpoint
`DOMAIN_ZK_HOST` | DomainMapping.Zookeeper.Host
`DOMAIN_ZK_PATH` | DomainMapping.Zookeeper.Path
`HAPROXY_TEMPLATE_PATH` | HAProxy.TemplatePath
`HAPROXY_OUTPUT_PATH` | HAProxy.OutputPath
`HAPROXY_RELOAD_CMD` | HAProxy.ReloadCommand


## REST APIs

POST /api/state/domains

Creates mapping from Marathon application ID to a DNS

```bash
curl -i -X POST -d '{"id":"app-1","value":"app1.example.com"}' http://localhost:8000/api/state/domains
```


PUT /api/state/domains/:id

Updates mapping of an existing Marathon application ID to a new DNS

```bash
curl -i -X PUT -d '{"id":"app-1","value":"app1-beta.example.com"}' http://localhost:8000/api/state/domains/app-1
```


DELETE /api/state/domains/:id

Deletes mapping of an existing Marathon ID DNS mapping

```bash
curl -i -X DELETE http://localhost:8000/api/state/domains/app-1
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

## Development and Contribution

We use [godep](https://github.com/tools/godep) managing Go package dependencies; Goconvey for unit testing; CommonJS and SASS for frontend development and build distribution.
 
* Golang 1.3
* Node.js 0.10.x+

Golang:

```bash
# Pakcage manager
go get github.com/tools/godep
go get -t github.com/smartystreets/goconvey

cd $GOPATH/github.com/QubitProducts/bamboo
godep restore

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
