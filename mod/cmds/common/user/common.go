// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package common

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "strings"

  u "badassops.ldap/utils"
  "badassops.ldap/ldap"
)

var (
  displayFields = []string{"uid", "givenName", "sn", "cn", "displayName",
    "gecos", "uidNumber", "gidNumber", "departmentNumber",
    "mail", "homeDirectory", "loginShell", "userPassword",
    "shadowWarning", "shadowMax", "sshPublicKey"}

)

func printUserRecord(c *ldap.Connection, userName string) {
  // the values are in days so we need to multiple by 86400
  value, _ := strconv.ParseInt(c.User.Field["shadowLastChange"], 10, 64)
   _, passChanged := u.GetReadableEpoch(value * 86400)

  value, _ = strconv.ParseInt(c.User.Field["shadowExpire"], 10, 64)
  _, passExpired := u.GetReadableEpoch(value * 86400)

  u.PrintLine(u.Purple)
  for _, field := range displayFields {
    u.PrintColor(u.Cyan, fmt.Sprintf("\t%s: %s\n", field, c.User.Field[field]))
  }

  u.PrintLine(u.Purple)
  u.PrintColor(u.Purple, fmt.Sprintf("\tUser %s groups:\n", userName))
  for _, group := range c.User.Groups {
    u.PrintColor(u.Cyan, fmt.Sprintf("\tdn: %s\n", group))
  }

  u.PrintLine(u.Purple)
  u.PrintColor(u.Purple, fmt.Sprintf("\tUser %s password information\n", userName))
  u.PrintColor(u.Cyan, fmt.Sprintf("\tPassword last changed on %s\n", passChanged))
  u.PrintColor(u.Red, fmt.Sprintf("\tPassword will expired on %s\n", passExpired))
}

func User(c *ldap.Connection, firstTime bool, showRecord bool) bool {
  reader := bufio.NewReader(os.Stdin)
  fmt.Printf("\tEnter user login name to be use: ")
  enterData, _ := reader.ReadString('\n')
  enterData = strings.TrimSuffix(enterData, "\n")

  if enterData == "" {
    u.PrintColor(u.Red, fmt.Sprintf("\n\tNo users was given aborting...\n"))
    return false
  }

  if firstTime {
    fmt.Printf("\tUse wildcard (default to N)? [y/n]: ")
    wildCard, _ := reader.ReadString('\n')
    wildCard = strings.TrimSuffix(wildCard, "\n")
    if u.GetYN(wildCard, false) == true {
      enterData = "*" + enterData + "*"
      c.SearchUser(enterData)
      fmt.Printf("\n\tSelect the userid from the above list:\n")
      User(c, false, showRecord)
      return true
    }
  }

  if c.GetUser(enterData, false) == 0 {
    u.PrintColor(u.Red, fmt.Sprintf("\n\tUser %s was not found, aborting...\n", enterData))
    return false
  }
  if showRecord {
    printUserRecord(c, enterData)
  }
  return true
}
