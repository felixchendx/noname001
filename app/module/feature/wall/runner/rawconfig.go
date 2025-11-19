package runner

type RawRunnerConfig struct {
	CfgDef      string      `yaml:"cfgdef"`
	Execute     bool        `yaml:"execute"`

	SegmentWall SegmentWall `yaml:"wall"`
}

type SegmentWall struct {
	Mode  string `yaml:"mode"`

	Walls []Wall `yaml:"walls"`
}
type Wall struct {
	Code           string `yaml:"code"`
	Name           string `yaml:"name"`
	State          string `yaml:"state"`
	Note           string `yaml:"note"`

	WallLayoutCode string `yaml:"wall_layout_code"`

	Items          []WallItem `yaml:"items"`
}
type WallItem struct {
	Index        int    `yaml:"index"`

	SourceNodeID string `yaml:"source_node"`
	StreamCode   string `yaml:"stream_code"`
}
