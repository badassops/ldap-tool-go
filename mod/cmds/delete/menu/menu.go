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

  deleteGroup "badassops.ldap/cmds/delete/group"
  deleteSudo "badassops.ldap/cmds/delete/sudo"
  deleteUser "badassops.ldap/cmds/delete/user"
)

func DeleteMenu(c *l.Connection) {
  reader := bufio.NewReader(os.Stdin)

  u.PrintHeader(u.Purple, "Delete", true)
  fmt.Printf("\tDelete (%s)ser, (%s)roup, (%s)udo role or (%s)uit?\n\t(default to User)? choice: ",
    u.CreateColorMsg(u.Green, "U"),
    u.CreateColorMsg(u.Green, "G"),
    u.CreateColorMsg(u.Blue, "S"),
    u.CreateColorMsg(u.Red,   "Q"),
  )

  choice, _ := reader.ReadString('\n')
  choice = strings.TrimSuffix(choice, "\n")
  switch strings.ToLower(choice) {
    case "user",  "u": deleteUser.Delete(c)
    case "group", "g": deleteGroup.Delete(c)
    case "sudo", "s": deleteSudo.Delete(c)
    case "quit",  "q":
        u.PrintRed("\n\tOperation cancelled\n")
        u.PrintLine(u.Purple)
        break
    default: deleteUser.Delete(c)
  }
}
