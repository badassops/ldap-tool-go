// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package cmds

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"

	"badassops.ldap/constants"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"

	// ldapv3 "gopkg.in/ldap.v2"
)

func printUserRecord(conn *ldap.Connection, userID string, wildcard bool)  {
	if wildcard {
			userID = "*" + userID + "*"
	}
	records, cnt := conn.CheckUser(true, userID)
	if cnt == 0 {
		utils.PrintColor(constants.Red, fmt.Sprintf("\tNo user match %s\n", userID))
		return
	}
	fmt.Printf("\n")
	if wildcard {
		for _, entry := range records.Entries {
			utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s \n", entry.DN))
		}
		return
	}
	FirstName := conn.User.FirstName.Data
	LastName := conn.User.LastName.Data

	// the values are in days so we need to multiple by 86400
	value, _ := strconv.ParseInt(conn.User.ShadowLastChange.Data, 10, 64)
	_, passChanged := utils.GetReadableEpoch(value * 86400)

	value, _ = strconv.ParseInt(conn.User.ShadowExpire.Data, 10, 64)
	_, passExpired := utils.GetReadableEpoch(value * 86400)

	utils.PrintLine(utils.Purple)
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tdn: %s\n", conn.User.DN))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tuid: %s\n", conn.User.UserName))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tgivenName: %s\n", FirstName))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tsn: %s\n", LastName))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tuidNumber: %d\n", conn.User.UID))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tgidNumber: %s\n", conn.User.GID.Data))

	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tcn: %s %s\n", FirstName, LastName))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tdisplayName: %s %s\n", FirstName, LastName))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tgecos: %s %s\n", FirstName, LastName))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tloginShell: %v\n", conn.User.Shell.Data))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\thomeDirectory: %v\n", conn.User.HomeDir))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tdepartmentNumber: %s\n", conn.User.GroupName))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tmail: %v\n", conn.User.Email.Data))

	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tuserPassword: %s\n", conn.User.Password.Data))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tshadowLastChange: %s\n", conn.User.ShadowLastChange.Data))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tshadowExpire: %s\n", conn.User.ShadowExpire.Data))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tshadowMax: %s\n", conn.User.ShadowMax.Data))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tshadowWarning: %s\n", conn.User.ShadowWarning.Data))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tsshPublicKey: %s\n", conn.User.ShadowWarning.Data))
	utils.PrintLine(utils.Purple)
	utils.PrintColor(utils.Purple, fmt.Sprintf("\tUser %s admins groups:\n", conn.User.UserName))
	for _, adminGroup := range conn.User.AdminGroups.Data {
		utils.PrintColor(utils.Cyan, fmt.Sprintf("\tdn: %s\n", adminGroup))
	}
	for _, adminGroup := range conn.User.VPNGroups.Data {
		utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s\n", adminGroup))
	}
	utils.PrintLine(utils.Purple)

	utils.PrintColor(utils.Purple, fmt.Sprintf("\tUser %s password information\n", conn.User.UserName))
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tPassword last changed on %s\n", passChanged))
	utils.PrintColor(utils.Red, fmt.Sprintf("\tPassword will expired on %s\n", passExpired))
	utils.PrintLine(utils.Purple)
}

func Search(conn *ldap.Connection, mode string) {
	utils.PrintHeader(constants.Purple, "Search " +  mode)
	reader := bufio.NewReader(os.Stdin)
	switch mode {
		case "user":
			fmt.Printf("\tEnter userid to be use: ")
			userID, _ := reader.ReadString('\n')
			userID = strings.TrimSuffix(userID, "\n")

			fmt.Printf("\tUse wildcard (default to N)? [y/n]: ")
			wildCard, _ := reader.ReadString('\n')
			wildCard = strings.TrimSuffix(wildCard, "\n")

			if utils.GetYN(wildCard, false) == false {
				printUserRecord(conn, userID, false)
				return
			} else {
				printUserRecord(conn, userID, true)
				fmt.Printf("\n\tSelect the userid from the above list: ")
				userID, _ := reader.ReadString('\n')
				userID = strings.TrimSuffix(userID, "\n")
				printUserRecord(conn, userID, false)
			}

		case "users":
			conn.SearchUsers()

		case "group":
		case "groups":
		case "admin":
		case "admins":
	}

}
