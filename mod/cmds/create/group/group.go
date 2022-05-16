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
	"strconv"

	"badassops.ldap/consts"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
	"badassops.ldap/cmds/common/group"
)

func createGroup(conn *ldap.Connection, groupName string, groupType string) bool {
	var groupID int
	switch groupType {
		case "posix", "p", "":
			fmt.Printf("\tEnter the group id (gid): ")
			reader := bufio.NewReader(os.Stdin)
			enterData, _ := reader.ReadString('\n')
			enterData = strings.TrimSuffix(enterData, "\n")
			groupID, _  = strconv.Atoi(enterData)
			found, groupName := common.CheckGroupID(conn, groupID)
			if found {
				utils.PrintColor(consts.Red,
				fmt.Sprintf("\n\tGiven group id %d already use by the group %s, aborting...\n", groupID, groupName))
				return false
			}
		default: groupID = 666
	}
	if !conn.AddGroup(groupName, groupType, groupID) {
		return false
	}
	return true
}

func Create(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "Create Group", true)
	found, groupName := common.CheckGroup(conn)
	if found {
		utils.PrintColor(consts.Red,
			fmt.Sprintf("\n\tGiven group %s already exist, aborting...\n", groupName))
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\tEnter group name type [(p)osix|(g)roupOfNames] (default to posix): ")
		enterData, _ := reader.ReadString('\n')
		enterData = strings.TrimSuffix(enterData, "\n")
		createGroup(conn, groupName, enterData)
	}
	utils.PrintLine(utils.Purple)
	return
}
