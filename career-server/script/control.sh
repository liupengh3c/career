#!/bin/sh
APP_NAME="career-server"
WORKROOT=$(cd $(dirname "$0") && pwd)
SUPERVISE=${WORKROOT}/supervise/supervise64
TMP=`cd $(dirname "$0") && pwd | awk -F "script" '{print $1}'`
SUPERVISE_CONF=${TMP}conf/supervise.conf
STATUS_DIR=${WORKROOT}/supervise/status/${APP_NAME}
CMD=${TMP}$APP_NAME
Start(){
    echo "begin to start service"
    echo "SUPERVISE=$SUPERVISE"
    echo "SUPERVISE_CONF=$SUPERVISE_CONF"
    echo "STATUS_DIR=$STATUS_DIR"
    echo "CMD=$CMD"
    cd ${TMP}
    echo $SUPERVISE -f "$CMD" -p "$STATUS_DIR" -F "$SUPERVISE_CONF"
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:${TMP}lib/
    if [ ! -d "$STATUS_DIR" ];then
    mkdir -p $STATUS_DIR
    fi
    $SUPERVISE -f "$CMD" -p "$STATUS_DIR" -F "$SUPERVISE_CONF" &>/dev/null </dev/null &
    sleep 2
    monitor_pid=$(ps axo pid,cmd | grep "${CMD}" | grep -v "supervise64"|grep -v "grep" | head -n 2 | tail -n 1 | awk '{print $1}')
    if [ -d /proc/${monitor_pid} ]
    then
        echo "${APP_NAME} start success: pid=${monitor_pid}"
    else
        echo "${APP_NAME} start failed"
    fi
    return
}

Stop(){
    # echo "begin to stop service"
    # echo `ps -ef | grep "${APP_NAME}" | grep -v "grep" | grep -v "control.sh" | awk '{print $0}'`
    # ID=`ps -ef | grep "${APP_NAME}" | grep -v "grep" | grep -v "control.sh" | awk '{print $2}'`
    # for id in $ID
    # do
    # kill -9 $id
    # echo "kill -9 $id"
    # done
    # sleep 3
    killall -9 supervise64 ${APP_NAME} || true
}

Restart(){
    echo "begin to restart service"
    Stop
    Start
}

cmd=$1
if [ $cmd == "start" ]
then
    Start
elif [ $cmd == "stop" ]
then
    Stop
else
    Restart
fi
