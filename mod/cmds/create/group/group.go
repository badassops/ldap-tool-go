// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package create

import (
	"bufio"
	"fmt"
	"os"
	//"regexp"
	"strconv"
	"strings"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"

	"github.com/badassops/packages-go/is"
	"github.com/badassops/packages-go/print"
)

var (
	// group fields
	// required: groupName and groupType
	// required if posix: gidNumber
	// autofilled: objectClass, cn
	// autofilled if not posix: member

	valueEntered string
	nextGID      int
	fields       = []string{"groupName", "groupType"}
	validTypes   = []string{"posix", "groupOfNames"}

	p = print.New()
	i = is.New()
)

func createGroup(c *l.Connection) bool {
	allGroupDN := c.GetAlGroups()

	for _, fieldName := range fields {
		if v.Template[fieldName].Value != "" {
			p.PrintYellow(fmt.Sprintf("\t ** Default to: %s **\n", v.Template[fieldName].Value))
		}

		fmt.Printf("\t%s: ", v.Template[fieldName].Prompt)

		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")

		fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
		switch fieldName {
		case "groupName":
			groupDN := fmt.Sprintf("cn=%s,%s", valueEntered, c.Config.ServerValues.GroupDN)
			if i.IsInList(allGroupDN, groupDN) {
				p.PrintRed(fmt.Sprintf("\n\tGiven group %s already exist, aborting...\n\n", valueEntered))
				return false
			}
			p.PrintPurple(fmt.Sprintf("\tUsing Group: %s\n\n", valueEntered))
			v.WorkRecord.Fields["cn"] = valueEntered
			v.WorkRecord.Fields["dn"] = groupDN

		case "groupType":
			switch valueEntered {
			case "", "p", "posix":
				v.WorkRecord.Fields["objectClass"] = "posixGroup"
				valueEntered = "posix"
			case "g", "groupOfNames":
				v.WorkRecord.Fields["objectClass"] = "groupOfNames"
				// hard coded, groupOfNames must have at least 1 member
				v.WorkRecord.Fields["member"] = fmt.Sprintf("uid=initial-user,%s", c.Config.ServerValues.GroupDN)
				valueEntered = "groupOfNames"
			}
			if !i.IsInList(validTypes, valueEntered) {
				p.PrintRed(fmt.Sprintf("\tWrong group type (%s) aborting...\n\n", valueEntered))
				return false
			}
		}

		if (len(valueEntered) == 0) && (v.Template[fieldName].NoEmpty == true) {
			p.PrintRed("\tNo value was entered aborting..%s.\n\n")
			return false
		}
	}

	if  v.WorkRecord.Fields["objectClass"] == "posixGroup" {
		v.WorkRecord.Fields["gidNumber"] = strconv.Itoa(c.GetNextGID())
		p.PrintPurple(fmt.Sprintf("\tOptional set groups's GID, press enter to use the next GID: %s\n",
			v.WorkRecord.Fields["gidNumber"]))
		fmt.Printf("\t%s: ", v.Template["gidNumber"].Prompt)
		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		if len(valueEntered) > 0 {
			gitNumberList := c.GetAlGroupsGID()
			if groupname, found := gitNumberList[valueEntered]; found {
				p.PrintRed(fmt.Sprintf("\n\tGiven group id %s already use by the group %s , aborting...\n",
					valueEntered, groupname))
				return false
			}
			v.WorkRecord.Fields["gidNumber"] = valueEntered
		}
	}
	return c.AddGroup()
}

func Create(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Create Group", 18, true))
	if createGroup(c) {
		p.PrintGreen(fmt.Sprintf("\tGroup %s created\n", v.WorkRecord.Fields["cn"]))
	} else {
		p.PrintRed(fmt.Sprintf("\tFailed to create the group %s, check the log file\n",
			 v.WorkRecord.Fields["cn"]))
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
