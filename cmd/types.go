package cmd

type HelpMe struct {
	Runtime  Runtime    `yaml:"runtime"`
}
type Runtime struct {
	Services []Service `yaml:"services"`
}
type Service struct {
	Service string `yaml:"service"`
	Run bool `yaml:"run"`
}
const QUAY string = "quay.io/helpernode"
const DEFAULTTAG string = "latest"

//TODO this needs to be removed and start/stop updated once we are reading from viper config file
var images =  map[string]string{
	"dns": "quay.io/helpernode/dns",
	"dhcp": "quay.io/helpernode/dhcp",
	"http": "quay.io/helpernode/http",
	"loadbalancer": "quay.io/helpernode/loadbalancer",
	"pxe": "quay.io/helpernode/pxe",
}

