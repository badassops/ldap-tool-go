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
  v "badassops.ldap/vars"
  ldapv3 "gopkg.in/ldap.v2"
)

func (c *Connection) modify(recordId, recordType string, request *ldapv3.ModifyRequest) bool {
    if err := c.Conn.Modify(request); err != nil {
        msg = fmt.Sprintf("Error modify the %s %s, error %s", recordType, recordId, err.Error())
        l.Log(msg, "ERROR")
        return false
    }
    msg = fmt.Sprintf("The %s %s has been modify", recordType, recordId)
    l.Log(msg, "INFO")
    return true
}

func (c *Connection) ModifyUser() bool {
  var errored int = 0
  var passChanged bool = false
  modifyRecord := ldapv3.NewModifyRequest(c.User.Field["dn"])
  for fieldName, fieldValue := range v.ModRecord.Field {
    if fieldName != "userPassword" {
      modifyRecord.Replace(fieldName, []string{fieldValue})
    }
    if fieldName == "userPassword" {
       c.User.Field["userPassword"] = fieldValue
       passChanged = true
    }
  }

  if !c.modify(c.User.Field["uid"], "user", modifyRecord) {
    errored++
  }

  if passChanged {
    if !c.setPassword() {
      errored++
    }
  }

  if errored != 0 {
    return false
  }
  // modify only if the user modification was successfuly
  c.ModifyUserGroup()
  return true
}

func (c *Connection) ModifyUserGroup() {
  // group type is always groupOfNames
  var modUser, groupName, groupCN string
  modUser = fmt.Sprintf("uid=%s,%s", c.User.Field["uid"], c.Config.ServerValues.UserDN)

  for _, groupName = range v.ModRecord.AddList {
    groupCN = fmt.Sprintf("cn=%s,%s", groupName, c.Config.ServerValues.GroupDN)
    modifyMember := ldapv3.NewModifyRequest(groupCN)
    modifyMember.Add("member", []string{modUser})
    c.modify(groupName, "group member (add)", modifyMember)
  }

  for _, groupName = range v.ModRecord.DelList {
    groupCN = fmt.Sprintf("cn=%s,%s", groupName, c.Config.ServerValues.GroupDN)
    modifyMember := ldapv3.NewModifyRequest(groupCN)
    modifyMember.Delete("member", []string{modUser})
    c.modify(groupName, "group member (remove)", modifyMember)
  }
}


func (c *Connection) ModifyGroupMember() bool {
    groupCN := fmt.Sprintf("cn=%s,%s", v.ModRecord.Field["groupName"], c.Config.ServerValues.GroupDN)
    state, groupType := c.GetGroupType(v.ModRecord.Field["groupName"])
    if !state {
      return false
    }

    if groupType == "groupOfNames" {
      memberField = "member"
    }
    if groupType == "posix" {
      memberField = "memberUid"
    }

  if len(v.ModRecord.AddList) != 0 || len (v.ModRecord.DelList) != 0 {
    modifyGroup := ldapv3.NewModifyRequest(groupCN)
    for _, addMember := range v.ModRecord.AddList {
      modifyGroup.Add(memberField, []string{addMember})
    }
    for _, delMember := range v.ModRecord.DelList {
      modifyGroup.Delete(memberField, []string{delMember})
    }
    if !c.modify(v.ModRecord.Field["groupName"], "modify member", modifyGroup) {
      return false
	}
  }
  return true
}
