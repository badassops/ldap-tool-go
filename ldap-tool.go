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

	"badassops.ldap/consts"
	"badassops.ldap/utils"

	"badassops.ldap/initializer"
	"badassops.ldap/configurator"
	"badassops.ldap/logs"
	"badassops.ldap/ldap"

	// the search function
	searchUser "badassops.ldap/cmds/search/user"
	searchGroup "badassops.ldap/cmds/search/group"

	// the base functions ; create, modify and delete
	//createUser "badassops.ldap/cmds/create/user"
	//modifyUser "badassops.ldap/cmds/modify/user"
	//removeUser "badassops.ldap/cmds/remove/user"

	// for future version
	//createGroup "badassops.ldap/cmds/create/group"
	//modifyGroup "badassops.ldap/cmds/modify/group"
	//removeGroup "badassops.ldap/cmds/remove/group"
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
	if ok := utils.IsUser("root"); !ok {
		utils.PrintColor(consts.Red, "The program has to be run as root or use sudo, aborting..\n")
		os.Exit(0)
	}
	if ok := utils.CheckFileSettings(config.ConfigFile, "root", []string{"0400", "0600"}); !ok {
		utils.PrintColor(consts.Red, "Aborting..\n")
		os.Exit(0)
	}

	// get the configuration
	config.InitializeConfigs()

	// initialize the user record
	config.InitializeUserRecord()

	// only if the given server was enabled
	if config.ServerValues.Enabled == false {
		utils.PrintColor(consts.Red, fmt.Sprintf("The given server %s is not enabled, aborting..\n", config.Env ))
		os.Exit(0)
	}

	// create the lock file to prevent an other script is running/started
	config.LockPID = LockPid

	// initialize the logger system
	LogConfig := &logs.LogConfig {
        LogsDir:        config.LogValues.LogsDir,
        LogFile:        config.LogValues.LogFile,
        LogMaxSize:     config.LogValues.LogMaxSize,
        LogMaxBackups:  config.LogValues.LogMaxBackups,
        LogMaxAge:      config.LogValues.LogMaxAge,
    }

	logs.InitLogs(LogConfig)
	logs.Log("System all clear", "INFO")

	// create lock all initializing has been done
	utils.LockIT(config.DefaultValues.LockFile, LockPid, info)

	// add a new ldap record
	conn := ldap.New(config)

	switch config.Cmd {
	//	case "create":	create.Create(conn)
	//	case "modify":	// cmds.Modify(conn)
	//	case "delete":	// cmds.Delete(conn)
		case "search":	searchUser.User(conn)
		case "users":	searchUser.Users(conn)
		case "group":	searchGroup.Group(conn)
		case "groups":	searchGroup.Groups(conn)
	}

	utils.ReleaseIT(config.DefaultValues.LockFile, LockPid)
	os.Exit(0)
}
