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
