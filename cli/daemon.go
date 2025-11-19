package cli

import (
	"fmt"
	
	ufc "github.com/urfave/cli/v2"

	"noname001/provisioning"
)

func daemonRootCommand() (*ufc.Command) {
	return &ufc.Command{
		Name: "daemon",
		Usage: "Daemon related commands",
		Flags: []ufc.Flag{},
		Subcommands: []*ufc.Command{
			_daemonSubCommandInstall(),
			// _daemonSubCommandGenerate(),
		},
	}
}

func _daemonSubCommandInstall() (*ufc.Command) {
	return &ufc.Command{
		Name: "install",
		Usage: "Install as systemd service",
		Category: "daemon",
		Flags: []ufc.Flag{
			flagOSUser, flagOSUserGroup,
			flagBinPath, flagNodeConfig,
		},
		Before: func(ctx *ufc.Context) (error) {
			// TODO: os checks & systemctl checks

			// https://stackoverflow.com/questions/29733575/how-to-find-the-user-that-executed-a-program-as-root-using-golang
			// need not check, better to just always show warning
			provisioning.PrintRootPrivilegeRequirement()

			return nil
		},
		Action: func(ctx *ufc.Context) (error) {
			var err error
			// TODO: on error, the exit status should only be 1 line and is exit 1

			// TODO: stop all shits before doing any install
			//       either make install idempotent or disable install, but use upgrade command

			// TODO: input safety check
			dat := map[string]string{
				"os_user": ctx.String("os_user"),
				"os_user_group": ctx.String("os_user_group"),

				"bin": ctx.String("bin"),
				"config": ctx.String("config"),
			}

			err = provisioning.InstallOSDependencies()
			if err != nil {
				fmt.Println(err)
				return err
			}

			err = provisioning.InstallAsSystemdService(dat)
			if err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		},
		After: func(ctx *ufc.Context) (error) {
			return nil
		},
	}
}
