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
  "strconv"

  u "badassops.ldap/utils"
  v "badassops.ldap/vars"
  ldapv3 "gopkg.in/ldap.v2"
)

// ** get info functions: user and group **

func (c *Connection) GetUser(userName string, checkOnly bool) int {
  attributes = []string{}
  searchBase = fmt.Sprintf("(&(objectClass=inetOrgPerson)(uid=%s))", userName)
  records, recordsCount = c.search(searchBase, attributes)
  if recordsCount == 1 {
    if checkOnly {
      return recordsCount
    }
    c.User.Field["dn"] = records.Entries[0].DN
    for _, field := range v.Fields {
      c.User.Field[field] = records.Entries[0].GetAttributeValue(field)
    }
    c.GetUserGroups("groupOfNames")
  }
  return recordsCount
}

func (c *Connection) GetGroup(groupName string) (int, string) {
  attributes := []string{}
  groupTypes = []string{"posix", "groupOfNames"}
  for _, groupType := range groupTypes {
    switch groupType {
      case "posix":
        searchBase = fmt.Sprintf("(&(objectClass=posixGroup)(cn=%s))", groupName)
      case "groupOfNames":
        searchBase = fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", groupName)
    }
    _, recordsCount = c.search(searchBase, attributes)
    if recordsCount > 0 {
      return recordsCount, groupType
    }
  }
  return recordsCount, "unknown"
}

func (c *Connection) GetGroupsNameByType(groupType string) []string {
  var searchBase string
  attributes := []string{"cn"}
  groupNames := []string{}
  switch groupType{
    case "posix":
      searchBase = fmt.Sprintf("(&(objectClass=posixGroup))")
    case "groupOfNames":
      searchBase = fmt.Sprintf("(&(objectClass=groupOfNames))")
  }
  records, _ := c.search(searchBase, attributes)
  for _, entry := range records.Entries {
    for _, groupCN := range entry.GetAttributeValues("cn") {
      groupNames = append(groupNames, groupCN)
    }
  }
  return groupNames
}

func (c *Connection) GetUserGroups(groupType string) []string {
  var groupList []string
  switch groupType{
    case "posix":
      searchBase = fmt.Sprintf("(&(objectClass=posixGroup)(memberUid=%s))", c.User.Field["uid"])
    case "groupOfNames":
      searchBase = fmt.Sprintf("(&(objectClass=groupOfNames)(member=%s))", c.User.Field["dn"])
  }
  records, recordsCount := c.search(searchBase, []string{"dn"})
  if recordsCount > 0 {
    for _, entry := range records.Entries {
      if u.InList(groupList, entry.DN) == false {
        groupList = append(groupList, entry.DN)
      }
    }
  }
  return groupList
}

func (c *Connection) GetNextUID() int {
  var startUID = c.Config.DefaultValues.UidStart
  var uidValue int
  attributes := []string{"uidNumber"}
  searchBase := fmt.Sprintf("(objectClass=person)")
  records, _ := c.search(searchBase, attributes)
  for _, entry := range records.Entries {
      uidValue, _ = strconv.Atoi(entry.GetAttributeValue("uidNumber"))
      if uidValue > startUID {
        startUID = uidValue
      }
  }
  return startUID + 1
}

func (c *Connection) GetNextGID() int {
  var startGID = c.Config.DefaultValues.GidStart
  var uidValue int
  attributes := []string{"gidNumber"}
  searchBase := fmt.Sprintf("(objectClass=posixGroup)")
  records, _ := c.search(searchBase, attributes)
  for _, entry := range records.Entries {
      uidValue, _ = strconv.Atoi(entry.GetAttributeValue("gidNumber"))
      if uidValue == c.Config.DefaultValues.GroupId {
        // we skip this special gid
        continue
      }
      if uidValue > startGID {
        startGID = uidValue
      }
  }
  return startGID + 1
}

func (c *Connection) GetGroupType(groupName string) (bool, string) {
    var typeGroup string
    var cnt int
    if cnt, typeGroup = c.GetGroup(groupName); cnt == 0 {
        return false, "errored"
    }
    return true, typeGroup
}

func (c* Connection) GetSudoCN(sudoCN string) ([]*ldapv3.Entry, int) {
  searchBase = fmt.Sprintf("(&(objectClass=top)(objectClass=sudoRole)(cn=%s))", sudoCN)
  records, recordsCount = c.search(searchBase, attributes)
  if recordsCount > 0 {
    return records.Entries, recordsCount
  }
  return nil, 0
}
