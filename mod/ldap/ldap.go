// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package ldap

import (
  "fmt"
  "net"
  "time"

  "crypto/tls"

  "badassops.ldap/vars"
  u "badassops.ldap/utils"
  "badassops.ldap/configurator"
  "badassops.ldap/logs"

  ldapv3 "gopkg.in/ldap.v2"
)

type (
  Connection struct {
    Conn     *ldapv3.Conn
    User     vars.UserRecord
    Group    map[string]string
    Config   *configurator.Config
    LockFile string
    LockPid  int
  }
)

var (
  // use in other file with same package name

  // these are the objectClasses needed for a user record
  userObjectClasses = []string{"top", "person",
    "organizationalPerson", "inetOrgPerson",
    "posixAccount", "shadowAccount", "ldapPublicKey"}

  msg          string
  searchBase   string
  groupType    string
  memberField  string
  groupTypes   []string
  recordsCount int
  idx           int

  groupObjectClasses = []string{"groupOfNames"}
  attributes         = []string{}

  records *ldapv3.SearchResult
  entry   *ldapv3.Entry
  addReq  *ldapv3.AddRequest
  delReq  *ldapv3.DelRequest
)

// function to initialize a user record
func New(config *configurator.Config) *Connection {
  // set variable for the ldap connection
  var ppolicy *ldapv3.ControlBeheraPasswordPolicy

  // check if we can search the server, timeout set to 15 seconds
  timeout := 15 * time.Second
  dialConn, err := net.DialTimeout("tcp", net.JoinHostPort(config.ServerValues.Server, "389"), timeout)
  if err != nil {
    u.ReleaseIT(config.DefaultValues.LockFile, config.LockPID)
    u.ExitWithMesssage(err.Error() + "\n")
  }
  dialConn.Close()

  ServerConn, err := ldapv3.Dial("tcp", fmt.Sprintf("%s:%d", config.ServerValues.Server ,389))
  if err != nil {
    u.ReleaseIT(config.DefaultValues.LockFile, config.LockPID)
    u.ExitWithMesssage(err.Error())
  }

  // now we need to reconnect with TLS
  if config.ServerValues.TLS {
    err := ServerConn.StartTLS(&tls.Config{InsecureSkipVerify: true})
    if err != nil {
      u.ReleaseIT(config.DefaultValues.LockFile, config.LockPID)
      u.ExitIfError(err)
    }
  }

  // setup control
  controls := []ldapv3.Control{}
  controls = append(controls, ldapv3.NewControlBeheraPasswordPolicy())

  // bind to the ldap server
  bindRequest := ldapv3.NewSimpleBindRequest(config.ServerValues.Admin, config.ServerValues.AdminPass, controls)
  request, err := ServerConn.SimpleBind(bindRequest)
  ppolicyControl := ldapv3.FindControl(request.Controls, ldapv3.ControlTypeBeheraPasswordPolicy)
  if ppolicyControl != nil {
    ppolicy = ppolicyControl.(*ldapv3.ControlBeheraPasswordPolicy)
   }
  if err != nil {
    errStr := "ERROR: Cannot bind: " + err.Error()
    if ppolicy != nil && ppolicy.Error >= 0 {
      errStr += ":" + ppolicy.ErrorString
    }
    u.ReleaseIT(config.DefaultValues.LockFile, config.LockPID)
    u.ExitWithMesssage(errStr)
  }

   // debug
  if config.Debug {
    logs.Log(fmt.Sprintf("Server : %s", config.ServerValues.Server), "DEBUG")
    logs.Log(fmt.Sprintf("__ BaseDN      : %s", config.ServerValues.BaseDN), "DEBUG")
    logs.Log(fmt.Sprintf("__ Admin       : %s", config.ServerValues.Admin), "DEBUG")
    logs.Log(fmt.Sprintf("__ AdminPass   : %s", config.ServerValues.AdminPass), "DEBUG")
    logs.Log(fmt.Sprintf("__ UserDN      : %s", config.ServerValues.UserDN), "DEBUG")
    logs.Log(fmt.Sprintf("__ GroupDN     : %s", config.ServerValues.GroupDN), "DEBUG")
    logs.Log(fmt.Sprintf("__ EmailDomain : %s", config.ServerValues.EmailDomain), "DEBUG")
    logs.Log(fmt.Sprintf("__ TLS         : %t", config.ServerValues.TLS), "DEBUG")
    logs.Log(fmt.Sprintf("__ isEnabled   : %t", config.ServerValues.Enabled), "DEBUG")
  }

  // the rest of the values will be filled during the process
  return &Connection {
    Conn:     ServerConn,
    Config:   config,
    User:     vars.User,
    Group:    vars.Group,
    LockFile: config.DefaultValues.LockFile,
    LockPid:  config.LockPID,
  }
}
