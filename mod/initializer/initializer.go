// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package initializer

import (
	"badassops.ldap/consts"
	"badassops.ldap/vars"
)

func Init() {
	// ldap fields that will be used
	vars.Fields	= []string{"uid", "givenName", "sn", "cn", "displayName",
		"gecos", "uidNumber", "gidNumber", "departmentNumber",
		"mail", "homeDirectory", "loginShell", "userPassword",
		"shadowLastChange", "shadowExpire", "shadowWarning", "shadowMax",
		"sshPublicKey"}

	vars.Logs.LogsDir		= vars.LogsDir
	vars.Logs.LogFile		= vars.LogFile
	vars.Logs.LogMaxSize	= vars.LogMaxSize
	vars.Logs.LogMaxBackups	= vars.LogMaxBackups
	vars.Logs.LogMaxAge		= vars.LogMaxAge

	vars.User.Strings	= make(map[string]vars.StringRecord)
	vars.User.Ints		= make(map[string]vars.IntRecord)
	vars.User.Groups	= []string{}

	for _, name := range vars.Fields {
		// NOTE shadowLastChange, shadowExpire is a string and will need to
		//		be converted to in64 during operation
		switch name {
			case "shadowWarning": vars.User.Ints[name] =
						vars.IntRecord{Value: consts.ShadowWarning, Changed: false}
			case "shadowMax": vars.User.Ints[name] =
						vars.IntRecord{Value: consts.ShadowMax, Changed: false}
			case "uidNumber", "gidNumber": vars.User.Ints[name] = vars.IntRecord{Changed: false}
			default:
				vars.User.Strings[name] = vars.StringRecord{Value: "", Changed: false}
		}
	}
}
