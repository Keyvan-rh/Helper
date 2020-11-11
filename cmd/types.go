package cmd

const VERSION string = "latest"

//TODO probably can delete this and use VERSION
const DEFAULTTAG string = "latest"

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

