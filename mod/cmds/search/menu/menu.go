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

	searchGroup "badassops.ldap/cmds/search/group"
	searchSudo "badassops.ldap/cmds/search/sudo"
	searchUser "badassops.ldap/cmds/search/user"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/print"
)

func SearchMenu(c *l.Connection) {
	p := print.New()
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Search", 20, true))
	fmt.Printf("\tSearch (%s)ser, (%s)ll Users, (%s)roup, all Group(%s)\n",
		p.MessageGreen("U"),
		p.MessageGreen("A"),
		p.MessageGreen("G"),
		p.MessageGreen("S"),
	)
	fmt.Printf("\t\t(%s)sudo role, (%s)all sudos role or (%s)uit?\n\t(default to User)? choice: ",
		p.MessageBlue("X"),
		p.MessageBlue("Z"),
		p.MessageRed("Q"),
	)

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSuffix(choice, "\n")
	switch strings.ToLower(choice) {
	case "user", "u":
		searchUser.User(c)
	case "users", "a":
		searchUser.Users(c)
	case "group", "g":
		searchGroup.Group(c)
	case "groups", "s":
		searchGroup.Groups(c)
	case "sudo", "x":
		searchSudo.Sudo(c)
	case "sudos", "z":
		searchSudo.Sudos(c)
	case "quit", "q":
		p.PrintRed("\n\t\tOperation cancelled\n")
		fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 40))
		break
	default:
		searchUser.User(c)
	}
}
