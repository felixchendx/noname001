package provisioning

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"noname001/filesystem"
)

// TODO: bin upgrade and downgrade
// symlink bin to /usr/local/bin

// TODO: unit template add checks and stuffs
// TODO: test restart and kills and 
const (
	UNIT_NAME    = "noname.service"
	INSTALL_PATH = "/etc/systemd/system"

	SYSTEMD_UNIT_TEMPLATE = `# Comment: ...

[Unit]
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User={{ $.os_user }}
Group={{ $.os_user_group }}
ExecStart={{ $.abs_bin_path }} node start --config {{ $.abs_config_path }}
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target

`

	INSTALLATION_SCRIPT = `#!/bin/bash

declare CURR_USER=$(whoami)
declare HOME_DIR=$(echo ~)
`
)


// func GenerateInstallationScript() (error) {

// 	return nil
// }

func PrintRootPrivilegeRequirement() {
	fmt.Println("======================================================")
	fmt.Println("=== Warning! This command requires root privilege. ===")
	fmt.Println("======================================================")
}

func InstallAsSystemdService(params map[string]string) (error) {
	var (
		err       error
		logBorder string = "\n=== Installing as systemd service: %s ===\n"
	)

	fmt.Printf(logBorder, "start")
	defer fmt.Printf(logBorder, "end")

	absBinPath, err := filepath.Abs(params["bin"])
	if err != nil {
		return err
	}
	absConfigPath, err := filepath.Abs(params["config"])
	if err != nil {
		return err
	}

	// final dat
	finalDat := map[string]string{
		"os_user"        : params["os_user"],
		"os_user_group"  : params["os_user_group"],
		"abs_bin_path"   : absBinPath,
		"abs_config_path": absConfigPath,
	}

	tmpl, err := template.New("SystemdUnitTemplate").Parse(SYSTEMD_UNIT_TEMPLATE)
	if err != nil {
		return err
	}

	sb := new(strings.Builder)
	if err := tmpl.Execute(sb, finalDat); err != nil {
		return err
	}


	// simple check for clean / existing service installation
	_, err = os.Open(filepath.Join(INSTALL_PATH, UNIT_NAME))
	if err == nil {
		fmt.Println("updating service installation...")

		// stop before tweaking da service
		err = _stopAndDisableNonameService()
		if err != nil {
			return err
		}
	} else {
		fmt.Println("new service installation...")

		// nothing to stop
	}

	writeErr := os.WriteFile(
		filepath.Join(INSTALL_PATH, UNIT_NAME),
		[]byte(sb.String()),
		filesystem.DEFAULT_FILE_PERMISSION,
	)
	if writeErr != nil {
		return writeErr
	}

	// restart service
	err = _enableAndStartNonameService()
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	_ = _statNonameService()

	return nil
}

func InstallOSDependencies() (error) {
	var (
		err       error
		logBorder string = "\n=== Installing OS dependencies: %s ===\n"
	)

	fmt.Printf(logBorder, "start")
	defer fmt.Printf(logBorder, "end")

	cmdBash := exec.Command("/usr/bin/env", "/bin/bash")
	cmdBashIn, _ := cmdBash.StdinPipe()
	cmdBash.Stdout = os.Stdout
	cmdBash.Stderr = os.Stderr

	err = cmdBash.Start()
	if err != nil { return err }
	
	// TEMP
	fmt.Fprintf(cmdBashIn, "apt-get update -y\n")
	fmt.Fprintf(cmdBashIn, "apt-get install ffmpeg libzmq3-dev -y\n") // python3.10-venv
	cmdBashIn.Close()

	err = cmdBash.Wait()
	if err != nil { return err }

	return nil
}

func _stopAndDisableNonameService() (err error) {
	fmt.Println("")

	cmdBash := exec.Command("/usr/bin/env", "/bin/bash")
	cmdBashIn, _ := cmdBash.StdinPipe()
	cmdBash.Stdout, cmdBash.Stderr = os.Stdout, os.Stderr

	err = cmdBash.Start()
	if err != nil { return err }

	fmt.Fprintf(cmdBashIn, "systemctl stop %s\n", UNIT_NAME)
	fmt.Fprintf(cmdBashIn, "systemctl disable %s\n", UNIT_NAME)
	cmdBashIn.Close()

	err = cmdBash.Wait()
	if err != nil { return err }

	return nil
}

func _enableAndStartNonameService() (err error) {
	fmt.Println("")

	cmdBash := exec.Command("/usr/bin/env", "/bin/bash")
	cmdBashIn, _ := cmdBash.StdinPipe()
	cmdBash.Stdout, cmdBash.Stderr = os.Stdout, os.Stderr

	err = cmdBash.Start()
	if err != nil { return err }

	fmt.Fprintf(cmdBashIn, "systemctl enable %s\n", UNIT_NAME)
	fmt.Fprintf(cmdBashIn, "systemctl start %s\n", UNIT_NAME)
	cmdBashIn.Close()

	err = cmdBash.Wait()
	if err != nil { return err }

	return nil
}

func _statNonameService() (err error) {
	fmt.Println("")

	cmdBash := exec.Command("/usr/bin/env", "/bin/bash")
	cmdBashIn, _ := cmdBash.StdinPipe()
	cmdBash.Stdout, cmdBash.Stderr = os.Stdout, os.Stderr

	err = cmdBash.Start()
	if err != nil { return err }

	// TODO: sensitive info in proc
	fmt.Fprintf(cmdBashIn, "systemctl status %s --no-page\n", UNIT_NAME)
	cmdBashIn.Close()

	err = cmdBash.Wait()
	if err != nil { return err }

	return nil
}
