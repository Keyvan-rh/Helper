#!/bin/bash
#
## Make sure the HELPERPOD_CONFIG_YAML env var has size
[[ ${#HELPERPOD_CONFIG_YAML} -eq 0 ]] && echo "FATAL: HELPERPOD_CONFIG_YAML env var not set!!!" && exit 254

#
## For now, let's test just writing the file out
echo ${HELPERPOD_CONFIG_YAML} | base64 -d > /usr/local/src/helperpod.yaml

#
## Convert the YAML to JSON because it's easier to work with
python3 -c 'import sys, yaml, json; json.dump(yaml.safe_load(sys.stdin), sys.stdout, indent=4)' < /usr/local/src/helperpod.yaml > /usr/local/src/helperpod.json

#
## TBD:
### * Create a Template haproxy.cfg and add it in the Containerfile
### * Using the YAML/JSON provided, create an haproxy.cfg file "on the fly"

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
	rm -f ${haproxyPidFile}
	/usr/sbin/haproxy -Ws -f ${haproxyConfig} -p ${haproxyPidFile}
fi
##
##
