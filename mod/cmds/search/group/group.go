// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package group

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


func Group(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "Search Group")
	common.Group(conn, true)
    utils.PrintLine(utils.Purple)
}

func Groups(conn *ldap.Connection) {
	var groupType string
	reader := bufio.NewReader(os.Stdin)
	utils.PrintHeader(consts.Purple, "Search Groups")
	fmt.Printf("\tGroup type [p]osix or [m]emberOf (default to posix) [p/n]: ")
	enterType, _ := reader.ReadString('\n')
	enterType = strings.TrimSuffix(enterType, "\n")
	switch enterType {
		case "p", "posix", "":	groupType = "posix"
		case "m", "memberof":	groupType = "memberof"
		default:				groupType = "posix"
	}
	utils.PrintColor(consts.Purple, fmt.Sprintf("\n\t__________ all group and the members __________\n"))
	conn.SearchGroup("*", groupType, true)
    utils.PrintLine(utils.Purple)
}
