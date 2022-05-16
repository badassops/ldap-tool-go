// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package modify

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	//"regexp"

	"badassops.ldap/consts"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
	"badassops.ldap/cmds/common/group"
)

func cleanUpData (data string) string {
	// remove the userDN part
	newData := strings.Split(data, ",")[0]
	// remove the uid= part
	return strings.TrimPrefix(newData, "uid=")
}

func Modify(conn *ldap.Connection) {
	var delList []string
	var addList []string

	utils.PrintHeader(consts.Purple, "Modify Group", true)
	selectedGroup := common.Group(conn, true)
	reader := bufio.NewReader(os.Stdin)
	utils.PrintColor(consts.Yellow,
		fmt.Sprintf("\n\tEnter the user(s) to be deleted, select from the list above, (default to skip)\n"))
	for true {
		fmt.Printf("\tUser : ")
		enterData, _ := reader.ReadString('\n')
		enterData = strings.TrimSuffix(enterData, "\n")
		if enterData != "" {
			cleanUpData(enterData)
			delList = append(delList, cleanUpData(enterData))
		} else {
			break
		}
	}

	utils.PrintColor(consts.Yellow,
		fmt.Sprintf("\n\tEnter the user(s) to be added, (default to skip)\n"))
	for true {
		fmt.Printf("\tUser : ")
		enterData, _ := reader.ReadString('\n')
		enterData = strings.TrimSuffix(enterData, "\n")
		if enterData != "" {
			addList = append(addList, cleanUpData(enterData))
		} else {
			break
		}
	}
	if !conn.ModifyGroup(selectedGroup, addList, delList) {
		utils.PrintColor(consts.Red,
			fmt.Sprintf("\n\tFailed to modify the group %s\n", selectedGroup))
	} else {
		utils.PrintColor(consts.Green,
			fmt.Sprintf("\n\tGiven group %s has been modified\n", selectedGroup))
	}
	utils.PrintLine(utils.Purple)
	return
}
