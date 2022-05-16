// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package create

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"badassops.ldap/consts"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
)

func Delete(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "Delete User", true)
	reader := bufio.NewReader(os.Stdin)
	utils.PrintColor(consts.Yellow,
			fmt.Sprintf("\tEnter userid (login name) to be use: "))
	valueEntered, _ := reader.ReadString('\n')
	valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
	if valueEntered == "" {
		utils.PrintColor(consts.Red, fmt.Sprintf("\n\tNo users was given aborting...\n"))
    } else {
		cnt := conn.CheckUser(valueEntered)
		if cnt == 0 {
			utils.PrintColor(consts.Red,
				fmt.Sprintf("\n\tGiven user %s doen not exist, aborting...\n\n", valueEntered))
			return
		} else {
			conn.User.Field["dn"] = fmt.Sprintf("uid=%s,%s", valueEntered, conn.Config.ServerValues.UserDN)
			utils.PrintColor(consts.Red,
				fmt.Sprintf("\n\tGiven user %s will be delete, this can not be undo!\n", valueEntered))
			utils.PrintColor(consts.Yellow,
				fmt.Sprintf("\tContinue (default to N)? [y/n]: "))
			continueDelete, _ := reader.ReadString('\n')
			continueDelete = strings.ToLower(strings.TrimSuffix(continueDelete, "\n"))
			if utils.GetYN(continueDelete, false) == true {
				if !conn.DeleteRecord() {
					utils.PrintColor(consts.Red,
						fmt.Sprintf("\n\tFailed to delete the user %s\n", conn.User.Field["uid"]))
				} else {
					utils.PrintColor(consts.Green,
						fmt.Sprintf("\n\tGiven user %s has been deleted\n", valueEntered))
				}
			} else {
				utils.PrintColor(consts.Blue,
					fmt.Sprintf("\n\tDeletion of the user %s cancelled\n", valueEntered))
			}
		}
	}
	utils.PrintLine(utils.Purple)
	return
}
