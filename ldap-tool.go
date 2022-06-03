// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version	:	0.1
//

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	//"time"

	// local
	"badassops.ldap/configurator"
	"badassops.ldap/initializer"
	"badassops.ldap/ldap"
	"badassops.ldap/logs"
	"badassops.ldap/vars"

	// on github
	"github.com/badassops/packages-go/is"
	"github.com/badassops/packages-go/lock"
	"github.com/badassops/packages-go/print"
	//"github.com/badassops/packages-go/spinner"

	// the menus
	searchMenu "badassops.ldap/cmds/search/menu"
	// createMenu "badassops.ldap/cmds/create/menu"
	deleteMenu "badassops.ldap/cmds/delete/menu"
	modifyMenu "badassops.ldap/cmds/modify/menu"
	// limit	  "badassops.ldap/cmds/limit"
)

func main() {
	LockPid := os.Getpid()
	progName, _ := os.Executable()
	progBase := filepath.Base(progName)

	i := is.New()
	p := print.New()
	// s := spinner.New(100)

	config := configurator.Configurator()

	// get given parameters
	config.InitializeArgs()

	// get the configuration
	config.InitializeConfigs()

	// initialize the user data dictionary
	initializer.Init(config)

	// make sure the configuration file has the proper settings
	runningUser, _ := i.IsRunningUser()
	if !i.IsInList(config.AuthValues.AllowUsers, runningUser) {
		p.PrintRed(fmt.Sprintf("The program has to be run as these user(s): %s or use sudo, aborting..\n",
			strings.Join(config.AuthValues.AllowUsers[:], ", ")))
		os.Exit(0)
	}
	ownerInfo, ownerOK := i.IsFileOwner(config.ConfigFile, config.AuthValues.AllowUsers)
	if !ownerOK {
		p.PrintRed(fmt.Sprintf("%s,\nAborting..\n", ownerInfo))
		os.Exit(0)
	}
	permInfo, permOK := i.IsFilePermission(config.ConfigFile, config.AuthValues.AllowMods)
	if !permOK {
		p.PrintRed(fmt.Sprintf("%s,\nAborting..\n", permInfo))
		os.Exit(0)
	}

	// only if the given server was enabled
	if config.ServerValues.Enabled == false {
		p.PrintRed(fmt.Sprintf("The given server %s is not enabled, aborting..\n", config.Server))
		os.Exit(0)
	}

	// go s.Run()
	// initialize the logger system
	LogConfig := &logs.LogConfig{
		LogsDir:       config.LogValues.LogsDir,
		LogFile:       config.LogValues.LogFile,
		LogMaxSize:    config.LogValues.LogMaxSize,
		LogMaxBackups: config.LogValues.LogMaxBackups,
		LogMaxAge:     config.LogValues.LogMaxAge,
	}

	logs.InitLogs(LogConfig)
	logs.Log("System all clear", "INFO")

	// create the lock file to prevent an other script is running/started
	l := lock.New(config.DefaultValues.LockFile)
	config.LockPID = LockPid
	if config.Cmd != "search" {
		// check lock file; lock file should not exist
		if _, fileExist, _ := i.IsExist(config.DefaultValues.LockFile, "file"); fileExist {
			lockPid, _ := l.LockGetPid()
			if progRunning, _ := i.IsRunning(progBase, lockPid); progRunning {
				// s.Stop()
				p.PrintRed(fmt.Sprintf("\nError there is already a process %s running, aborting...\n", progBase))
				os.Exit(0)
			}
		}
		// save to create new or overwrite the lock file
		if err := l.LockIt(LockPid); err != nil {
			// s.Stop()
			p.PrintRed(fmt.Sprintf("\nError creating the lock file, error %s, aborting..\n", err.Error()))
			os.Exit(0)
		}
	}

	// start the LDAP connection
	conn := ldap.New(config)

	// time.Sleep(1 * time.Second)
	// s.Stop()

	if config.ServerValues.ReadOnly == true {
		p.PrintRed(fmt.Sprintf("\tThe server %s is set to be ready only.\n\tOnly the Search options is available...\n",
			config.ServerValues.Server))
		p.PrintGreen("\tPress enter to continue to search: ")
		fmt.Scanln()
		config.Cmd = "search"
	}

	// semi-hardcoded
	// if config.ServerValues.Admin != "cn=admin," + config.ServerValues.BaseDN {
	//   switch config.Cmd {
	// 	case "search":
	// 	  limit.UserRecord() //conn)
	// 	case "modify":
	// 	  limit.ModifyUserPasswordSSHKey() //conn)
	// 	default:
	// 	  u.PrintRed("\n\tThis command is only available for admin...\n\n")
	//   }
	// } else {
	switch config.Cmd {
	case "search":
		searchMenu.SearchMenu(conn)
	// case "create":
	//   createMenu.CreateMenu() //conn)
	case "modify":
		modifyMenu.ModifyMenu(conn)
	case "delete":
		deleteMenu.DeleteMenu(conn)
	}
	// }

	if config.Cmd != "search" {
		l.LockRelease()
	}

	p.TheEnd()
	fmt.Printf("\t%s\n", p.PrintLine(vars.Purple, 50))
	logs.Log("System Normal shutdown", "INFO")
	os.Exit(0)
}
