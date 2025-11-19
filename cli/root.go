package cli

import (
	"fmt"
	"os"

	ufc "github.com/urfave/cli/v2"

	"noname001/version"
)

// https://cli.urfave.org/v2/examples/full-api-example/

func init() {
	ufc.VersionPrinter = func(cCtx *ufc.Context) {
		fmt.Printf("%s %s\n", cCtx.App.Name, cCtx.App.Version)
	}
}

func Root() {
	rootFlags := []ufc.Flag{}

	rootCommands := []*ufc.Command{
		// TODO: dependency command
		nodeRootCommand(),

		daemonRootCommand(),

		// TODO: emergency deactivation command
	}

	root := &ufc.App{
		Name: version.BIN,
		Usage: version.BIN,
		Version: version.FullVersion(),

		Flags: rootFlags,
		Commands: rootCommands,

		EnableBashCompletion: false,
		HideHelp: false,
		HideVersion: false,

		// CommandNotFound: func(cCtx *cli.Context, command string) {
        //     fmt.Fprintf(cCtx.App.Writer, "Thar be no %q here.\n", command)
        // },
		// Metadata: map[string]interface{}{
        //     "layers":          "many",
        //     "explicable":      false,
        //     "whatever-values": 19.99,
        // },
	}

	// do err handling in their respective commands, due to different mode / usage
	rootErr := root.Run(os.Args)
	_ = rootErr
}
