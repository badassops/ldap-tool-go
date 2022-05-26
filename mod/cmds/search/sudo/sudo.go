// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package sudo

import (
  "fmt"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
  cs "badassops.ldap/cmds/common/sudo"
)

func Sudo(c *l.Connection) {
  u.PrintHeader(u.Purple, "Search Sudo Rules", true)
  cs.Sudo(c, true, true)
  u.PrintLine(u.Purple)
}

func Sudos(c *l.Connection) {
  u.PrintHeader(u.Purple, "Search Sudo Rules", true)
  records := c.SearchSudoRoles()
  if records != nil {
    for _, entry := range records {
      u.PrintBlue(fmt.Sprintf("\tDN: %s\n", entry.DN))
      for _, attributes := range entry.Attributes {
        for _, value := range attributes.Values {
          if attributes.Name != "objectClass" {
            if attributes.Name == "cn" {
              if u.InList(c.Config.SudoValues.ExcludeSudo, value) {
                u.PrintYellow("\t\tThis entry can not be modified\n")
              }
            }
            u.PrintCyan(fmt.Sprintf("\t\t%s : %s \n", attributes.Name, value ))
          }
        }
      }
    fmt.Printf("\n")
    }
  }
  u.PrintLine(u.Purple)
}
