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

  modifyUser "badassops.ldap/cmds/modify/user"
  modifyGroup "badassops.ldap/cmds/modify/group"
  modifySudo "badassops.ldap/cmds/modify/sudo"
)

func ModifyMenu(c *l.Connection) {
  reader := bufio.NewReader(os.Stdin)

  u.PrintHeader(u.Purple, "Modify", true)
  fmt.Printf("\tModify (%s)ser, (%s)roup, (%s)udo rule or (%s)uit?\n\t(default to User)? choice: ",
    u.CreateColorMsg(u.Green, "U"),
    u.CreateColorMsg(u.Green, "G"),
    u.CreateColorMsg(u.Blue,  "S"),
    u.CreateColorMsg(u.Red,   "Q"),
  )

  choice, _ := reader.ReadString('\n')
  choice = strings.TrimSuffix(choice, "\n")
  switch strings.ToLower(choice) {
    case "user",  "u": modifyUser.Modify(c)
    case "group", "g": modifyGroup.Modify(c)
    case "sudo",  "s": modifySudo.Modify(c)
    case "quit",  "q":
        u.PrintRed("\n\tOperation cancelled\n")
        u.PrintLine(u.Purple)
        break
    default: modifyUser.Modify(c)
  }
}
