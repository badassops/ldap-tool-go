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

func (c *Connection) Delete(recordId, recordType string) bool {
	delReq := ldapv3.NewDelRequest(v.WorkRecord.DN, []ldapv3.Control{})
	if err := c.Conn.Del(delReq); err != nil {
		msg := fmt.Sprintf("Error deleting the %s %s, error %s", recordType, recordId, err.Error())
		l.Log(msg, "ERROR")
		return false
	}
	msg := fmt.Sprintf("The %s %s has been deleted", recordType, recordId)
	l.Log(msg, "INFO")
	return true
}
