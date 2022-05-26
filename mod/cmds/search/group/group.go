// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package group

import (

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
  cg "badassops.ldap/cmds/common/group"
)


func Group(c *l.Connection) {
  u.PrintHeader(u.Purple, "Search Group", true)
  cg.Group(c, true)
  u.PrintLine(u.Purple)
}

func Groups(c *l.Connection) {
  u.PrintHeader(u.Purple, "Search Groups", true)
  u.PrintPurple("\n\t  _________ all group and the members __________\n")
  c.SearchGroups()
  u.PrintLine(u.Purple)
}
