// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package limit

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	l "badassops.ldap/ldap"
	u "badassops.ldap/utils"
	v "badassops.ldap/vars"
)

var (
	fields = []string{"userPassword", "sshPublicKey"}

	// input
	valueEntered string
)

func createModifyUserPasswordSSHKey(c *l.Connection) int {
	changeCount := 0
	u.PrintPurple(fmt.Sprintf("\tUsing user: %s\n", c.User.Field["uid"]))
	u.PrintYellow(fmt.Sprintf("\tPress enter to leave the value unchanged\n"))
	u.PrintLine(u.Purple)

	for _, fieldName := range fields {
		// these will be valid once the field was filled since they depends
		// on some of the fields value
		switch fieldName {
		case "userPassword":
			passWord := u.GenerateRandom(
				c.Config.DefaultValues.PassComplex,
				c.Config.DefaultValues.PassLenght)
			u.PrintCyan(fmt.Sprintf("\tCurrent value (encrypted!): %s\n", c.User.Field[fieldName]))
			u.PrintYellow(fmt.Sprintf("\t\tsuggested password: %s\n", passWord))

		case "sshPublicKey":
			u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", c.User.Field[fieldName]))

		}

		fmt.Printf("\t%s: ", v.Template[fieldName].Prompt)
		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		if len(valueEntered) != 0 {
			v.ModRecord.Field[fieldName] = valueEntered
			changeCount++
		}
	}
	return changeCount
}

func ModifyUserPasswordSSHKey(c *l.Connection) {
	reg, _ := regexp.Compile("^uid=|,ou=users,.*")
	userID := reg.ReplaceAllString(c.Config.ServerValues.Admin, "")
	u.PrintHeader(u.Purple, fmt.Sprintf("Modify User's Password / SSH Public Key %s", userID), true)
	if c.GetUser(userID, false) == 0 {
		u.PrintColor(u.Red, fmt.Sprintf("\n\tUser %s was not found, aborting...\n", userID))
		return
	}
	if createModifyUserPasswordSSHKey(c) == 0 {
		u.PrintBlue(fmt.Sprintf("\n\tNo field were changed, no modification was made for the user %s\n", userID))
	} else {
		c.User.Field["uid"] = userID
		if !c.ModifyUser() {
			u.PrintRed(fmt.Sprintf("\n\tFailed modify the user %s, check the log file\n", c.User.Field["uid"]))
		} else {
			u.PrintGreen(fmt.Sprintf("\n\tUser %s modified successfully\n", c.User.Field["uid"]))
		}
	}
	u.PrintLine(u.Purple)
}
