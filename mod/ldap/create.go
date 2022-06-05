//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package ldap

import (
	"fmt"

	l "badassops.ldap/logs"
	v "badassops.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

// create a ldap user
//func (c *Connection) AddUser() bool {
//	newUserReq := ldapv3.NewAddRequest(v.WorkRecord.Fields["dn"]
//	newUserReq.Attribute("objectClass", []string{v.WorkRecord.Fields["objectClass"]})
//	for _, field := range v.Fields {
//		if field != "groups" {
//			newUserReq.Attribute(field, []string{c.User.Field[field]})
//		}
//	}
// we ignore errors adding group
//c.addUserTogroupOfNamesGroup()
//c.addUserToPosixGroup()
// set password
//  return c.setPassword()
//}

// create a ldap group
func (c *Connection) CreateGroup() bool {
	newGroupReq := ldapv3.NewAddRequest(v.WorkRecord.Fields["dn"])
	newGroupReq.Attribute("objectClass", []string{v.WorkRecord.Fields["objectClass"]})
	newGroupReq.Attribute("cn", []string{v.WorkRecord.Fields["cn"]})
	if v.WorkRecord.Fields["objectClass"] == "posixGroup" {
		newGroupReq.Attribute("gidNumber", []string{v.WorkRecord.Fields["gidNumber"]})
	}
	if v.WorkRecord.Fields["objectClass"] == "groupOfNames" {
		newGroupReq.Attribute("member", []string{v.WorkRecord.Fields["member"]})
	}
	if err := c.Conn.Add(newGroupReq); err != nil {
		msg = fmt.Sprintf("Error creating the group %s error %s",
			v.WorkRecord.Fields["cn"], err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The group %s has been created", v.WorkRecord.Fields["cn"])
	l.Log(msg, "INFO")
	return true
}

// create a ldap sudo rule
func (c *Connection) CreateSudoRule() bool {
	newSudoRuleReq := ldapv3.NewAddRequest(v.WorkRecord.Fields["dn"])
	newSudoRuleReq.Attribute("objectClass", []string{v.WorkRecord.Fields["objectClass"]})
	for _, field := range v.SudoFields {
		if len(v.WorkRecord.Fields[field]) > 0 {
			newSudoRuleReq.Attribute(field, []string{v.WorkRecord.Fields[field]})
		}
	}
	if err := c.Conn.Add(newSudoRuleReq); err != nil {
		msg = fmt.Sprintf("Error creating the sudo rule %s error %s",
			v.WorkRecord.Fields["cn"], err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The sudo rule %s has been created", v.WorkRecord.Fields["cn"])
	l.Log(msg, "INFO")
	return true
}
