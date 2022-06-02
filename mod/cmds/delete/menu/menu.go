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

	deleteGroup "badassops.ldap/cmds/delete/group"
	deleteSudo "badassops.ldap/cmds/delete/sudo"
	deleteUser "badassops.ldap/cmds/delete/user"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/print"
)

func DeleteMenu(c *l.Connection) {
	p := print.New()
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Delete", 20, true))
	fmt.Printf("\tDelete (%s)ser, (%s)roup, (%s)udo role or (%s)uit?\n\t(default to User)? choice: ",
		p.MessageGreen("U"),
		p.MessageGreen("G"),
		p.MessageBlue("S"),
		p.MessageRed("Q"),
	)

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSuffix(choice, "\n")
	switch strings.ToLower(choice) {
	case "user", "u":
		deleteUser.Delete(c)
	case "group", "g":
		deleteGroup.Delete(c)
	case "sudo", "s":
		deleteSudo.Delete(c)
	case "quit", "q":
		p.PrintRed("\n\tOperation cancelled\n")
		fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 40))
		break
	default:
		deleteUser.Delete(c)
	}
}
