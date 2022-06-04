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

// set an user's ldap passwod
func (c *Connection) SetPassword() bool {
	// once the record is create we need to hash the password
	passwordReq := ldapv3.NewPasswordModifyRequest(
		v.WorkRecord.DN, "", v.WorkRecord.Fields["userPassword"])
	if _, err := c.Conn.PasswordModify(passwordReq); err != nil {
		msg = fmt.Sprintf("Failed setting password for the user %s, error %s", v.WorkRecord.ID, err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("Successfully setting the password for for user %s to %s",
		v.WorkRecord.ID, v.WorkRecord.Fields["userPassword"])
	l.Log(msg, "INFO")
	return true
}
