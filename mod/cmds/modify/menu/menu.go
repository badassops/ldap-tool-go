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

	modifyGroup "badassops.ldap/cmds/modify/group"
	modifySudo "badassops.ldap/cmds/modify/sudo"
	//modifyUser "badassops.ldap/cmds/modify/user"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/print"
)

func ModifyMenu(c *l.Connection) {
	p := print.New()
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Modify", 20, true))
	fmt.Printf("\tModify (%s)ser, (%s)roup, (%s)udo rule or (%s)uit?\n\t(default to User)? choice: ",
		p.MessageGreen("U"),
		p.MessageGreen("G"),
		p.MessageBlue("S"),
		p.MessageRed("Q"),
	)

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSuffix(choice, "\n")
	switch strings.ToLower(choice) {
	case "user", "u":
		//modifyUser.Modify(c)
	case "group", "g":
		modifyGroup.Modify(c)
	case "sudo", "s":
		modifySudo.Modify(c)
	case "quit", "q":
		p.PrintRed("\n\tOperation cancelled\n")
		fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 40))
		break
	default:
		//modifyUser.Modify(c)
	}
}
