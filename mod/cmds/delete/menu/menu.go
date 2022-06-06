//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
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

	"badassops.ldap/ldap"
	"badassops.ldap/vars"
)

func DeleteMenu(conn *ldap.Connection, funcs *vars.Funcs) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Delete", 20, true))
	fmt.Printf("\tDelete (%s)ser, (%s)roup, (%s)udo role or (%s)uit?\n\t(default to User)? choice: ",
		funcs.P.MessageGreen("U"),
		funcs.P.MessageGreen("G"),
		funcs.P.MessageBlue("S"),
		funcs.P.MessageRed("Q"),
	)

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSuffix(choice, "\n")
	switch strings.ToLower(choice) {
	case "user", "u":
		deleteUser.Delete(conn, funcs)
	case "group", "g":
		deleteGroup.Delete(conn, funcs)
	case "sudo", "s":
		deleteSudo.Delete(conn, funcs)
	case "quit", "q":
		funcs.P.PrintRed("\n\tOperation cancelled\n")
		fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 40))
		break
	default:
		deleteUser.Delete(conn, funcs)
	}
}
