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
	"strconv"
	"strings"

	"badassops.ldap/consts"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
)

var (
	displayFields = []string{"uid", "givenName", "sn", "cn", "displayName",
        "gecos", "uidNumber", "gidNumber", "departmentNumber",
        "mail", "homeDirectory", "loginShell", "userPassword",
        "shadowWarning", "shadowMax", "sshPublicKey"}

)

func printUserRecord(conn *ldap.Connection, userName string) {
	// the values are in days so we need to multiple by 86400
	value, _ := strconv.ParseInt(conn.User.Field["shadowLastChange"], 10, 64)
	 _, passChanged := utils.GetReadableEpoch(value * 86400)

	value, _ = strconv.ParseInt(conn.User.Field["shadowExpire"], 10, 64)
	_, passExpired := utils.GetReadableEpoch(value * 86400)

	utils.PrintLine(utils.Purple)
	for _, field := range displayFields {
		utils.PrintColor(utils.Cyan, fmt.Sprintf("\t%s: %s\n", field, conn.User.Field[field]))
	}

	utils.PrintLine(utils.Purple)
	utils.PrintColor(utils.Purple, fmt.Sprintf("\tUser %s groups:\n", userName))
	for _, group := range conn.User.Groups {
		utils.PrintColor(utils.Cyan, fmt.Sprintf("\tdn: %s\n", group))
	}

	utils.PrintLine(utils.Purple)
	utils.PrintColor(utils.Purple, fmt.Sprintf("\tUser %s password information\n", userName))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tPassword last changed on %s\n", passChanged))
	utils.PrintColor(utils.Red, fmt.Sprintf("\tPassword will expired on %s\n", passExpired))
}

func User(conn *ldap.Connection, firstTime bool, showRecord bool) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\tEnter user login name to be use: ")
	enterData, _ := reader.ReadString('\n')
	enterData = strings.TrimSuffix(enterData, "\n")

	if enterData == "" {
		utils.PrintColor(consts.Red, fmt.Sprintf("\n\tNo users was given aborting...\n"))
		return false
	}

	if firstTime {
		fmt.Printf("\tUse wildcard (default to N)? [y/n]: ")
		wildCard, _ := reader.ReadString('\n')
		wildCard = strings.TrimSuffix(wildCard, "\n")
		if utils.GetYN(wildCard, false) == true {
			enterData = "*" + enterData + "*"
			conn.SearchUser(enterData)
			fmt.Printf("\n\tSelect the userid from the above list:\n")
			User(conn, false, showRecord)
			return true
		}
	}

	if conn.GetUser(enterData) == 0 {
		utils.PrintColor(consts.Red, fmt.Sprintf("\n\tUser %s was not found, aborting...\n", enterData))
		return false
	}
	if showRecord {
		printUserRecord(conn, enterData)
	}
	return true
}
