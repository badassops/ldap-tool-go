// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package user

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"

	"badassops.ldap/constants"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
	"badassops.ldap/cmds/search/common"
)

func printUserRecord(conn *ldap.Connection, userID string, wildcard bool) int {
	if wildcard {
			userID = "*" + userID + "*"
	}
	records, cnt := conn.CheckUser(userID)
	if cnt == 0 {
		return 0
	}
	if wildcard {
		for _, entry := range records.Entries {
			utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s \n", entry.DN))
		}
		return cnt
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
	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tsshPublicKey: %s\n", conn.User.SSHPublicKey.Data))
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
	return cnt
}

func Search(conn *ldap.Connection, mode string) {
	utils.PrintHeader(constants.Purple, "Search " +  mode)
	reader := bufio.NewReader(os.Stdin)
	givenValue := ""
	wildCard := false
	cnt := 0
	if strings.HasSuffix(mode, "s") == false {
		givenValue, wildCard = common.EnterValue(mode)
		if givenValue == "" {
			return
		}
	}
	switch mode {
		case "user":
			if !wildCard {
				cnt = printUserRecord(conn, givenValue, false)
			} else {
				utils.PrintLine(utils.Purple)
				cnt = printUserRecord(conn, givenValue, true)
				fmt.Printf("\n\tSelect the userid from the above list: ")
				givenValue, _ = reader.ReadString('\n')
				givenValue = strings.TrimSuffix(givenValue, "\n")
				if givenValue == "" {
					utils.PrintColor(utils.Red, fmt.Sprintf("\tNo user was given aborting...\n"))
					return
				}
				cnt = printUserRecord(conn, givenValue, false)
			}

		case "users":
			cnt = conn.SearchUsers()

	}
	if cnt == 0 {
		utils.PrintColor(constants.Red, fmt.Sprintf("\tNo user match %s\n", givenValue))
		return
	}
	utils.PrintLine(utils.Purple)
}
