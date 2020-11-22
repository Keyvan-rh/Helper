package cmd

const VERSION string = "latest"

//TODO probably can delete this and use VERSION
const DEFAULTTAG string = "latest"
var logLevel string


type HelpMe struct {
	Runtime  Runtime    `yaml:"a_runtime"`

}
type Runtime struct {
	dns string `yaml:"dns"`
	dhcp string `yaml:"dhcp"`
	http string `yaml:"http"`
	loadbalancer string `yaml:"loadbalancer"`
	pxe string `yaml:"pxe"`
}
type Service struct {
	Service string `yaml:"service"`
	Run bool `yaml:"run"`
}
const QUAY string = "quay.io/helpernode"

// Define ports needed for preflight check
var ports = [10]string{"67", "546", "53", "80", "443", "69", "6443", "22623", "8080", "9000"}

//Default images
//TODO Add disconnected
var images =  map[string]string{
	"dns": "quay.io/helpernode/dns",
	"dhcp": "quay.io/helpernode/dhcp",
	"http": "quay.io/helpernode/http",
	"loadbalancer": "quay.io/helpernode/loadbalancer",
	"pxe": "quay.io/helpernode/pxe",
}

// Define systemd services we will check
var systemdsvc = map[string]string {
	"resolved": "systemd-resolved.service",
	"dnsmasq" : "dnsmasq.service",
}

var fwrule = [13]string {
	"6443/tcp",
	"22623/tcp",
	"8080/tcp",
	"9000/tcp",
	"9090/tcp",
	"67/udp",
	"546/udp",
	"53/tcp",
	"53/udp",
	"80/tcp",
	"443/tcp",
	"22/tcp",
	"69/udp",
}

var clients = map[string]string {
	"oc": "openshift-client-linux.tar.gz",
	"openshift-install": "openshift-install-linux.tar.gz",
	"helm": "helm.tar.gz",
}