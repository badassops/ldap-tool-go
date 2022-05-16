// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"badassops.ldap/consts"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
)


func CheckGroupID (conn *ldap.Connection, groupID int) (bool, string) {
	for _, groupMap := range conn.Config.GroupValues.GroupsMap {
		if groupMap.Gid == groupID {
			return true, groupMap.Name
		}
	}
	return false, "no-found"
}

func CheckProtectedGroup(conn *ldap.Connection, groupName string) bool {
	var groupsName []string
	groupsName = append(groupsName, conn.Config.GroupValues.Groups...)
	groupsName = append(groupsName, conn.Config.GroupValues.SpecialGroups...)

	for _, protectGroup := range groupsName {
		if groupName == protectGroup {
			return true
		}
	}
	return false
}

func CheckGroup(conn *ldap.Connection) (bool, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\tEnter the group name: ")
	enterData, _ := reader.ReadString('\n')
	enterData = strings.TrimSuffix(enterData, "\n")

	if enterData == "" {
		utils.PrintColor(consts.Red, fmt.Sprintf("\n\tNo group was given aborting...\n"))
		return false, "none-given"
	}

	if cnt, _ := conn.GetGroup(enterData) ; cnt == 0 {
		return false, enterData
	}
	return true, enterData
}

// Need to be global since the function is called recursive
var (
	groupType string
	enterData string
)

func Group(conn *ldap.Connection, firstTime bool) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\tEnters the group name to be use: ")
	enterData, _ = reader.ReadString('\n')
	enterData = strings.TrimSuffix(enterData, "\n")

	if enterData == "" {
		utils.PrintColor(consts.Red, fmt.Sprintf("\n\tNo group was given aborting...\n"))
		if firstTime {
			return "not-found"
		} else {
			// need to break the recursive
			utils.ReleaseIT(conn.LockFile, conn.LockPid)
			utils.PrintLine(utils.Purple)
			os.Exit(1)
		}
	}

	if firstTime {
		fmt.Printf("\tUse wildcard (default to N)? [y/n]: ")
		wildCard, _ := reader.ReadString('\n')
		wildCard = strings.TrimSuffix(wildCard, "\n")
		if utils.GetYN(wildCard, false) == true {
			enterData = "*" + enterData + "*"
			// conn.SearchGroup(enterData, groupType, false)
			conn.SearchGroup(enterData, false)
			fmt.Printf("\n\tSelect the group name from the above list:\n")
			Group(conn, false)
		}
	} else {
		// from recursive
		return "recursive"
	}

	utils.PrintLine(utils.Purple)
	if cnt := conn.SearchGroup(enterData, true) ; cnt == 0 {
		utils.PrintColor(consts.Red, fmt.Sprintf("\n\tGroup %s was not found, aborting...\n", enterData))
	}
	return enterData
}
