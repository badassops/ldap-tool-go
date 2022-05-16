// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package modify

import (
	"fmt"

	"badassops.ldap/consts"
	"badassops.ldap/utils"
	"badassops.ldap/ldap"
)

func Modify(conn *ldap.Connection) {
	utils.PrintHeader(consts.Purple, "Modify User", true)
	utils.PrintColor(consts.Green, fmt.Sprintf("\n\t coming soon\n"))
	utils.PrintLine(utils.Purple)
	return
}
