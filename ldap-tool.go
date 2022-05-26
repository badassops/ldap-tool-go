// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package main

import (
  "fmt"
  "os"
  "path/filepath"
  "strings"

  u "badassops.ldap/utils"

  "badassops.ldap/initializer"
  "badassops.ldap/configurator"
  "badassops.ldap/logs"
  "badassops.ldap/ldap"

  // the menus
  searchMenu "badassops.ldap/cmds/search/menu"
  createMenu "badassops.ldap/cmds/create/menu"
  modifyMenu "badassops.ldap/cmds/modify/menu"
  deleteMenu "badassops.ldap/cmds/delete/menu"
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

  // make sure the configuration file has the proper settings
  if !u.InList(config.AuthValues.AllowUsers, u.RunningUser()) {
    u.PrintRed(fmt.Sprintf("The program has to be run as these user(s): %s or use sudo, aborting..\n",
			strings.Join(config.AuthValues.AllowUsers[:], ", ")))
    os.Exit(0)
  }
  if !u.CheckFileSettings(config.ConfigFile, config.AuthValues.AllowUsers, config.AuthValues.AllowMods) {
    u.PrintRed("Aborting..\n")
    os.Exit(0)
  }

  // only if the given server was enabled
  if config.ServerValues.Enabled == false {
    u.PrintRed(fmt.Sprintf("The given server %s is not enabled, aborting..\n", config.Env ))
    os.Exit(0)
  }

  // create the lock file to prevent an other script is running/started
  config.LockPID = LockPid

  // initialize the logger system
  LogConfig := &logs.LogConfig {
    LogsDir:       config.LogValues.LogsDir,
    LogFile:       config.LogValues.LogFile,
    LogMaxSize:    config.LogValues.LogMaxSize,
    LogMaxBackups: config.LogValues.LogMaxBackups,
    LogMaxAge:     config.LogValues.LogMaxAge,
  }

  logs.InitLogs(LogConfig)
  logs.Log("System all clear", "INFO")

  // create lock all initializing has been done, but not for search
  if config.Cmd != "search" {
    u.LockIT(config.DefaultValues.LockFile, LockPid, info)
  }

  // add a new ldap record
  conn := ldap.New(config)

  if config.ServerValues.ReadOnly  == true {
    u.PrintRed(fmt.Sprintf("\tThe server %s is set to be ready only.\n\tOnly the Search options is available...\n",
        config.ServerValues.Server))
    u.PrintGreen(fmt.Sprintf("\tPress enter to continue to search: "))
    fmt.Scanln()
    config.Cmd = "search"
  }

  switch config.Cmd {
    case "search":
      searchMenu.SearchMenu(conn)

    case "create":
      createMenu.CreateMenu(conn)

    case "modify":
      modifyMenu.ModifyMenu(conn)

    case "delete":
      deleteMenu.DeleteMenu(conn)
  }

  if config.Cmd != "search" {
    u.ReleaseIT(config.DefaultValues.LockFile, LockPid)
  }

  u.TheEnd()
  u.PrintLine(u.Purple)
  logs.Log("System Normal shutdown", "INFO")
  os.Exit(0)
}
