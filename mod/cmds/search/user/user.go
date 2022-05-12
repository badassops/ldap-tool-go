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

	"badassops.ldap/consts"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
	"badassops.ldap/cmds/search/common"
)

func User(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "Search User")
	common.User(conn, true)
	utils.PrintLine(utils.Purple)
}

func Users(conn *ldap.Connection) {
	reader := bufio.NewReader(os.Stdin)
	utils.PrintHeader(consts.Purple, "Search Users")
	fmt.Printf("\tPrint full name and department (default to N)? [y/n]: ")
	userInfo, _ := reader.ReadString('\n')
	userInfo = strings.TrimSuffix(userInfo, "\n")
	utils.PrintLine(utils.Purple)
	if utils.GetYN(userInfo, false) == true {
		conn.SearchUsers(true)
	} else {
		conn.SearchUsers(false)
	}
	utils.PrintLine(utils.Purple)
}
