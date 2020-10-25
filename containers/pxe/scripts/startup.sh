#!/bin/bash
#
## This is the startup script for the LoadBalancer container

#
## Variables for PXE
tftpBootDir=/var/lib/tftpboot
pxeConfig=${tftpBootDir}/pxelinux.cfg
rhcosDir=${tftpBootDir}/rhcos
bootstrapPxeTemplate=/usr/local/src/pxe-bootstrap.j2
masterPxeTemplate=/usr/local/src/pxe-master.j2
workerPxeTemplate=/usr/local/src/pxe-worker.j2
helperPodYaml=/usr/local/src/helperpod.yaml
ansibleLog=/var/log/helperpod_ansible_run.log

#
## Make sure the HELPERPOD_CONFIG_YAML env var has size
[[ ${#HELPERPOD_CONFIG_YAML} -eq 0 ]] && echo "FATAL: HELPERPOD_CONFIG_YAML env var not set!!!" && exit 254

#
## Take the HELPERPOD_CONFIG_YAML env variable and write out the YAML file.
echo ${HELPERPOD_CONFIG_YAML} | base64 -d > ${helperPodYaml}

#
## Create PXE dir in the tftp dir
mkdir -m 0755 -p ${pxeConfig}

#
## Create a directory for RHCOS artifacts
mkdir -m 0755 -p ${rhcosDir}

#
## Copy over files needed for TFTP
cp -a /usr/share/syslinux/* ${tftpBootDir}/

#
## Downloading OCP4 installer initramfs and kernel files. Setting the proper permissions to 0555
wget https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/${OCPVERSION%.*}/${OCPRHCOSVERSION}/rhcos-${OCPRHCOSVERSION}-x86_64-installer-initramfs.x86_64.img -O ${rhcosDir}/initramfs.img
wget https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/${OCPVERSION%.*}/${OCPRHCOSVERSION}/rhcos-${OCPRHCOSVERSION}-x86_64-installer-kernel-x86_64 -O ${rhcosDir}/kernel
chmod 0555 ${rhcosDir}/{initramfs.img,kernel}

#
## Create pxe/tftp files based on the template and yaml passed in.
rm -f  ${ansibleLog}
ansible localhost -c local -e @${helperPodYaml} -m template -a 'src=${bootstrapPxeTemplate} dest="${pxeConfig}/01-{{ bootstrap.macaddr | lower | regex_replace (':', '-') }}" mode=0555' >> ${ansibleLog} 2>&1
ansible localhost -c local -e @${helperPodYaml} -m template -a 'src=${masterPxeTemplate} dest="${pxeConfig}/01-{{ item.macaddr | regex_replace (':', '-')}}" mode=0555 with_items="{{ masters | lower }}"' >> ${ansibleLog} 2>&1
ansible localhost -c local -e @${helperPodYaml} -m template -a 'src=${workerPxeTemplate} dest="${pxeConfig}/01-{{ item.macaddr | regex_replace (':', '-')}}" mode=0555 with_items="{{ workers | lower }}" when:"(workers is defined) and (workers | length > 0)"' >> ${ansibleLog} 2>&1

#
## PXE is a "best effort" service that is kind of "old". So putting this here as a placeholder until someone has time to write a "checker"
if false ; then
	echo "=============================="
	echo "FATAL: Invalid PXE/TFTP config"
	echo "=============================="
	exit 254
else
	echo "============================"
	echo "Starting PXE/TFTP service..."
	echo "============================"
	/usr/sbin/in.tftpd -L --verbosity 4  --secure ${tftpBootDir}
fi
##
##
