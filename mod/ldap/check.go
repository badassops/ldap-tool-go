// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package ldap

// ** check functions: user and group **

func (c *Connection) CheckUser(userName string) int {
  return c.GetUser(userName, false)
}

func (c *Connection) CheckGroup(groupName string) int {
  recordsCount, _ := c.GetGroup(groupName)
  return recordsCount
}

func (c *Connection) CheckGroupID(groupID int) (bool, string) {
  for _, groupMap := range c.Config.GroupValues.GroupsMap {
    if groupMap.Gid == groupID {
      return true, groupMap.Name
    }
  }
  return false, "no-found"
}

func (c *Connection) CheckProtectedGroup(groupName string) bool {
  var groupsName []string
  groupsName = append(groupsName, c.Config.GroupValues.Groups...)
  groupsName = append(groupsName, c.Config.GroupValues.SpecialGroups...)

  for _, protectGroup := range groupsName {
    if groupName == protectGroup {
      return true
    }
  }
  return false
}
