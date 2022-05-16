// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package delete

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"badassops.ldap/consts"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
	"badassops.ldap/cmds/common/group"
)

func Delete(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "delete group", true)
	reader := bufio.NewReader(os.Stdin)
	utils.PrintColor(consts.Yellow,
			fmt.Sprintf("\tEnter group name to be use: "))
	valueEntered, _ := reader.ReadString('\n')
	valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
	if valueEntered == "" {
		utils.PrintColor(consts.Red, fmt.Sprintf("\n\tNo users was given aborting...\n"))
    } else {
		if cnt, _ := conn.GetGroup(valueEntered) ; cnt == 0 {
			utils.PrintColor(consts.Red, fmt.Sprintf("\n\tGiven group %s does not found, aborting...\n", valueEntered))
			return
		}
		if common.CheckProtectedGroup(conn, valueEntered) {
			utils.PrintColor(consts.Red,
				fmt.Sprintf("\n\tGiven group %s is protected and can not be deleted, aborting...\n\n", valueEntered))
			return
		}
		utils.PrintColor(consts.Red,
			fmt.Sprintf("\n\tGiven group %s will be delete, this can not be undo!\n", valueEntered))
		utils.PrintColor(consts.Yellow,
			fmt.Sprintf("\tContinue (default to N)? [y/n]: "))
		continueDelete, _ := reader.ReadString('\n')
		continueDelete = strings.ToLower(strings.TrimSuffix(continueDelete, "\n"))
		if utils.GetYN(continueDelete, false) == true {
			if !conn.DeleteGroup(valueEntered) {
				utils.PrintColor(consts.Red,
					fmt.Sprintf("\n\tFailed to delete the group %s\n", valueEntered))
			} else {
				utils.PrintColor(consts.Green,
					fmt.Sprintf("\n\tGiven group %s has been deleted\n", valueEntered))
			}
		} else {
			utils.PrintColor(consts.Blue,
				fmt.Sprintf("\n\tDeletion of the group %s cancelled\n", valueEntered))
		}
	}
	utils.PrintLine(utils.Purple)
	return
}
