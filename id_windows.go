//go:build windows
// +build windows

package machineid

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// machineID returns the key MachineGuid in registry `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Cryptography`.
// If there is an error running the commad an empty string is returned.
func machineID() (string, error) {
	if id, e := getSmbiosUUID(); e == nil {
		return id, nil
	}
	if id, e := deviceID(); e == nil {
		return id, nil
	}
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return "", err
	}
	return s, nil
}
func deviceID() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\SQMClient`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("MachineId")
	if err != nil {
		return "", err
	}
	return s, nil

}
func getSmbiosUUID() (uuid string, err error) {
	defer func() {
		panc := recover()
		if panc != nil {
			err = fmt.Errorf("getSmbiosUUID: %v", panc)
		}
	}()
	cmd := exec.Command("wmic", "csproduct", "get", "UUID")
	output, err := cmd.Output()
	flag := false
	if err == nil {
		uuid = strings.TrimRight(strings.Split(string(output), "\n")[1], "\r\n")
		if uuid == "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF" {
			err = fmt.Errorf("smbios return FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF")
			flag = false
		} else {
			flag = true
		}
	}
	if !flag {
		cmd = exec.Command("wmic", "bios", "get", "serialnumber")
		output, err = cmd.Output()
		if err == nil {
			uuid = strings.TrimRight(strings.Split(string(output), "\n")[1], "\r\n")
			flag = true
		}
	}
	return
}
