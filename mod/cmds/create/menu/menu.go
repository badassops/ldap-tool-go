// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package menu

import (
  "bufio"
  "fmt"
  "os"
  "strings"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"

  createGroup "badassops.ldap/cmds/create/group"
  createSudo  "badassops.ldap/cmds/create/sudo"
  createUser "badassops.ldap/cmds/create/user"
)

func CreateMenu(c *l.Connection) {
  reader := bufio.NewReader(os.Stdin)

  u.PrintHeader(u.Purple, "Create", true)
  fmt.Printf("\tCreate (%s)ser, (%s)roup, (%s)sudo role or (%s)uit?\n\t(default to User)? choice: ",
    u.CreateColorMsg(u.Green, "U"),
    u.CreateColorMsg(u.Green, "G"),
    u.CreateColorMsg(u.Blue, "S"),
    u.CreateColorMsg(u.Red,   "Q"),
  )

  choice, _ := reader.ReadString('\n')
  choice = strings.TrimSuffix(choice, "\n")
  switch strings.ToLower(choice) {
    case "user",  "u": createUser.Create(c)
    case "group", "g": createGroup.Create(c)
    case "sudo",  "s": createSudo.Create(c)
    case "quit",  "q":
        u.PrintRed("\n\tOperation cancelled\n")
        u.PrintLine(u.Purple)
        break
    default: createUser.Create(c)
  }
}
