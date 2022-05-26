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
  "strings"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
  v "badassops.ldap/vars"
)


// Need to be global since the function is called recursive
var (
  groupType string
  enterData string
)

func Group(c *l.Connection, firstTime bool) bool {
  reader := bufio.NewReader(os.Stdin)
  fmt.Printf("\tEnters the group name to be use: ")
  enterData, _ = reader.ReadString('\n')
  enterData = strings.TrimSuffix(enterData, "\n")

  if enterData == "" {
    u.PrintRed(fmt.Sprintf("\n\tNo group was given aborting...\n"))
    if firstTime {
      return false
    } else {
      // need to break the recursive
      u.ReleaseIT(c.LockFile, c.LockPid)
      u.PrintLine(u.Purple)
      os.Exit(1)
    }
  }

  if firstTime {
    fmt.Printf("\tUse wildcard (default to N)? [y/n]: ")
    wildCard, _ := reader.ReadString('\n')
    wildCard = strings.TrimSuffix(wildCard, "\n")
    if u.GetYN(wildCard, false) == true {
      enterData = "*" + enterData + "*"
      c.SearchGroup(enterData, false)
      fmt.Printf("\n\tSelect the group name from the above list:\n")
      Group(c, false)
    }
  } else {
    // from recursive
    return true
  }

  u.PrintLine(u.Purple)
  if cnt := c.SearchGroup(enterData, true) ; cnt == 0 {
    u.PrintRed(fmt.Sprintf("\n\tGroup %s was not found, aborting...\n", enterData))
    return false
  }
  v.ModRecord.Field["groupName"] = enterData
  return true
}
