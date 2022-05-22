// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package modify

import (
  "bufio"
  "fmt"
  "os"
  "strings"

  u "badassops.ldap/utils"
  l "badassops.ldap/ldap"
  v "badassops.ldap/vars"
  gc "badassops.ldap/cmds/common/group"
)

var (
  valueEntered string

  // list of the member(s) to be added or deleted
  addList []string
  delList []string
  modCount int = 0
)

func modifyGroup(c *l.Connection) bool {
  if !gc.Group(c, true) {
    return false
  }
  u.PrintRed(fmt.Sprintf("\n\tEnter the user(s) to be deleted, select from the list above, (default to skip)\n"))
  reader := bufio.NewReader(os.Stdin)
  for true {
    fmt.Printf("\tUser : ")
    valueEntered, _ = reader.ReadString('\n')
    valueEntered = strings.TrimSuffix(valueEntered, "\n")
    if valueEntered != "" {
      delList = append(delList, valueEntered)
    } else {
      break
    }
  }

  u.PrintGreen(fmt.Sprintf("\n\tEnter the user(s) to be added, (default to skip)\n"))
  for true {
    fmt.Printf("\tUser : ")
    valueEntered, _ := reader.ReadString('\n')
    valueEntered = strings.TrimSuffix(valueEntered, "\n")
    if valueEntered != "" {
      addList = append(addList, valueEntered)
    } else {
      break
    }
  }

  v.ModRecord.AddList = addList
  v.ModRecord.DelList = delList
  modCount = len(v.ModRecord.AddList) + len (v.ModRecord.DelList)
  if modCount == 0 {
    return false
  }
  return (c.ModifyGroupMember())
}

func Modify(c *l.Connection) {
  u.PrintHeader(u.Purple, "Modify Group", true)
  if !modifyGroup(c) {
      if modCount == 0 {
        u.PrintBlue(fmt.Sprintf("\n\tNo change, no modification was made for the group %s\n",
          v.ModRecord.Field["groupName"]))
      } else {
        u.PrintRed(fmt.Sprintf("\n\tFailed modify the group %s, check the log file\n",
          v.ModRecord.Field["groupName"]))
      }
  } else {
        u.PrintGreen(fmt.Sprintf("\n\tGroup %s modified successfully\n", v.ModRecord.Field["groupName"]))
  }
  u.PrintLine(u.Purple)
}
