// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package initializer

import (
	"fmt"
	"strconv"

	c "badassops.ldap/configurator"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/epoch"
	"github.com/badassops/packages-go/print"
)

var (
	msg string
)

func Init(c *c.Config) {
	// ldap fields that will be used
	v.UserFields = []string{"uid", "givenName", "sn", "cn", "displayName",
		"gecos", "uidNumber", "gidNumber", "departmentNumber",
		"mail", "homeDirectory", "loginShell", "userPassword",
		"shadowLastChange", "shadowExpire", "shadowWarning", "shadowMax",
		"sshPublicKey", "groups"}

	v.DisplayUserFields = []string{"uid", "givenName", "sn", "cn", "displayName",
		"gecos", "uidNumber", "gidNumber", "departmentNumber",
		"mail", "homeDirectory", "loginShell", "userPassword",
		"shadowWarning", "shadowMax", "sshPublicKey"}

	v.GroupFields = []string{"cn", "groupName", "groupType", "gidNumber", "memberUid", "member"}

	v.SudoFields = []string{"cn", "sudoCommand", "sudoHost", "sudoOption",
		"sudoOrder", "sudoRunAsUser"}

	v.UserObjectClass = []string{"top", "person",
		"organizationalPerson", "inetOrgPerson",
		"posixAccount", "shadowAccount", "ldapPublicKey"}

	v.GroupObjectClass = []string{"posix", "groupOfNames"}

	v.SudoObjectClass = []string{"top", "sudoRole"}

	v.Logs.LogsDir = v.LogsDir
	v.Logs.LogFile = v.LogFile
	v.Logs.LogMaxSize = v.LogMaxSize
	v.Logs.LogMaxBackups = v.LogMaxBackups
	v.Logs.LogMaxAge = v.LogMaxAge

	v.WorkRecord.Fields = make(map[string]string)
	v.WorkRecord.Group = make(map[string]string)
	v.WorkRecord.SudoAddList = make(map[string][]string)
	v.WorkRecord.SudoDelList = make(map[string][]string)
	v.Template = make(map[string]v.Record)

	// set to expire by default as today + ShadowMax
	e := epoch.New()
	p := print.New()
	currExpired := strconv.FormatInt(e.Days()+int64(c.DefaultValues.ShadowMax), 10)

	// user
	v.Template["uid"] =
		v.Record{
			Prompt:   "Enter userid (login name) to be use",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	v.Template["givenName"] =
		v.Record{
			Prompt:   "Enter First name",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	v.Template["sn"] =
		v.Record{
			Prompt:   "Enter Last name",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	v.Template["mail"] =
		v.Record{
			Prompt:   "Enter email",
			Value:    "",
			NoEmpty:  false,
			UseValue: true,
		}

	v.Template["uidNumber"] =
		v.Record{
			Prompt:   "Enter user's UID",
			Value:    "",
			NoEmpty:  false,
			UseValue: true,
		}

	v.Template["departmentNumber"] =
		v.Record{
			Prompt:   "Enter department",
			Value:    c.DefaultValues.GroupName,
			NoEmpty:  false,
			UseValue: true,
		}

	v.Template["loginShell"] =
		v.Record{
			Prompt:   "Enter shell",
			Value:    c.DefaultValues.Shell,
			NoEmpty:  false,
			UseValue: true,
		}

	v.Template["userPassword"] =
		v.Record{
			Prompt:   "Enter password",
			Value:    "",
			NoEmpty:  false,
			UseValue: true,
		}

	v.Template["shadowMax"] =
		v.Record{
			Prompt:   "Enter the max password age",
			Value:    strconv.Itoa(c.DefaultValues.ShadowAge),
			NoEmpty:  false,
			UseValue: true,
		}

	v.Template["shadowWarning"] =
		v.Record{
			Prompt:   "Enter the days notification before the password expires",
			Value:    strconv.Itoa(c.DefaultValues.ShadowWarning),
			NoEmpty:  false,
			UseValue: true,
		}

	v.Template["sshPublicKey"] =
		v.Record{
			Prompt:   "Enter SSH the Public Key",
			Value:    "is missing",
			NoEmpty:  false,
			UseValue: false,
		}

	v.Template["shadowExpire"] =
		v.Record{
			Prompt:   fmt.Sprintf("Reset password expired, Y/N"),
			Value:    currExpired,
			NoEmpty:  false,
			UseValue: false,
		}

	v.Template["shadowLastChange"] =
		v.Record{
			Prompt:   "Date the password was last changed",
			Value:    strconv.FormatInt(e.Days(), 10),
			NoEmpty:  false,
			UseValue: false,
		}

	// share in group and sudo rule
	v.Template["cn"] =
		v.Record{
			Prompt:   "Auto filled based on the groupDN value",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	// group
	v.Template["groupName"] =
		v.Record{
			Prompt:   "Enter the group name",
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	v.Template["groupType"] =
		v.Record{
			Prompt:   "Group type (p)osix or (g)roupOfNames (default to posix)",
			Value:    "posix",
			NoEmpty:  false,
			UseValue: true,
		}

	// onlty use for posix group
	v.Template["gidNumber"] =
		v.Record{
			Prompt:   "Group ID/number of the posix group",
			Value:    "",
			NoEmpty:  false,
			UseValue: true,
		}

	// these are automatically filled
	v.Template["objectClass"] =
		v.Record{
			Prompt:   "Auto filled based on group type, posix or groupOfNames (default to posix)",
			Value:    "",
			NoEmpty:  true,
			UseValue: true,
		}

	// the defaul is always used
	v.Template["member"] =
		v.Record{
			Prompt:   "Auto filled based on the groupDN value",
			Value:    fmt.Sprintf("uid=initial-member,%s", c.ServerValues.UserDN),
			NoEmpty:  true,
			UseValue: false,
		}

	// sudo rules
	v.Template["sudoCommand"] =
		v.Record{
			Prompt: fmt.Sprintf("%sfully qualified path or ALL%s\n\tEnter the command allow with this rule",
				v.Yellow, v.Off),
			Value:    "",
			NoEmpty:  true,
			UseValue: false,
		}

	msg = p.MessageYellow("default to ALL")
	v.Template["sudoHost"] =
		v.Record{
			Prompt:   fmt.Sprintf("%s\n\tThe host the command is allowed", msg),
			Value:    "ALL",
			NoEmpty:  false,
			UseValue: true,
		}

	msg = p.MessageYellow("exmple %s!authenticate")
	msg = msg + p.MessageCyan(" or no password required")
	v.Template["sudoOption"] =
		v.Record{
			Prompt:   fmt.Sprintf("%s\n\tSudo option with the command", msg),
			Value:    "",
			NoEmpty:  false,
			UseValue: false,
		}

	msg = p.MessageYellow("default to 4, use 3 and not higher than 10")
	v.Template["sudoOrder"] =
		v.Record{
			Prompt:   fmt.Sprintf("%s\n\tThe order of the rule", msg),
			Value:    "4",
			NoEmpty:  false,
			UseValue: true,
		}

	msg = p.MessageYellow("default to ")
	msg = p.MessageRed("root")
	v.Template["sudoRunAsUser"] =
		v.Record{
			Prompt:   fmt.Sprintf("%s\n\tRun the command as the user", msg),
			Value:    "root",
			NoEmpty:  false,
			UseValue: true,
		}
}
