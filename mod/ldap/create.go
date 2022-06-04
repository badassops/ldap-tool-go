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

func (c *Connection) AddGroup() bool {
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
