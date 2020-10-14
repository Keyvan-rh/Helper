#!/bin/bash
#
## Make sure the HELPERPOD_CONFIG_YAML env var has size
[[ ${#HELPERPOD_CONFIG_YAML} -eq 0 ]] && echo "FATAL: HELPERPOD_CONFIG_YAML env var not set!!!" && exit 254

#
## For now, let's test just writing the file out
echo ${HELPERPOD_CONFIG_YAML} | base64 -d > /usr/local/src/helperpod.yaml

#
## Start HAProxy
haproxyConfig=/etc/haproxy/haproxy.cfg 
haproxyPidFile=/run/haproxy.pid
## This is how you test for a valid file...cloud be useful later
### /usr/sbin/haproxy -f ${haproxyConfig} -c -q
rm -f ${haproxyPidFile}
/usr/sbin/haproxy -Ws -f ${haproxyConfig} -p ${haproxyPidFile}
##
##
