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

  l "badassops.ldap/logs"
  v "badassops.ldap/vars"

  ldapv3 "gopkg.in/ldap.v2"
)

// create functions: user and group

func (c *Connection) add(recordId, recordType string, request *ldapv3.AddRequest) bool {
  if err := c.Conn.Add(request); err != nil {
    msg = fmt.Sprintf("Error adding the %s %s, error %s", recordType, recordId, err.Error())
    l.Log(msg, "ERROR")
    return false
  }
  msg = fmt.Sprintf("The %s %s has been added", recordType, recordId)
  l.Log(msg, "INFO")
  return true
}

func (c *Connection) AddUser() bool {
  addReq := ldapv3.NewAddRequest(c.User.Field["dn"])
  addReq.Attribute("objectClass", userObjectClasses)
  for _, field := range v.Fields {
    if field != "groups" {
      addReq.Attribute(field, []string{c.User.Field[field]})
    }
  }

  if !c.add(c.User.Field["uid"], "user", addReq) {
    return false
  }

  // we ignore errors adding group
  c.addUserTogroupOfNamesGroup()
  c.addUserToPosixGroup()

  // set password
  return c.setPassword()
}

func (c *Connection) AddGroup() bool {
  addReq := ldapv3.NewAddRequest(c.Group["cn"])
  addReq.Attribute("objectClass", []string{c.Group["objectClass"]})
  addReq.Attribute("cn", []string{c.Group["groupName"]})

  if c.Group["groupType"] == "posix" {
    addReq.Attribute("gidNumber", []string{c.Group["gidNumber"]})
  }

  if c.Group["groupType"] == "groupOfNames" {
    addReq.Attribute("member", []string{c.Group["member"]})
  }

  return c.add(c.Group["groupName"], "group", addReq)
}

func (c *Connection) addUserTogroupOfNamesGroup() bool {
  // adding the user to the groups
  var errorCnt int = 0
  for _, group := range c.User.Groups {
    groupCN := fmt.Sprintf("cn=%s,%s", group, c.Config.ServerValues.GroupDN)
    addtoGroup := ldapv3.NewModifyRequest(groupCN)
    addtoGroup.Add("member", []string{c.User.Field["dn"]})
    if !c.modify("group", fmt.Sprintf("%s member to group %s", group, c.User.Field["uid"]), addtoGroup) {
      errorCnt++
    }
  }
  if errorCnt != 0 {
    return false
  }
  return true
}

func (c *Connection) addUserToPosixGroup() bool {
  var groupName string
  var groupCN string
  for _, mapValues := range c.Config.GroupValues.GroupsMap {
    if strconv.Itoa(mapValues.Gid) == c.User.Field["gidNumber"] {
      groupName = mapValues.Name
      groupCN = fmt.Sprintf("cn=%s,%s", groupName, c.Config.ServerValues.GroupDN)
      break
    }
  }
  if len(groupName) > 0 {
    addtoGroup := ldapv3.NewModifyRequest(groupCN)
    addtoGroup.Add("memberUid", []string{c.User.Field["uid"]})
    return c.modify("group", fmt.Sprintf("%s member to group %s", groupName, c.User.Field["uid"]), addtoGroup)
  }
  return false
}
