#!/bin/bash
PRG="goframehttp"

PRG_HOME="./"
LOG="/data/logs/$PRG.log"
CONFIG="config.toml"
BOOTPRG="$PRG_HOME$PRG $PRG_HOME$CONFIG"
CMD="eval ps -C $PRG --no-header |awk '{print \$1}'"
RETVAL=0


start()
{
PID=$($CMD)
if [[ -n $PID ]] ; then
    echo "$PRG is running,PID is $PID"
    return $RETVAL
elif [ ! -f $PRG_HOME$PRG ] || [ ! -x $PRG_HOME$PRG ] || [ ! -f $PRG_HOME$CONFIG ];then    #检查是否有执行权限
    echo "Start $PRG Failed! Cannot find $PRG_HOME$PRG or does not have execute permission!"
    RETVAL=1
    return $RETVAL
else
    echo $"Starting $PRG..."
    $BOOTPRG >>$LOG 2>&1 &
    RETVAL=$?
    sleep 1
    PID=$($CMD)
    if [[ -z $PID ]] ; then
        echo "$PRG Start Failed! "
        return $RETVAL
    fi
    echo "Start $PRG Success. PID is $PID"
fi
}

stop()
{
PID=$($CMD)
if [[ -z $PID ]] ; then
    echo "$PRG is not running."
    return $RETVAL
else
    echo "Stopping $PRG..."
    kill $PID
    RETVAL=$?
    sleep 1
    PID=$($CMD)
    if [[ ! -z $PID ]] ; then
        echo "Stop $PRG Failed! $RETVAL"
        return $RETVAL
    fi
    echo "$PRG Stopped."
    return $RETVAL
fi
}


restart()
{
stop
start
return $?
}

status()
{
PID=$($CMD)
if [[ -z $PID ]] ; then
    echo "$PRG is not running."
    return $RETVAL
else
    echo "$PRG is running,PID is $PID"
    return $RETVAL
fi
}


case $1 in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
esac
exit
