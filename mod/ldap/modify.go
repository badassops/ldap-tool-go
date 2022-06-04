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

// remove an user from a group
func (c *Connection) RemoveFromGroups() bool {
	// posix uses user name
	// groupOfNames uses user's full dn
	removeReq := ldapv3.NewModifyRequest(v.WorkRecord.DN)
	removeReq.Delete(v.WorkRecord.MemberType, []string{v.WorkRecord.ID})
	if err := c.Conn.Modify(removeReq); err != nil {
		msg = fmt.Sprintf("Error removing the user %s from group %s, error %s",
			v.WorkRecord.ID, v.WorkRecord.DN, err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The %s %s has been modify", v.WorkRecord.ID, v.WorkRecord.DN)
	l.Log(msg, "INFO")
	return true
}

// add an user to a group
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

// modify an user ldap record
func (c *Connection) ModifyUser() bool {
	var passChanged bool = false
	modifyRecord := ldapv3.NewModifyRequest(v.WorkRecord.DN)
	for fieldName, fieldValue := range v.WorkRecord.Fields {
		if fieldName != "userPassword" {
			modifyRecord.Replace(fieldName, []string{fieldValue})
		}
		if fieldName == "userPassword" {
			passChanged = true
		}
	}
	if err := c.Conn.Modify(modifyRecord); err != nil {
		msg = fmt.Sprintf("Error modifying the user %s, error %s",
			v.WorkRecord.ID, err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	if passChanged {
		return c.SetPassword()
	}
	return true
}

// delete a sudo rule
func (c *Connection) DeleteSudoRule() bool {
	delSudoRule := ldapv3.NewModifyRequest(v.WorkRecord.DN)
	for fieldName, _ := range v.WorkRecord.SudoDelList {
		for _, value := range v.WorkRecord.SudoDelList[fieldName] {
			delSudoRule.Delete(fieldName, []string{value})
		}
	}
	if err := c.Conn.Modify(delSudoRule); err != nil {
		msg = fmt.Sprintf("Error deleting some of the entries of sudo rule %s, error %s",
			v.WorkRecord.ID, err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The sudo rule %s entries has been modified", v.WorkRecord.ID)
	l.Log(msg, "INFO")
	return true
}

// add a sudo rule
func (c *Connection) AddSudoRule() bool {
	addSudoRule := ldapv3.NewModifyRequest(v.WorkRecord.DN)
	for fieldName, _ := range v.WorkRecord.SudoAddList {
		for _, value := range v.WorkRecord.SudoAddList[fieldName] {
			addSudoRule.Add(fieldName, []string{value})
		}
	}
	if err := c.Conn.Modify(addSudoRule); err != nil {
		msg = fmt.Sprintf("Error adding some of the entries of sudo rule %s, error %s",
			v.WorkRecord.ID, err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The sudo rule %s entries has been modified", v.WorkRecord.ID)
	l.Log(msg, "INFO")
	return true
}
