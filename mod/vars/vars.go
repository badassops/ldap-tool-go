//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
package vars

import (
  "fmt"
  "os"
  "path"
  "strconv"
  "time"
)

type Record struct {
  Value    string // default value from the configuration
  Prompt   string
  NoEmpty  bool
  UseValue bool
}

type UserRecord struct {
  Field    map[string]string
  Groups   []string
}

type Log struct {
  LogsDir       string
  LogFile       string
  LogMaxSize    int
  LogMaxBackups int
  LogMaxAge     int
}

type ModifyRecord struct {
  Field   map[string]string
  AddList []string
  DelList []string
}

var (
  MyVersion  = "0.1.0"
  now      = time.Now()
  MyProgname  = path.Base(os.Args[0])
  myAuthor  = "Luc Suryo"
  myCopyright = "Copyright 2019 - " + strconv.Itoa(now.Year()) + " ©Badassops LLC"
  myLicense  = "License 3-Clause BSD, https://opensource.org/licenses/BSD-3-Clause ♥"
  myEmail    = "<luc@badassops.com>"
  MyInfo = fmt.Sprintf("%s (version %s)\n%s\n%s\nWritten by %s %s\n",
    MyProgname, MyVersion, myCopyright, myLicense, myAuthor, myEmail)
  MyDescription = "Simple script to manage LDAP users"

  // ldap logs
  Logs Log

  // the ldap fields
  Fields []string

  // user record in the ldap server
  User UserRecord

  // group record in the ldap server
  Group map[string]string

  // ldap record to be use for create or modify a user
  Template map[string]Record

  // ldap record to be use for create or modify a group
  GroupTemplate map[string]Record

  // formodify user
  ModRecord ModifyRecord

  // we sets these under variable
  LogsDir       = "/var/log/ldap-go"
  LogFile       = fmt.Sprintf("%s.log", MyProgname)
  LogMaxSize    = 128 // megabytes
  LogMaxBackups = 14  // 14 files
  LogMaxAge     = 14  // 14 days
)
