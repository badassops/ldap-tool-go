// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package initializer

import (
	"badassops.ldap/constants"
)

func Init() {
	constants.UserLdapData.LoginName			= ""
	constants.UserLdapData.FirstName			= ""
	constants.UserLdapData.LastName				= ""
	constants.UserLdapData.UidNumber			= ""
	constants.UserLdapData.GidNumber			= ""
	constants.UserLdapData.Email				= ""
	constants.UserLdapData.Department			= ""
	constants.UserLdapData.Shell				= ""
	constants.UserLdapData.Password				= ""
	constants.UserLdapData.ShadowMax			= constants.ShadowMax
	constants.UserLdapData.ShadowExpired		= ""
	constants.UserLdapData.ShadowWarning		= constants.ShadowWarning
	constants.UserLdapData.ShadowLastChange		= ""
	constants.UserLdapData.SSHPublicKey			= ""
	constants.UserLdapData.AdminGroups			= []string{}
	constants.UserLdapData.VPNGroups			= []string{}
	constants.UserLdapData.HomeDirectory		= ""

	constants.ServerLog.LogsDir					= constants.LogsDir
	constants.ServerLog.LogFile					= constants.LogFile
	constants.ServerLog.LogMaxSize				= constants.LogMaxSize
	constants.ServerLog.LogMaxBackups			= constants.LogMaxBackups
	constants.ServerLog.LogMaxAge				= constants.LogMaxAge
}
