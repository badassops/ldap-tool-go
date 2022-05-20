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

  l "badassops.ldap/logs"

  ldapv3 "gopkg.in/ldap.v2"
)

// delete functions: user and group

func (c *Connection) deleteRecord(recordId, recordType string, request *ldapv3.DelRequest) bool {
  if err := c.Conn.Del(request); err != nil {
    msg = fmt.Sprintf("Error deleting the %s %s, error %s", recordType, recordId, err.Error())
    l.Log(msg, "ERROR")
    return false
  }
  msg = fmt.Sprintf("The %s %s has been deleted", recordType, recordId)
  l.Log(msg, "INFO")
  return true
}

func (c *Connection) DeleteUser() bool {
  delReq := ldapv3.NewDelRequest(c.User.Field["dn"], []ldapv3.Control{})
  return c.deleteRecord(c.User.Field["uid"], "user", delReq)
}

func (c *Connection) DeleteGroup() bool {
  delReq := ldapv3.NewDelRequest(c.Group["cn"], []ldapv3.Control{})
  return c.deleteRecord(c.Group["groupName"], "group", delReq)
}

func (c *Connection) RemoveFromGroups() {
  userGroups := c.GetUserGroups("posfix")
  for _, userGroup := range userGroups {
    delReq := ldapv3.NewModifyRequest(userGroup)
    delReq.Delete("memberUid", []string{c.User.Field["uid"]})
    c.modify(c.User.Field["uid"], "member of group" , delReq)
  }
}
