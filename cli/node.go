package cli

import (
	"fmt"
	"os"
	"strings"
	
	ufc "github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"

	"noname001/config/rawconfig"
	"noname001/root"

	"noname001/app/base/sanitation"
)

// TODO: any panic will cause unclean resource handling
//       do kill prev processes before start

func nodeRootCommand() (*ufc.Command) {
	return &ufc.Command{
		Name: "node",
		Usage: "Start this Node",
		Flags: []ufc.Flag{
		},
		Subcommands: []*ufc.Command{
			_nodeSubCommandStart(),
		},
	}
}

func _nodeSubCommandStart() (*ufc.Command) {
	flagRunner := &ufc.StringFlag{Name: "runner", Hidden: true, Value: "", Required: false}
	flagDebug := &ufc.BoolFlag{Name: "debug", Value: false, Required: false}

	return &ufc.Command{
		Name: "start",
		Usage: "Run this Node in normal mode",
		Category: "run",
		Flags: []ufc.Flag{
			flagNodeConfig, flagDebug, flagRunner,
		},
		Action: func(ctx *ufc.Context) (error) {
			var err error
			// TODO: dep check here, so it need to be done only once
			// > apt-get update -y
			// > apt-get install pkg-config libzmq3-dev

			// TODO: verify
			// libzmq3-dev stuffs and sqlite stuffs is dev only requirement
			// final binary does not need any apt install
			// consider embed ffmpeg and ffprobe for posterity
			// at least until there's reachable server that provides depended binaries

			rawConfigRoot := _loadRawConfig(ctx.String("config"))
			rawRunnerConfigRoot := _loadRawRunner(ctx.String("runner"))

			err = root.Start(rawConfigRoot, rawRunnerConfigRoot)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func _loadRawConfig(configFilepath string) (*rawconfig.ConfigRoot) {
	// https://ini.unknwon.io/docs/

	iniConfig, err := ini.Load(configFilepath)
	if err != nil {
		fmt.Println("config:", err)
		os.Exit(1)
	}

	rawConfigRoot := &rawconfig.ConfigRoot{}
	err = iniConfig.StrictMapTo(rawConfigRoot)
	if err != nil {
		fmt.Println("config:", err)
		os.Exit(1)
	}

	_validateRequired(rawConfigRoot)

	return rawConfigRoot
}

func _loadRawRunner(runnerConfigFilepath string) (*rawconfig.RunnerConfigRoot) {
	if runnerConfigFilepath == "" { return nil }
	
	fileContent, err := os.ReadFile(runnerConfigFilepath)
	if err != nil {
		fmt.Println("runner:", err)
		os.Exit(1)
	}

	rawRunnerConfigRoot := &rawconfig.RunnerConfigRoot{}
	err2 := yaml.Unmarshal(fileContent, rawRunnerConfigRoot)
	if err2 != nil {
		fmt.Println("runner:", err)
		os.Exit(1)
	}

	return rawRunnerConfigRoot
}

// temp: move me to when converting raw -> sanitized config
func _validateRequired(rawConfigRoot *rawconfig.ConfigRoot) {
	rootDir := strings.TrimSpace(rawConfigRoot.Global.RootDirectory)
	if rootDir == "" {
		fmt.Println("config: [global] root_directory required.")
		os.Exit(1)
	}

	dir, err := os.MkdirTemp(rootDir, "tempe_*")
	if err != nil {
		fmt.Printf("config: unable to write to root_directory '%s'\n", rootDir)
		os.Exit(1)
	}
	os.RemoveAll(dir)

	nodeID := strings.TrimSpace(rawConfigRoot.Node.ID)
	if nodeID == "" {
		fmt.Println("config: [node] id required.")
		os.Exit(1)
	}

	isIllegal, invalidChar := sanitation.Code_ContainsIllegalChar(nodeID)
	if isIllegal {
		fmt.Printf("config: [node] id contains illegal char '%s'. Legal chars: %s\n", invalidChar, sanitation.CODE__LEGAL_CHARS)
		os.Exit(1)
	}
}
