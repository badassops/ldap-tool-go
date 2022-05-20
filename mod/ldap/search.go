// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package ldap

import (
  "fmt"

  u "badassops.ldap/utils"

  ldapv3 "gopkg.in/ldap.v2"
)

// ** search functions: user(s) and group(s) **

func (c *Connection) printSearchUser(searchBase string, baseInfo bool) int {
  attributes = []string{}
  records, recordsCount = c.search(searchBase, attributes)
  if recordsCount > 0 {
    for idx, entry = range records.Entries {
      switch u.FuncName() {
        case "SearchUser":
          u.PrintBlue(fmt.Sprintf("\tuid: %s\n",
            records.Entries[idx].GetAttributeValue("uid")))

        case "SearchUsers":
          u.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
          if baseInfo {
            userBaseInfo := fmt.Sprintf("\tFull namae: %s %s\t\tdepartmentNumber %s\n\n",
              records.Entries[idx].GetAttributeValue("givenName"),
              records.Entries[idx].GetAttributeValue("sn"),
                records.Entries[idx].GetAttributeValue("departmentNumber"))
            u.PrintCyan(userBaseInfo)
          }
      }
    }
    if u.FuncName() == "SearchUsers" {
      u.PrintYellow(fmt.Sprintf("\n\tTotal records: %d \n", recordsCount))
    }
  }
  return recordsCount
}

func (c *Connection) printSearchGroup(searchBase, groupName, groupType string, baseInfo bool) int {
  attributes = []string{}
  switch groupType{
    case "posix":
      searchBase = fmt.Sprintf("(&(objectClass=posixGroup)(cn=%s))", groupName)
      memberField = "memberUid"
    case "groupOfNames":
      searchBase = fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", groupName)
      memberField = "member"
  }
  records, recordsCount = c.search(searchBase, attributes)
  if recordsCount > 0 {
    for idx, entry = range records.Entries {
      switch u.FuncName() {
        case "SearchGroup":
          if baseInfo {
            u.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
            u.PrintBlue(fmt.Sprintf("\tcn: %s\n",
            records.Entries[idx].GetAttributeValue("cn")))
            if groupType == "posix" {
              u.PrintCyan(fmt.Sprintf("\tgidNumber: %s\n", entry.GetAttributeValue("gidNumber")))
            }
            for _, member := range entry.GetAttributeValues(memberField) {
              u.PrintCyan(fmt.Sprintf("\t%s: %s\n", memberField, member))
            }
          } else {
            u.PrintBlue(fmt.Sprintf("\tcn: %s\n", records.Entries[idx].GetAttributeValue("cn")))
          }

        case "SearchGroups":
          u.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
          u.PrintBlue(fmt.Sprintf("\tcn: %s\n", records.Entries[idx].GetAttributeValue("cn")))
          if groupType == "posix" {
            u.PrintCyan(fmt.Sprintf("\tgidNumber: %s\n", entry.GetAttributeValue("gidNumber")))
          }
          for _, member := range entry.GetAttributeValues(memberField) {
            u.PrintCyan(fmt.Sprintf("\t%s: %s\n", memberField, member))
          }
          u.PrintColor(u.Yellow,
            fmt.Sprintf("\tTotal members: %d \n", len(entry.GetAttributeValues(memberField))))
          fmt.Printf("\n")
      }
    }
    if u.FuncName() == "SearchGroups" {
      u.PrintYellow(fmt.Sprintf("\n\tTotal %s groups records: %d \n\n", groupType, recordsCount))
    }
  }
  return recordsCount
}

func (c *Connection) SearchUser(userName string) int {
  searchBase = fmt.Sprintf("(&(objectClass=inetOrgPerson)(uid=%s))", userName)
  return c.printSearchUser(searchBase, false)
}

func (c *Connection) SearchUsers(baseInfo bool) int {
  searchBase = fmt.Sprintf("(objectClass=person)")
  return c.printSearchUser(searchBase, baseInfo)
}

func (c *Connection) SearchGroup(groupName string, baseInfo bool) int {
  _, groupType = c.GetGroup(groupName)
  return c.printSearchGroup(searchBase, groupName, groupType, baseInfo)
}

func (c *Connection) SearchGroups() {
  groupTypes := []string{"posix", "groupOfNames"}
  baseInfo := false
  for _, groupType := range groupTypes {
    groupName := "*"
    c.printSearchGroup(searchBase, groupName, groupType, baseInfo)
  }
}

func (c *Connection) SearchUsersGroups(user string) []string {
  var groupsList []string
  searchBase = fmt.Sprintf("(&(objectClass=posixGroup))")
  attributes = []string{"cn", "memberUid"}
  records, _ = c.search(searchBase, attributes)
  for idx, entry := range records.Entries {
    for _, member := range entry.GetAttributeValues("memberUid") {
      if member == "user" {
        groupsList = append(groupsList, records.Entries[idx].GetAttributeValue("cn"))
      }
    }
  }
  return groupsList
}

func (c *Connection) search(searchBase string, searchAttribute []string) (*ldapv3.SearchResult, int) {
  searchRecords := ldapv3.NewSearchRequest(
    c.Config.ServerValues.BaseDN,
    ldapv3.ScopeWholeSubtree,
    ldapv3.NeverDerefAliases, 0, 0, false,
    searchBase,
    searchAttribute,
    nil,
  )
  sr, err := c.Conn.Search(searchRecords)
  if err != nil {
    c.Conn.Close()
    u.ReleaseIT(c.LockFile, c.LockPid)
    u.ExitWithMesssage(err.Error())
  }
  if len(sr.Entries) > 0 {
    return sr, len(sr.Entries)
  }
  return sr, 0
}
