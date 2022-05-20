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

  "badassops.ldap/vars"
  "badassops.ldap/configurator"
  u "badassops.ldap/utils"
)

func Init(c *configurator.Config) {
  // ldap fields that will be used
  vars.Fields  = []string{"uid", "givenName", "sn", "cn", "displayName",
    "gecos", "uidNumber", "gidNumber", "departmentNumber",
    "mail", "homeDirectory", "loginShell", "userPassword",
    "shadowLastChange", "shadowExpire", "shadowWarning", "shadowMax",
    "sshPublicKey", "groups"}

  vars.Logs.LogsDir       = vars.LogsDir
  vars.Logs.LogFile       = vars.LogFile
  vars.Logs.LogMaxSize    = vars.LogMaxSize
  vars.Logs.LogMaxBackups = vars.LogMaxBackups
  vars.Logs.LogMaxAge     = vars.LogMaxAge
  vars.User.Field         = make(map[string]string)
  vars.Group              = make(map[string]string)
  vars.Template           = make(map[string]vars.Record)
  vars.GroupTemplate      = make(map[string]vars.Record)

  // set to expire by default as today + ShadowMax
  currExpired := strconv.FormatInt(u.GetEpoch("days") + int64(c.DefaultValues.ShadowMax), 10)

  // the fields are always needed
  vars.Template["uid"] =
    vars.Record{
      Prompt: "Enter userid (login name) to be use",
      Value: "",
      NoEmpty: true,
      UseValue: false,
    }

  vars.Template["givenName"] =
    vars.Record{
      Prompt: "Enter First name",
      Value: "",
      NoEmpty: true,
      UseValue: false,
    }

  vars.Template["sn"] =
    vars.Record{
      Prompt: "Enter Last name",
      Value: "",
      NoEmpty: true,
      UseValue: false,
    }

  vars.Template["mail"] =
    vars.Record{
      Prompt: "Enter email",
      Value: "",
      NoEmpty: false,
      UseValue: true,
    }

  vars.Template["uidNumber"] =
    vars.Record{
      Prompt: "Enter user's UID",
      Value: "",
      NoEmpty: false,
      UseValue: true,
    }

  vars.Template["departmentNumber"] =
    vars.Record{
      Prompt: "Enter department",
      Value: c.DefaultValues.GroupName,
      NoEmpty: false,
      UseValue: true,
    }

  vars.Template["loginShell"] =
    vars.Record{
      Prompt: "Enter shell",
      Value: c.DefaultValues.Shell,
      NoEmpty: false,
      UseValue: true,
    }

  vars.Template["userPassword"] =
    vars.Record{
      Prompt: "Enter password",
      Value: "",
      NoEmpty: false,
      UseValue: true,
    }

  vars.Template["shadowMax"] =
    vars.Record{
      Prompt: "Enter the max password age",
      Value: strconv.Itoa(c.DefaultValues.ShadowAge),
      NoEmpty: false,
      UseValue: true,
    }

  vars.Template["shadowWarning"] =
    vars.Record{
      Prompt: "Enter the days notification before the password expires",
      Value: strconv.Itoa(c.DefaultValues.ShadowWarning),
      NoEmpty: false,
      UseValue: true,
    }

  vars.Template["sshPublicKey"] =
    vars.Record{
      Prompt: "Enter SSH the Public Key",
      Value: "is missing",
      NoEmpty: false,
      UseValue: false,
    }

  // these are use during modification
  vars.Template["shadowExpire"] =
    vars.Record{
      Prompt: fmt.Sprintf("Reset password expired to (%d days from now) Y/N", c.DefaultValues.ShadowMax),
      Value: currExpired,
      NoEmpty: false,
      UseValue: false,
  }

  vars.Template["shadowLastChange"] =
    vars.Record{
      Prompt: "Date the password was last changed",
      Value: strconv.FormatInt(u.GetEpoch("days"), 10),
      NoEmpty: false,
      UseValue: false,
  }

  // these are for the ldap group
  vars.GroupTemplate["groupName"] =
    vars.Record{
      Prompt: "Enter the group name",
      Value: "",
      NoEmpty: true,
      UseValue: false,
  }

  vars.GroupTemplate["groupType"] =
    vars.Record{
      Prompt: "Group type (p)osix or (g)roupOfNames (default to posix)",
      Value: "posix",
      NoEmpty: false,
      UseValue: true,
  }

  // onlty use for posix group
  vars.GroupTemplate["gidNumber"] =
    vars.Record{
      Prompt: "Group ID/number of the posix group",
      Value: "",
      NoEmpty: false,
      UseValue: false,
  }

  // these are automatically filled
  vars.GroupTemplate["objectClass"] =
    vars.Record{
      Prompt: "Auto filled based on group type, posix or groupOfNames (default to posix)",
      Value: "",
      NoEmpty: true,
      UseValue: true,
  }

  vars.GroupTemplate["cn"] =
    vars.Record{
      Prompt: "Auto filled based on the groupDN value",
      Value: "",
      NoEmpty: true,
      UseValue: false,
  }

  // the defaul is always used
  vars.GroupTemplate["member"] =
    vars.Record{
      Prompt: "Auto filled based on the groupDN value",
      Value: fmt.Sprintf("uid=initial-member,%s", c.ServerValues.UserDN),
      NoEmpty: true,
      UseValue: false,
  }

}
