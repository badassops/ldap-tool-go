// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	createGroup "badassops.ldap/cmds/create/group"
	//createSudo "badassops.ldap/cmds/create/sudo"
	//createUser "badassops.ldap/cmds/create/user"

    l "badassops.ldap/ldap"
    v "badassops.ldap/vars"
    "github.com/badassops/packages-go/print"
)

func CreateMenu(c *l.Connection) {
	reader := bufio.NewReader(os.Stdin)

	u.PrintHeader(u.Purple, "Create", true)
	fmt.Printf("\tCreate (%s)ser, (%s)roup, (%s)sudo role or (%s)uit?\n\t(default to User)? choice: ",
		p.MessageGreen, "U"),
		p.MessageGreen, "G"),
		p.MessageBlue, "S"),
		p.MessageRed, "Q"),
	)

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSuffix(choice, "\n")
	switch strings.ToLower(choice) {
	case "user", "u":
		//createUser.Create(c)
	case "group", "g":
		createGroup.Create(c)
	case "sudo", "s":
		//createSudo.Create(c)
	case "quit", "q":
		p.PrintRed("\n\t\tOperation cancelled\n")
		fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 40))
		break
	default:
		//createUser.Create(c)
	}
}
