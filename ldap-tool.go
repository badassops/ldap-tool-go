// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"badassops.ldap/utils"
	"badassops.ldap/initializer"

	"badassops.ldap/constants"
	"badassops.ldap/logs"
	"badassops.ldap/configurator"
	"badassops.ldap/ldap"
	"badassops.ldap/cmds"
)

func main() {
	LockPid  := os.Getpid()
	progName, _ := os.Executable()
	info := filepath.Base(progName)

	// initialize the user data dictionary
	initializer.Init()

	// get given parameters
	config := configurator.Configurator()
	config.InitializeArgs()

	// these are hardcoded!
	if ok := utils.IsUser("luc"); !ok {
		utils.PrintColor(constants.Red, "The program has to be run as root or use sudo, aborting..\n")
		os.Exit(0)
	}
	if ok := utils.CheckFileSettings(config.ConfigFile, "luc", []string{"0400", "0600"}); !ok {
		utils.PrintColor(constants.Red, "Aborting..\n")
		os.Exit(0)
	}

	// get the configuration
	config.InitializeConfigs()

	// only if the given server was enabled
	if config.Enabled == false {
		utils.PrintColor(constants.Red, fmt.Sprintf("The given server %s is not enabled, aborting..\n", config.Env ))
		os.Exit(0)
	}

	// create the lock file to prevent an other script is running/started
	config.LockPID = LockPid

	// initialize the logger system
	LogConfig := &logs.LogConfig {
        LogsDir:        config.LogsDir,
        LogFile:        config.LogFile,
        LogMaxSize:     config.LogMaxSize,
        LogMaxBackups:  config.LogMaxBackups,
        LogMaxAge:      config.LogMaxAge,
    }

	logs.InitLogs(LogConfig)
	logs.Log("System all clear", "INFO")

	// create lock all initializing has been done
	utils.LockIT(config.LockFile, LockPid, info)

	// add a new ldap record
	conn := ldap.New(config)

	switch config.Cmd {
		case "create":	// cmds.Create(conn)
		case "modify":	// cmds.Modify(conn)
		case "delete":	// cmds.Delete(conn)
		case "search":	cmds.Search(conn, "user")
		case "group":	cmds.Search(conn, "group")
		case "admin":	cmds.Search(conn, "admin")
		case "users":	cmds.Search(conn, "users")
		case "groups":	cmds.Search(conn, "groups")
		case "admins":	cmds.Search(conn, "admins")
	}

	utils.ReleaseIT(config.LockFile, LockPid)
	os.Exit(0)
}
