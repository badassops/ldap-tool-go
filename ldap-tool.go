// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	createUser "badassops.ldap/cmds/create/user"
	createGroup "badassops.ldap/cmds/create/group"
	modifyUser "badassops.ldap/cmds/modify/user"
	modifyGroup "badassops.ldap/cmds/modify/group"
	deleteUser "badassops.ldap/cmds/delete/user"
	deleteGroup "badassops.ldap/cmds/delete/group"
)

func main() {
	LockPid  := os.Getpid()
	progName, _ := os.Executable()
	info := filepath.Base(progName)

	// get given parameters
	config := configurator.Configurator()
	config.InitializeArgs()

	// get the configuration
	config.InitializeConfigs()

	// initialize the user data dictionary
	initializer.Init(config)

	// these are hardcoded!
	if ok := utils.IsUser("root"); !ok {
		utils.PrintColor(consts.Red, "The program has to be run as root or use sudo, aborting..\n")
		os.Exit(0)
	}
	if ok := utils.CheckFileSettings(config.ConfigFile, "root", []string{"0400", "0600"}); !ok {
		utils.PrintColor(consts.Red, "Aborting..\n")
		os.Exit(0)
	}

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

	reader := bufio.NewReader(os.Stdin)
	switch config.Cmd {
		case "search":
			utils.PrintHeader(consts.Purple, "Search", true)
			fmt.Printf("\tSearch (%s)ser, (%s)ll Users, (%s)roup, all Group(%s) or (%s)uit?\n\t(default to User)? choice: ",
				utils.CreateColorMsg(consts.Green, "U"),
				utils.CreateColorMsg(consts.Green, "A"),
				utils.CreateColorMsg(consts.Green, "G"),
				utils.CreateColorMsg(consts.Green, "S"),
				utils.CreateColorMsg(consts.Red, "Q"))

			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSuffix(choice, "\n")
			switch strings.ToLower(choice) {
				case "user", "u":	searchUser.User(conn)
				case "users", "a":	searchUser.Users(conn)
				case "group", "g":	searchGroup.Group(conn)
				case "groups", "s":	searchGroup.Groups(conn)
				case "quit", "q":
						utils.PrintColor(consts.Red, "\tOperation cancelled\n")
						break
				default: searchUser.User(conn)
			}

		case "create":
			utils.PrintHeader(consts.Purple, "Create", true)
			fmt.Printf("\tCreate (%s)ser, (%s)roup or (%s)uit?\n\t(default to User)? choice: ",
				utils.CreateColorMsg(consts.Green, "U"),
				utils.CreateColorMsg(consts.Green, "G"),
				utils.CreateColorMsg(consts.Red, "Q"))

			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSuffix(choice, "\n")
			switch strings.ToLower(choice) {
				case "user", "u":	createUser.Create(conn)
				case "group", "g":	createGroup.Create(conn)
				case "quit", "q":
						utils.PrintColor(consts.Red, "\tOperation cancelled\n")
						break
				default: createUser.Create(conn)
		}

		case "modify":
			utils.PrintHeader(consts.Purple, "Modify", true)
			fmt.Printf("\tModify (%s)ser, (%s)roup or (%s)uit?\n\t(default to User)? choice: ",
				utils.CreateColorMsg(consts.Green, "U"),
				utils.CreateColorMsg(consts.Green, "G"),
				utils.CreateColorMsg(consts.Red, "Q"))

			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSuffix(choice, "\n")
			switch strings.ToLower(choice) {
				case "user", "u":	modifyUser.Modify(conn)
				case "group", "g":	modifyGroup.Modify(conn)
				case "quit", "q":
						utils.PrintColor(consts.Red, "\tOperation cancelled\n")
						break
				default: modifyUser.Modify(conn)
		}

		case "delete":
			utils.PrintHeader(consts.Purple, "Delete", true)
			fmt.Printf("\tDelete (%s)ser, (%s)roup or (%s)uit?\n\t(default to User)? choice: ",
				utils.CreateColorMsg(consts.Green, "U"),
				utils.CreateColorMsg(consts.Green, "G"),
				utils.CreateColorMsg(consts.Red, "Q"))

			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSuffix(choice, "\n")
			switch strings.ToLower(choice) {
				case "user", "u":	deleteUser.Delete(conn)
				case "group", "g":	deleteGroup.Delete(conn)
				case "quit", "q":
						utils.PrintColor(consts.Red, "\tOperation cancelled\n")
						break
				default: deleteUser.Delete(conn)
		}
	}

	utils.ReleaseIT(config.DefaultValues.LockFile, LockPid)
	logs.Log("System Normal shutdown", "INFO")
	os.Exit(0)
}
