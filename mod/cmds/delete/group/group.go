// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package delete

import (
  "bufio"
  "fmt"
  "os"
  "strings"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
)

var (
  valueEntered string
  continueDelete string
)

func deleteGroup(c *l.Connection) {
  reader := bufio.NewReader(os.Stdin)
  u.PrintYellow(fmt.Sprintf("\tEnter group name to be use: "))
  valueEntered, _ = reader.ReadString('\n')
  valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
  if valueEntered == "" {
    u.PrintRed(fmt.Sprintf("\n\tNo users was given aborting...\n"))
    return
    }

  if cnt := c.CheckGroup(valueEntered) ; cnt == 0 {
    u.PrintRed(fmt.Sprintf("\n\tGiven group %s does not found, aborting...\n", valueEntered))
    return
  }

  if c.CheckProtectedGroup(valueEntered) {
    u.PrintRed(fmt.Sprintf("\n\tGiven group %s is protected and can not be deleted, aborting...\n\n",
      valueEntered))
    return
  }

  c.Group["groupName"] = valueEntered
  c.Group["cn"] = fmt.Sprintf("cn=%s,%s", valueEntered, c.Config.ServerValues.GroupDN) 

  u.PrintRed(fmt.Sprintf("\n\tGiven group %s will be delete, this can not be undo!\n", valueEntered))
  u.PrintYellow(fmt.Sprintf("\tContinue (default to N)? [y/n]: "))
  continueDelete, _ = reader.ReadString('\n')
  continueDelete = strings.ToLower(strings.TrimSuffix(continueDelete, "\n"))
  if u.GetYN(continueDelete, false) == true {
    if !c.DeleteGroup() {
      u.PrintRed(fmt.Sprintf("\n\tFailed to delete the group %s, check the log file\n", valueEntered))
    } else {
      u.PrintGreen(fmt.Sprintf("\n\tGiven group %s has been deleted\n", valueEntered))
    }
  } else {
    u.PrintBlue(fmt.Sprintf("\n\tDeletion of the group %s cancelled\n", valueEntered))
  }
  return
}

func Delete(c *l.Connection) {
  u.PrintHeader(u.Purple, "Delete Group", true)
  deleteGroup(c)
  u.PrintLine(u.Purple)
}
