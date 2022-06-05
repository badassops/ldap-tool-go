//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package limit

import (
	"fmt"
	"regexp"
	"strconv"

	l "badassops.ldap/ldap"
	u "badassops.ldap/utils"
)

var (
	displayFields = []string{"uid", "givenName", "sn", "cn", "displayName",
		"gecos", "uidNumber", "gidNumber", "departmentNumber",
		"mail", "homeDirectory", "loginShell", "userPassword",
		"shadowWarning", "shadowMax", "sshPublicKey"}
)

func printUserRecord(c *l.Connection, userName string) {
	// the values are in days so we need to multiple by 86400
	value, _ := strconv.ParseInt(c.User.Field["shadowLastChange"], 10, 64)
	_, passChanged := u.GetReadableEpoch(value * 86400)

	value, _ = strconv.ParseInt(c.User.Field["shadowExpire"], 10, 64)
	_, passExpired := u.GetReadableEpoch(value * 86400)

	for _, field := range displayFields {
		u.PrintCyan(fmt.Sprintf("\t%s: %s\n", field, c.User.Field[field]))
	}

	u.PrintLine(u.Purple)
	c.User.Groups = c.GetUserGroups("groupOfNames")
	u.PrintPurple(fmt.Sprintf("\tUser %s groups:\n", userName))
	for _, group := range c.User.Groups {
		u.PrintCyan(fmt.Sprintf("\tdn: %s\n", group))
	}

	u.PrintLine(u.Purple)
	u.PrintPurple(fmt.Sprintf("\tUser %s password information\n", userName))
	u.PrintCyan(fmt.Sprintf("\tPassword last changed on %s\n", passChanged))
	u.PrintRed(fmt.Sprintf("\tPassword will expired on %s\n", passExpired))
}

func UserRecord(c *l.Connection) {
	reg, _ := regexp.Compile("^uid=|,ou=users,.*")
	userID := reg.ReplaceAllString(c.Config.ServerValues.Admin, "")
	u.PrintHeader(u.Purple, fmt.Sprintf("Search User %s", userID), true)
	if c.GetUser(userID, false) == 0 {
		u.PrintColor(u.Red, fmt.Sprintf("\n\tUser %s was not found, aborting...\n", userID))
		return
	}
	printUserRecord(c, userID)
	u.PrintLine(u.Purple)
}
