// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package cmds

import (
	"fmt"
	//"badassops.ldap/utils"
	//"badassops.ldap/configurator"
	"badassops.ldap/ldap"
)

func Search(conn *ldap.Connection, mode string) {
	fmt.Printf("Hello World: %s\n", mode)
	fmt.Printf("%v\n", conn)
}
