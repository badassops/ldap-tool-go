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

	// "badassops.ldap/consts"
	"badassops.ldap/vars"
	"badassops.ldap/configurator"
	"badassops.ldap/utils"
)

func Init(conf *configurator.Config) {
	// ldap fields that will be used
	vars.Fields	= []string{"uid", "givenName", "sn", "cn", "displayName",
		"gecos", "uidNumber", "gidNumber", "departmentNumber",
		"mail", "homeDirectory", "loginShell", "userPassword",
		"shadowLastChange", "shadowExpire", "shadowWarning", "shadowMax",
		"sshPublicKey", "groups"}

	vars.Logs.LogsDir		= vars.LogsDir
	vars.Logs.LogFile		= vars.LogFile
	vars.Logs.LogMaxSize	= vars.LogMaxSize
	vars.Logs.LogMaxBackups	= vars.LogMaxBackups
	vars.Logs.LogMaxAge		= vars.LogMaxAge
	vars.User.Field			= make(map[string]string)
	vars.Template			= make(map[string]vars.Record)

	// set to expire by default as today + ShadowMax
	currExpired := strconv.FormatInt(utils.GetEpoch("days") + int64(conf.DefaultValues.ShadowMax), 10)

	// the fields are always needed
	vars.Template["uid"] =
		vars.Record{
			Prompt: "Enter userid (login name) to be use",
			Value: "",
			NoEmpty: true,
			UseValue: false,
			Changed: false,
		}

	vars.Template["givenName"] =
		vars.Record{
			Prompt: "Enter First name",
			Value: "",
			NoEmpty: true,
			UseValue: false,
			Changed: false,
		}

	vars.Template["sn"] =
		vars.Record{
			Prompt: "Enter Last name",
			Value: "",
			NoEmpty: true,
			UseValue: false,
			Changed: false,
		}

	vars.Template["mail"] =
		vars.Record{
			Prompt: "Enter email",
			Value: "",
			NoEmpty: false,
			UseValue: true,
			Changed: false,
		}

	vars.Template["uidNumber"] =
		vars.Record{
			Prompt: "Enter user's UID",
			Value: "",
			NoEmpty: false,
			UseValue: true,
			Changed: false,
		}

	vars.Template["departmentNumber"] =
		vars.Record{
			Prompt: "Enter department",
			Value: conf.DefaultValues.GroupName,
			NoEmpty: false,
			UseValue: true,
			Changed: false,
		}

	vars.Template["loginShell"] =
		vars.Record{
			Prompt: "Enter shell",
			Value: conf.DefaultValues.Shell,
			NoEmpty: false,
			UseValue: true,
			Changed: false,
		}

	vars.Template["userPassword"] =
		vars.Record{
			Prompt: "Enter password",
			Value: "",
			NoEmpty: false,
			UseValue: true,
			Changed: false,
		}

	vars.Template["shadowMax"] =
		vars.Record{
			Prompt: "Enter the max password age",
			Value: strconv.Itoa(conf.DefaultValues.ShadowAge),
			NoEmpty: false,
			UseValue: true,
			Changed: false,
		}

	vars.Template["shadowWarning"] =
		vars.Record{
			Prompt: "Enter the days notification before the password expires",
			Value: strconv.Itoa(conf.DefaultValues.ShadowWarning),
			NoEmpty: false,
			UseValue: true,
			Changed: false,
		}

	vars.Template["sshPublicKey"] =
		vars.Record{
			Prompt: "Enter SSH the Public Key",
			Value: "is missing",
			NoEmpty: false,
			UseValue: false,
			Changed: false,
		}

	// these are use during modification
	vars.Template["shadowExpire"] =
		vars.Record{
			Prompt: "Date the password will expired",
			Value: currExpired,
			NoEmpty: false,
			UseValue: false,
			Changed: false,
	}

	vars.Template["shadowLastChange"] =
		vars.Record{
			Prompt: "Date the password was last changed",
			Value: strconv.FormatInt(utils.GetEpoch("days"), 10),
			NoEmpty: false,
			UseValue: false,
			Changed: false,
	}
}
