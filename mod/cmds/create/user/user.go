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
	//"strconv"

	"badassops.ldap/constants"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
//	"badassops.ldap/cmds/search/common"
)

//func createUserRecord() {
//userid
//firstName
//lastName
//}
func createRecord(conn *ldap.Connection) bool {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\tThe userid / login name is case sensitive, it will be made all lowercase)\n")
	fmt.Printf("\tEnter userid (login name) to be use: ")
	enterData, _ := reader.ReadString('\n')
	enterData = strings.ToLower(strings.TrimSuffix(enterData, "\n"))
	if enterData == "" {
		return false
	}

	_, cnt := conn.CheckUser(enterData)
	if cnt != 0 {
		return false
	}

	fmt.Printf("\tEnter First name: ")
	enterData, _ = reader.ReadString('\n')
	enterData = strings.Title(strings.TrimSuffix(enterData, "\n"))
	if enterData == "" {
		return false
	}

	return true
}

func Create(conn *ldap.Connection)  {
	utils.PrintHeader(constants.Purple, "create user")
	createRecord(conn)


	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tdn: %s\n", conn.User.DN))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tuid: %s\n", conn.User.UserName))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tgivenName: %s\n", FirstName))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tsn: %s\n", LastName))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tuidNumber: %d\n", conn.User.UID))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tgidNumber: %s\n", conn.User.GID.Data))

	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tcn: %s %s\n", FirstName, LastName))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tdisplayName: %s %s\n", FirstName, LastName))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tgecos: %s %s\n", FirstName, LastName))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tloginShell: %v\n", conn.User.Shell.Data))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\thomeDirectory: %v\n", conn.User.HomeDir))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tdepartmentNumber: %s\n", conn.User.GroupName))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tmail: %v\n", conn.User.Email.Data))

	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tuserPassword: %s\n", conn.User.Password.Data))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tshadowLastChange: %s\n", conn.User.ShadowLastChange.Data))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tshadowExpire: %s\n", conn.User.ShadowExpire.Data))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tshadowMax: %s\n", conn.User.ShadowMax.Data))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tshadowWarning: %s\n", conn.User.ShadowWarning.Data))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tsshPublicKey: %s\n", conn.User.SSHPublicKey.Data))
	// utils.PrintLine(utils.Purple)
	// utils.PrintColor(utils.Purple, fmt.Sprintf("\tUser %s admins groups:\n", conn.User.UserName))
	// for _, adminGroup := range conn.User.AdminGroups.Data {
	// 	utils.PrintColor(utils.Cyan, fmt.Sprintf("\tdn: %s\n", adminGroup))
	// }
	// for _, adminGroup := range conn.User.VPNGroups.Data {
	// 	utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s\n", adminGroup))
	// }
	// utils.PrintLine(utils.Purple)

	// utils.PrintColor(utils.Purple, fmt.Sprintf("\tUser %s password information\n", conn.User.UserName))
	// utils.PrintColor(utils.Cyan, fmt.Sprintf("\tPassword last changed on %s\n", passChanged))
	// utils.PrintColor(utils.Red, fmt.Sprintf("\tPassword will expired on %s\n", passExpired))
}
