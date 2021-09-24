#! /bin/bash

### BEGIN INIT INFO
# Provides:          aes
# Required-Start:    $local_fs $network
# Required-Stop:     $local_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: aes67 service
# Description:       Run AES67 service
### END INIT INFO

pname="aes"
exe="/opt/aes/aes67-daemon"
args="-c /opt/aes/daemon.conf"
pidfile="/var/run/${pname}.pid"
lockfile="/var/lock/subsys/${pname}"
log="/var/log/$pname"

[ -x $exe ] || exit 0

common_opts="--quiet --pidfile $pidfile"

# Carry out specific functions when asked to by the system
case "$1" in
  start)
    echo "Starting $pname ..."
    insmod /opt/aes/MergingRavennaALSA.ko > /dev/null 2>&1 || true
    start-stop-daemon --start $common_opts --make-pidfile --background --startas /bin/bash -- -c "$exe $args > $log 2>&1"
    ;;
  stop)
    echo "Shutting down $pname ..."
    start-stop-daemon --stop $common_opts --signal INT --remove-pidfile
    ;;
  *)
    echo "Usage: /etc/init.d/$pname {start|stop}"
    exit 1
    ;;
esac

exit 0
