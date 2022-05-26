// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package modify

import (
  "fmt"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
  v "badassops.ldap/vars"
  cs "badassops.ldap/cmds/common/sudo"
)

var (
  allowedField = []string{"sudoCommand", "sudoHost", "sudoOption",
    "sudoOrder", "sudoRunAsUser" }
)

func createSudoModRecord(c *l.Connection) {
  records, _ := c.GetSudoCN(v.ModRecord.Field["cn"])
  for _, entry := range records {
    u.PrintBlue(fmt.Sprintf("\tDN: %s\n", entry.DN))
    for _, attributes := range entry.Attributes {
      for _, value := range attributes.Values {
        if attributes.Name != "objectClass" {
          switch attributes.Name {
            case "sudoCommand":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))
              fmt.Printf("\tEnter %s value: \n", attributes.Name)

            case "sudoHost":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))
              fmt.Printf("\tEnter %s value: \n", attributes.Name)

            case "sudoOption":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))
              fmt.Printf("\tEnter %s value: \n", attributes.Name)

            case "sudoOrder":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))
              fmt.Printf("\tEnter %s value: \n", attributes.Name)

            case "sudoRunAsUser":
              u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value ))
              fmt.Printf("\tEnter %s value: \n", attributes.Name)

          }
        }
      }
    }
  }
}

func Modify(c *l.Connection) {
  u.PrintHeader(u.Purple, "Modify Sudo rule", true)
  if cs.Sudo(c, true, false) {
    if u.InList(c.Config.SudoValues.ExcludeSudo, v.ModRecord.Field["cn"]) {
      u.PrintRed(fmt.Sprintf("\t\tThe sudo rule %s can not be modified\n", v.ModRecord.Field["cn"]))
    } else {
      v.ModRecord.Field["dn"] = fmt.Sprintf("%s,%s",
        v.ModRecord.Field["cn"], c.Config.SudoValues.SudoersBase)
      //fmt.Printf("\n\tModifying the sudo rule %s\n", v.ModRecord.Field["dn"])
      createSudoModRecord(c)
    }
  }
  u.PrintLine(u.Purple)
}
