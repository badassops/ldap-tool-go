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

  v "badassops.ldap/vars"
  c "badassops.ldap/configurator"
  u "badassops.ldap/utils"
)

func Init(c *c.Config) {
  // ldap fields that will be used
  v.Fields  = []string{"uid", "givenName", "sn", "cn", "displayName",
    "gecos", "uidNumber", "gidNumber", "departmentNumber",
    "mail", "homeDirectory", "loginShell", "userPassword",
    "shadowLastChange", "shadowExpire", "shadowWarning", "shadowMax",
    "sshPublicKey", "groups"}

  v.Sudoers = []string{"cn", "sudoCommand", "sudoHost", "sudoOption",
    "sudoOrder", "sudoRunAsUser" }

  v.SudoObjectClass = []string{"top", "sudoRole"}

  v.Logs.LogsDir       = v.LogsDir
  v.Logs.LogFile       = v.LogFile
  v.Logs.LogMaxSize    = v.LogMaxSize
  v.Logs.LogMaxBackups = v.LogMaxBackups
  v.Logs.LogMaxAge     = v.LogMaxAge
  v.User.Field         = make(map[string]string)
  v.Group              = make(map[string]string)
  v.Template           = make(map[string]v.Record)
  v.GroupTemplate      = make(map[string]v.Record)
  v.ModRecord.Field    = make(map[string]string)
  v.SudoTemplate       = make(map[string]v.Record)

  // set to expire by default as today + ShadowMax
  currExpired := strconv.FormatInt(u.GetEpoch("days") + int64(c.DefaultValues.ShadowMax), 10)

  // the fields are always needed
  v.Template["uid"] =
    v.Record{
      Prompt: "Enter userid (login name) to be use",
      Value: "",
      NoEmpty: true,
      UseValue: false,
    }

  v.Template["givenName"] =
    v.Record{
      Prompt: "Enter First name",
      Value: "",
      NoEmpty: true,
      UseValue: false,
    }

  v.Template["sn"] =
    v.Record{
      Prompt: "Enter Last name",
      Value: "",
      NoEmpty: true,
      UseValue: false,
    }

  v.Template["mail"] =
    v.Record{
      Prompt: "Enter email",
      Value: "",
      NoEmpty: false,
      UseValue: true,
    }

  v.Template["uidNumber"] =
    v.Record{
      Prompt: "Enter user's UID",
      Value: "",
      NoEmpty: false,
      UseValue: true,
    }

  v.Template["departmentNumber"] =
    v.Record{
      Prompt: "Enter department",
      Value: c.DefaultValues.GroupName,
      NoEmpty: false,
      UseValue: true,
    }

  v.Template["loginShell"] =
    v.Record{
      Prompt: "Enter shell",
      Value: c.DefaultValues.Shell,
      NoEmpty: false,
      UseValue: true,
    }

  v.Template["userPassword"] =
    v.Record{
      Prompt: "Enter password",
      Value: "",
      NoEmpty: false,
      UseValue: true,
    }

  v.Template["shadowMax"] =
    v.Record{
      Prompt: "Enter the max password age",
      Value: strconv.Itoa(c.DefaultValues.ShadowAge),
      NoEmpty: false,
      UseValue: true,
    }

  v.Template["shadowWarning"] =
    v.Record{
      Prompt: "Enter the days notification before the password expires",
      Value: strconv.Itoa(c.DefaultValues.ShadowWarning),
      NoEmpty: false,
      UseValue: true,
    }

  v.Template["sshPublicKey"] =
    v.Record{
      Prompt: "Enter SSH the Public Key",
      Value: "is missing",
      NoEmpty: false,
      UseValue: false,
    }

  // these are use during modification
  v.Template["shadowExpire"] =
    v.Record{
      Prompt: fmt.Sprintf("Reset password expired to (%d days from now) Y/N", c.DefaultValues.ShadowMax),
      Value: currExpired,
      NoEmpty: false,
      UseValue: false,
  }

  v.Template["shadowLastChange"] =
    v.Record{
      Prompt: "Date the password was last changed",
      Value: strconv.FormatInt(u.GetEpoch("days"), 10),
      NoEmpty: false,
      UseValue: false,
  }

  // these are for the ldap group
  v.GroupTemplate["groupName"] =
    v.Record{
      Prompt: "Enter the group name",
      Value: "",
      NoEmpty: true,
      UseValue: false,
  }

  v.GroupTemplate["groupType"] =
    v.Record{
      Prompt: "Group type (p)osix or (g)roupOfNames (default to posix)",
      Value: "posix",
      NoEmpty: false,
      UseValue: true,
  }

  // onlty use for posix group
  v.GroupTemplate["gidNumber"] =
    v.Record{
      Prompt: "Group ID/number of the posix group",
      Value: "",
      NoEmpty: false,
      UseValue: true,
  }

  // these are automatically filled
  v.GroupTemplate["objectClass"] =
    v.Record{
      Prompt: "Auto filled based on group type, posix or groupOfNames (default to posix)",
      Value: "",
      NoEmpty: true,
      UseValue: true,
  }

  v.GroupTemplate["cn"] =
    v.Record{
      Prompt: "Auto filled based on the groupDN value",
      Value: "",
      NoEmpty: true,
      UseValue: false,
  }

  // the defaul is always used
  v.GroupTemplate["member"] =
    v.Record{
      Prompt: "Auto filled based on the groupDN value",
      Value: fmt.Sprintf("uid=initial-member,%s", c.ServerValues.UserDN),
      NoEmpty: true,
      UseValue: false,
  }

  // the sudoer template
  v.SudoTemplate["cn"] =
    v.Record{
      Prompt: fmt.Sprintf("%sUse %s%%%s%s to indicate a group%s\n\tEnter the sudoers CN",
       u.Yellow, u.Cyan, u.Off, u.Yellow, u.Off),
      Value: "",
      NoEmpty: true,
      UseValue: false,
  }

  v.SudoTemplate["sudoCommand"] =
    v.Record{
      Prompt: fmt.Sprintf("%sfully qualified path or ALL%s\n\tEnter the command allow with this rule",
        u.Yellow, u.Off),
      Value: "",
      NoEmpty: true,
      UseValue: false,
  }

  v.SudoTemplate["sudoHost"] =
    v.Record{
      Prompt: fmt.Sprintf("%sdefault to ALL%s\n\tThe host the command is allowed",
        u.Yellow, u.Off),
      Value: "ALL",
      NoEmpty: false,
      UseValue: true,
  }

  v.SudoTemplate["sudoOption"] =
    v.Record{
      Prompt: fmt.Sprintf("%sexmple %s!authenticate%s%s for no password required%s\n\tSudo option with the command",
        u.Yellow, u.Cyan, u.Off, u.Yellow, u.Off),
      Value: "",
      NoEmpty: false,
      UseValue: false,
  }

  v.SudoTemplate["sudoOrder"] =
    v.Record{
      Prompt: fmt.Sprintf("%sdefault to 4, use 3 and not higher than 10%s\n\tThe order of the rule",
        u.Yellow, u.Off),
      Value: "4",
      NoEmpty: false,
      UseValue: true,
  }

  v.SudoTemplate["sudoRunAsUser"] =
    v.Record{
      Prompt: fmt.Sprintf("%sdefault to %sroot%s\n\tRun the command as the user", u.Yellow, u.Red, u.Off),
      Value: "root",
      NoEmpty: false,
      UseValue: true,
  }
}
