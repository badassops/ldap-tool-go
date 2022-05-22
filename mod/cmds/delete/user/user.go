// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package create

import (
  "bufio"
  "fmt"
  "os"
  "strings"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
)

var (
  valueEntered   string
  continueDelete string
)

func deleteUser(c *l.Connection) {
  reader := bufio.NewReader(os.Stdin)
  u.PrintYellow(fmt.Sprintf("\tEnter userid (login name) to be use: "))
  valueEntered, _ = reader.ReadString('\n')
  valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
  if valueEntered == "" {
    u.PrintRed(fmt.Sprintf("\n\tNo users was given aborting...\n"))
    return
    }

  if cnt := c.CheckUser(valueEntered); cnt == 0 {
    u.PrintRed(fmt.Sprintf("\n\tGiven user %s doen not exist, aborting...\n\n", valueEntered))
    return
  }

  c.User.Field["dn"] = fmt.Sprintf("uid=%s,%s", valueEntered, c.Config.ServerValues.UserDN)
  c.User.Field["uid"] = valueEntered

  u.PrintRed(fmt.Sprintf("\n\tGiven user %s will be delete, this can not be undo!\n", valueEntered))
  u.PrintYellow(fmt.Sprintf("\tContinue (default to N)? [y/n]: "))
  continueDelete, _ = reader.ReadString('\n')
  continueDelete = strings.ToLower(strings.TrimSuffix(continueDelete, "\n"))
  if u.GetYN(continueDelete, false) == true {
    if !c.DeleteUser() {
      u.PrintRed(fmt.Sprintf("\n\tFailed to delete the user %s, check the log file\n", c.User.Field["uid"]))
    } else {
      u.PrintGreen(fmt.Sprintf("\n\tGiven user %s has been deleted\n", valueEntered))
    }
    // ignore errors
    c.RemoveFromGroups()
  } else {
    u.PrintBlue(fmt.Sprintf("\n\tDeletion of the user %s cancelled\n", valueEntered))
  }
  return
}

func Delete(c *l.Connection) {
  u.PrintHeader(u.Purple, "Delete User", true)
  deleteUser(c)
  u.PrintLine(u.Purple)
}
