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

var (
	msg string
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

func (c *Connection) RemoveFromGroups() bool {
	// posix uses user name
	// groupOfNames uses user's full dn
	delReq := ldapv3.NewModifyRequest(v.WorkRecord.DN)
	delReq.Delete(v.WorkRecord.MemberType, []string{v.WorkRecord.ID})
	if err := c.Conn.Modify(delReq); err != nil {
		msg = fmt.Sprintf("Error removing the user %s from group %s, error %s",
			v.WorkRecord.ID, v.WorkRecord.DN, err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The %s %s has been modify", v.WorkRecord.ID, v.WorkRecord.DN)
	l.Log(msg, "INFO")
	return true
}

func (c *Connection) AddToGroup() bool {
	// posix uses user name
	// groupOfNames uses user's full dn
	addReq := ldapv3.NewModifyRequest(v.WorkRecord.DN)
	addReq.Add(v.WorkRecord.MemberType, []string{v.WorkRecord.ID})
	if err := c.Conn.Modify(addReq); err != nil {
		msg = fmt.Sprintf("Error removing the user %s from group %s, error %s",
			v.WorkRecord.ID, v.WorkRecord.DN, err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The %s %s has been modify", v.WorkRecord.ID, v.WorkRecord.DN)
	l.Log(msg, "INFO")
	return true
}

// func (c *Connection) ModifyUser() bool {
//   var errored int = 0
//   var passChanged bool = false
//   modifyRecord := ldapv3.NewModifyRequest(c.User.Field["dn"])
//   for fieldName, fieldValue := range v.ModRecord.Field {
//     if fieldName != "userPassword" {
//       modifyRecord.Replace(fieldName, []string{fieldValue})
//     }
//     if fieldName == "userPassword" {
//        c.User.Field["userPassword"] = fieldValue
//        passChanged = true
//     }
//   }
//
//   if !c.modify(c.User.Field["uid"], "user", modifyRecord) {
//     errored++
//   }
//
//   if passChanged {
//     if !c.setPassword() {
//       errored++
//     }
//   }
//
//   if errored != 0 {
//     return false
//   }
//   // modify only if the user modification was successfuly
//   c.ModifyUserGroup()
//   return true
// }

// func (c *Connection) ModifyUserGroup() {
//   // group type is always groupOfNames
//   var modUser, groupName, groupCN string
//   modUser = fmt.Sprintf("uid=%s,%s", c.User.Field["uid"], c.Config.ServerValues.UserDN)
//
//   for _, groupName = range v.ModRecord.AddList {
//     groupCN = fmt.Sprintf("cn=%s,%s", groupName, c.Config.ServerValues.GroupDN)
//     modifyMember := ldapv3.NewModifyRequest(groupCN)
//     modifyMember.Add("member", []string{modUser})
//     c.modify(groupName, "group member (add)", modifyMember)
//   }
//
//   for _, groupName = range v.ModRecord.DelList {
//     groupCN = fmt.Sprintf("cn=%s,%s", groupName, c.Config.ServerValues.GroupDN)
//     modifyMember := ldapv3.NewModifyRequest(groupCN)
//     modifyMember.Delete("member", []string{modUser})
//     c.modify(groupName, "group member (remove)", modifyMember)
//   }
// }

// func (c *Connection) ModifySudoRule() bool {
//   var modifiedOK bool = true
//   if len(v.ModSudo.DelList) > 0 {
//     modifyDelRule := ldapv3.NewModifyRequest(v.ModSudo.DN)
//     for fieldName, _ := range v.ModSudo.DelList {
//       for _, value := range v.ModSudo.DelList[fieldName] {
//         modifyDelRule.Delete(fieldName, []string{value})
//       }
//     }
//     if !c.modify(v.ModSudo.DN, "delete rule ", modifyDelRule) {
//       modifiedOK = modifiedOK && false
//     }
//   }
//
//   if len(v.ModSudo.AddList) > 0 {
//     modifyAddRule := ldapv3.NewModifyRequest(v.ModSudo.DN)
//     for fieldName, value := range v.ModSudo.AddList {
//       modifyAddRule.Add(fieldName, value)
//     }
//     if !c.modify(v.ModSudo.DN, "add rule ", modifyAddRule) {
//       modifiedOK = modifiedOK && false
//     }
//   }
//   return modifiedOK
// }
