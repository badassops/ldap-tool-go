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
  v "badassops.ldap/vars"
  cs "badassops.ldap/cmds/common/sudo"
)

var (
  valueEntered string
  continueDelete string
)

func deleteSudoRule(c *l.Connection) {
  if cs.Sudo(c, true, false) {
    if u.InList(c.Config.SudoValues.ExcludeSudo, v.ModRecord.Field["cn"]) {
      u.PrintRed(fmt.Sprintf("\t\tThe sudo rule %s can not be deleted\n", v.ModRecord.Field["cn"]))
    } else {
      u.PrintRed(fmt.Sprintf("\n\tSudo rule %s will be delete, this can not be undo!\n",
        v.ModRecord.Field["cn"]))
      u.PrintYellow(fmt.Sprintf("\tContinue (default to N)? [y/n]: "))
      reader := bufio.NewReader(os.Stdin)
      continueDelete, _ = reader.ReadString('\n')
      continueDelete = strings.ToLower(strings.TrimSuffix(continueDelete, "\n"))
      if u.GetYN(continueDelete, false) == true {
        v.ModRecord.Field["dn"] = fmt.Sprintf("cn=%s,ou=%s",
          v.ModRecord.Field["cn"], c.Config.SudoValues.SudoersBase)
        if !c.DeleteSudoRule() {
          u.PrintRed(fmt.Sprintf("\n\tFailed to delete the sudo rule %s, check the log file\n", valueEntered))
        } else {
          u.PrintGreen(fmt.Sprintf("\n\tGiven sudo rule %s has been deleted\n", valueEntered))
        }
      } else {
        u.PrintBlue(fmt.Sprintf("\n\tDeletion of the sudo rule %s cancelled\n", valueEntered))
      }
    }
  }
}

func Delete(c *l.Connection) {
  u.PrintHeader(u.Purple, "Delete Sudo rule", true)
  deleteSudoRule(c)
  u.PrintLine(u.Purple)
}
