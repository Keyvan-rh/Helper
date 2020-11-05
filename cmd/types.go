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

