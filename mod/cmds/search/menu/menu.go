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

  searchUser "badassops.ldap/cmds/search/user"
  searchGroup "badassops.ldap/cmds/search/group"
  searchSudo "badassops.ldap/cmds/search/sudo"
)

func SearchMenu(c *l.Connection) {
  reader := bufio.NewReader(os.Stdin)

  u.PrintHeader(u.Purple, "Search", true)
  fmt.Printf("\tSearch (%s)ser, (%s)ll Users, (%s)roup, all Group(%s)\n",
    u.CreateColorMsg(u.Green, "U"),
    u.CreateColorMsg(u.Green, "A"),
    u.CreateColorMsg(u.Green, "G"),
    u.CreateColorMsg(u.Green, "S"),
  )
  fmt.Printf("\t\t(%s)sudo role, (%s)all sudos role or (%s)uit?\n\t(default to User)? choice: ",
    u.CreateColorMsg(u.Blue,  "X"),
    u.CreateColorMsg(u.Blue,  "Z"),
    u.CreateColorMsg(u.Red,   "Q"),
  )

  choice, _ := reader.ReadString('\n')
  choice = strings.TrimSuffix(choice, "\n")
  switch strings.ToLower(choice) {
    case "user",   "u": searchUser.User(c)
    case "users",  "a": searchUser.Users(c)
    case "group",  "g": searchGroup.Group(c)
    case "groups", "s": searchGroup.Groups(c)
    case "sudo",   "x": searchSudo.Sudo(c)
    case "sudos",  "z": searchSudo.Sudos(c)
    case "quit",   "q":
        u.PrintRed("\n\tOperation cancelled\n")
        u.PrintLine(u.Purple)
        break
    default: searchUser.User(c)
  }
}
