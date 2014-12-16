#!/bin/sh
### BEGIN INIT INFO
# Provides:          bamboo-server
# Required-Start:    $network $zookeeper $marathon $mesos
# Required-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:      1 2 3 4 5
# Short-Description:  bamboo server
# Description:       bamboo server
### END INIT INFO

SERVICE_NAME=bamboo-service
DAEMON=/opt/bamboo/bamboo
DAEMON_OPTS="-config /var/bamboo/production.json -log /var/bamboo/log/bamboo-server.log"
PIDFILE=/var/run/bamboo.pid

if [ ! -x $DAEMON ]; then
  echo "ERROR: Can't execute $DAEMON."
  exit 1
fi

start_service() {
  echo -n " * Starting $SERVICE_NAME... "
  cd /opt/bamboo/
  sleep 10
  start-stop-daemon --quiet --background --start --pidfile "$PIDFILE" --make-pidfile --exec $DAEMON -- $DAEMON_OPTS
  e=$?
  if [ $e -eq 1 ]; then
    echo "already running"
    return
  fi

  if [ $e -eq 255 ]; then
    echo "couldn't start :("
    exit 1
  fi

  echo "done"
}

stop_service() {
  echo -n " * Stopping $SERVICE_NAME... "
  start-stop-daemon -Kq -R 10 -p $PIDFILE
  e=$?
  if [ $e -eq 1 ]; then
    echo "not running"
    return
  fi

  echo "done"
}

status_service() {
    if [ -f $PIDFILE ]; then
        PID=`cat $PIDFILE`
        if [ -z "`ps axf | grep ${PID} | grep -v grep`" ]; then
            printf "%s\n" "$SERVICE_NAME dead but pidfile exists"
            exit 1 
        else
            printf "%-50s" "$SERVICE_NAME is running"
        fi
    else
        printf "%s\n" "$SERVICE_NAME not running"
        exit 3 
    fi
}

case "$1" in
  status)
    status_service
    ;;
  start)
    start_service
    ;;
  stop)
    stop_service
    ;;
  restart)
    stop_service
    start_service
    ;;
  *)
    echo "Usage: service $SERVICE_NAME {start|stop|restart|status}" >&2
    exit 1   
    ;;
esac

exit 0
