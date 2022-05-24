// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package main

import (
  "bufio"
  "fmt"
  "os"
  "path/filepath"
  "strings"

  u "badassops.ldap/utils"

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

  reader := bufio.NewReader(os.Stdin)
  switch config.Cmd {
    case "search":
      u.PrintHeader(u.Purple, "Search", true)
      fmt.Printf("\tSearch (%s)ser, (%s)ll Users, (%s)roup, all Group(%s) or (%s)uit?\n\t(default to User)? choice: ",
        u.CreateColorMsg(u.Green, "U"),
        u.CreateColorMsg(u.Green, "A"),
        u.CreateColorMsg(u.Green, "G"),
        u.CreateColorMsg(u.Green, "S"),
        u.CreateColorMsg(u.Red,   "Q"),
      )

      choice, _ := reader.ReadString('\n')
      choice = strings.TrimSuffix(choice, "\n")
      switch strings.ToLower(choice) {
        case "user",   "u": searchUser.User(conn)
        case "users",  "a": searchUser.Users(conn)
        case "group",  "g": searchGroup.Group(conn)
        case "groups", "s": searchGroup.Groups(conn)
        case "quit",   "q":
            u.PrintRed("\n\tOperation cancelled\n")
            u.PrintLine(u.Purple)
            break
        default: searchUser.User(conn)
      }

    case "create":
      u.PrintHeader(u.Purple, "Create", true)
      fmt.Printf("\tCreate (%s)ser, (%s)roup or (%s)uit?\n\t(default to User)? choice: ",
        u.CreateColorMsg(u.Green, "U"),
        u.CreateColorMsg(u.Green, "G"),
        u.CreateColorMsg(u.Red,   "Q"),
      )

      choice, _ := reader.ReadString('\n')
      choice = strings.TrimSuffix(choice, "\n")
      switch strings.ToLower(choice) {
        case "user",  "u": createUser.Create(conn)
        case "group", "g": createGroup.Create(conn)
        case "quit",  "q":
            u.PrintRed("\n\tOperation cancelled\n")
            u.PrintLine(u.Purple)
            break
        default: createUser.Create(conn)
      }

    case "modify":
      u.PrintHeader(u.Purple, "Modify", true)
      fmt.Printf("\tModify (%s)ser, (%s)roup or (%s)uit?\n\t(default to User)? choice: ",
        u.CreateColorMsg(u.Green, "U"),
        u.CreateColorMsg(u.Green, "G"),
        u.CreateColorMsg(u.Red,   "Q"),
      )

      choice, _ := reader.ReadString('\n')
      choice = strings.TrimSuffix(choice, "\n")
      switch strings.ToLower(choice) {
        case "user",  "u": modifyUser.Modify(conn)
        case "group", "g": modifyGroup.Modify(conn)
        case "quit",  "q":
            u.PrintRed("\n\tOperation cancelled\n")
            u.PrintLine(u.Purple)
            break
        default: modifyUser.Modify(conn)
      }

    case "delete":
      u.PrintHeader(u.Purple, "Delete", true)
      fmt.Printf("\tDelete (%s)ser, (%s)roup or (%s)uit?\n\t(default to User)? choice: ",
        u.CreateColorMsg(u.Green, "U"),
        u.CreateColorMsg(u.Green, "G"),
        u.CreateColorMsg(u.Red,   "Q"),
      )

      choice, _ := reader.ReadString('\n')
      choice = strings.TrimSuffix(choice, "\n")
      switch strings.ToLower(choice) {
        case "user",  "u": deleteUser.Delete(conn)
        case "group", "g": deleteGroup.Delete(conn)
        case "quit",  "q":
            u.PrintRed("\n\tOperation cancelled\n")
            u.PrintLine(u.Purple)
            break
        default: deleteUser.Delete(conn)
      }
  }

  if config.Cmd != "search" {
    u.ReleaseIT(config.DefaultValues.LockFile, LockPid)
  }

  u.TheEnd()
  u.PrintLine(u.Purple)
  logs.Log("System Normal shutdown", "INFO")
  os.Exit(0)
}
