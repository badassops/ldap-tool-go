// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package user

import (
  "bufio"
  "fmt"
  "os"
  "strings"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
  cu "badassops.ldap/cmds/common/user"
)

var (
  valueEntered string
)

func User(c *l.Connection) {
  u.PrintHeader(u.Purple, "Search User", true)
  cu.User(c, true, true)
  u.PrintLine(u.Purple)
}

func Users(c *l.Connection) {
  u.PrintHeader(u.Purple, "Search Users", true)
  fmt.Printf("\tPrint full name and department (default to N)? [y/n]: ")
  reader := bufio.NewReader(os.Stdin)
  valueEntered, _ = reader.ReadString('\n')
  valueEntered = strings.TrimSuffix(valueEntered, "\n")
  u.PrintLine(u.Purple)
  if u.GetYN(valueEntered, false) == true {
    c.SearchUsers(true)
  } else {
    c.SearchUsers(false)
  }
  u.PrintLine(u.Purple)
}
