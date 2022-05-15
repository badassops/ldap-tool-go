// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package group

import (
	"fmt"

	"badassops.ldap/consts"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
	"badassops.ldap/cmds/search/common"
)


func Group(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "Search Group", false)
	common.Group(conn, true)
    utils.PrintLine(utils.Purple)
}

func Groups(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "Search Groups", false)
	utils.PrintColor(consts.Purple, fmt.Sprintf("\n\t__________ all group and the members __________\n"))
	conn.SearchGroups()
    utils.PrintLine(utils.Purple)
}
