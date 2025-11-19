package cli

import (
	ufc "github.com/urfave/cli/v2"
)

var flagOSUser = &ufc.StringFlag{
	Name: "os_user",
	Aliases: []string{"user"},
	Required: true,
}

var flagOSUserGroup = &ufc.StringFlag{
	Name: "os_user_group",
	Aliases: []string{"group"},
	Required: true,
}

var flagBinPath = &ufc.StringFlag{
	Name: "bin",
	Value: "noname",
	Usage: "Path to executable bin `FILE`",
	Required: true,
}

var flagNodeConfig = &ufc.StringFlag{
	Name: "config",
	Aliases: []string{"c"},
	Value: "conf.ini",
	Usage: "Load configuration from `FILE`",
	// Category: "required",
	Required: true,
	// EnvVars: []string{"NODE_CONFIG"},
	// FilePath: "/etc/mysql/password",
	// Action: func(ctx *ufc.Context, v int) (error) {
	// 	if v >= 65536 {
	// 		return fmt.Errorf("Flag port value %v out of range[0-65535]", v)
	// 	}
	// 	return nil
	// },
}
