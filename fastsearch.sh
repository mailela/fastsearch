#!/bin/bash

#每分钟检测fastsearch运行
#*/1 * * * * /data/fastsearch/fastsearch.sh > /dev/null 2>&1

#每3点 重启fastsearch
#0 3 * * * /etc/init.d/fastsearch.d restart

count=`ps -fe |grep "fastsearch"|grep "config.yaml" -c`

echo $count
if [ $count -lt 1 ]; then
	echo "restart"
	echo $(date +%Y-%m-%d_%H:%M:%S) >/data/fastsearch/restart.log 
	/etc/init.d/fastsearch.d restart
else
	echo "is running"
fi