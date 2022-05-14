// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package initializer

import (
	"strconv"

	"badassops.ldap/consts"
	"badassops.ldap/vars"
	"badassops.ldap/utils"
	"badassops.ldap/configurator"
)

func Init(conf *configurator.Config) {
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
	vars.RecordFields	= make(map[int]vars.RecordField)

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

	vars.RecordFields[0] =
		vars.RecordField{ FieldName: "uid",
			Prompt: "Enter userid (login name) to be use: ",
			Default: "",
			NoEmpty: true,
			UseDefault: false,
		}

	vars.RecordFields[1] =
		vars.RecordField{ FieldName: "givenName",
			Prompt: "Enter First name: ",
			Default: "",
			NoEmpty: true,
			UseDefault: false,
		}

	vars.RecordFields[2] =
		vars.RecordField{ FieldName: "sn",
			Prompt: "Enter Last name: ",
			Default: "",
			NoEmpty: true,
			UseDefault: false,
		}

	vars.RecordFields[3] =
		vars.RecordField{ FieldName: "mail",
			Prompt: "Enter email: ",
			Default: "",
			NoEmpty: false,
			UseDefault: true,
		}

	vars.RecordFields[4] =
		vars.RecordField{ FieldName: "uidNumber",
			Prompt: "Enter user's UID: ",
			Default: "",
			NoEmpty: false,
			UseDefault: true,
		}

	vars.RecordFields[5] =
		vars.RecordField{ FieldName: "departmentNumber",
			Prompt: "Enter department: ",
			Default: conf.DefaultValues.GroupName,
			NoEmpty: false,
			UseDefault: true,
		}

	vars.RecordFields[6] =
		vars.RecordField{ FieldName: "loginShell",
			Prompt: "Enter shell: ",
			Default: conf.DefaultValues.Shell,
			NoEmpty: false,
			UseDefault: true,
		}

	vars.RecordFields[7] =
		vars.RecordField{ FieldName: "userPassword",
			Prompt: "Enter password: ",
			Default: utils.GenerateRandom(true, 25),
			NoEmpty: false,
			UseDefault: true,
		}

	vars.RecordFields[8] =
		vars.RecordField{ FieldName: "shadowMax",
			Prompt: "Enter new max password age",
			Default: strconv.Itoa(conf.DefaultValues.ShadowMax),
			NoEmpty: false,
			UseDefault: true,
		}

	vars.RecordFields[9] =
		vars.RecordField{ FieldName: "sshPublicKey",
			Prompt: "Enter SSH the Public Key",
			Default: "none",
			NoEmpty: false,
			UseDefault: false,
		}
}
