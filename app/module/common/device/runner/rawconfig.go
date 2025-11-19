package runner

type RawRunnerConfig struct {
	CfgDef  string `yaml:"cfgdef"`
	Execute bool   `yaml:"execute"`

	SegmentDevice SegmentDevice `yaml:"device"`
}

type SegmentDevice struct {
	Mode string `yaml:"mode"`

	Devices []Device `yaml:"devices"`
}
type Device struct {
	Code     string `yaml:"code"`
	Name     string `yaml:"name"`
	State    string `yaml:"state"`
	Note     string `yaml:"note"`
	Protocol string `yaml:"protocol"`
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Brand    string `yaml:"brand"`

	FallbackRTSPPort string `yaml:"fallback_rtsp_port"`
}
