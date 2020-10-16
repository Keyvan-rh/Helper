#!/bin/bash
#
## Make sure the HELPERPOD_CONFIG_YAML env var has size
[[ ${#HELPERPOD_CONFIG_YAML} -eq 0 ]] && echo "FATAL: HELPERPOD_CONFIG_YAML env var not set!!!" && exit 254

#
## For now, let's test just writing the file out
echo ${HELPERPOD_CONFIG_YAML} | base64 -d > /usr/local/src/helperpod.yaml

#
## Create haproxy.cfg based on the template and yaml passed in.
ansible  localhost --connection=local -e @/usr/local/src/helperpod.yaml -m template -a "src=/usr/local/src/haproxy.cfg.j2 dest=/etc/haproxy/haproxy.cfg" > /var/log/helperpod_ansible_run.log 2>&1

#
## Set HAProxy variables
haproxyConfig=/etc/haproxy/haproxy.cfg 
haproxyPidFile=/run/haproxy.pid

#
## Test for the validity of the config file. Run the HAProxy process if it passes
if ! /usr/sbin/haproxy -f ${haproxyConfig} -c -q ; then
	echo "============================="
	echo "FATAL: Invalid HAProxy config"
	echo "============================="
	exit 254
else
	echo "==========================="
	echo "Starting HAproxy service..."
	echo "==========================="
	rm -f ${haproxyPidFile}
	/usr/sbin/haproxy -Ws -f ${haproxyConfig} -p ${haproxyPidFile}
fi
##
##
