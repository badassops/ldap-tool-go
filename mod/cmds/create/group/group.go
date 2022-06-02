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
)

func createGroupRecord(c *l.Connection) bool {
	i := is.New()
	p := print.New()

	for _, fieldName := range fields {
		if v.GroupTemplate[fieldName].Value != "" {
			u.PrintYellow(fmt.Sprintf("\t ** Default to: %s **\n", v.GroupTemplate[fieldName].Value))
		}

		if c.Config.Debug {
			fmt.Printf("\t(%s) - %s: ", fieldName, v.GroupTemplate[fieldName].Prompt)
		} else {
			fmt.Printf("\t%s: ", v.GroupTemplate[fieldName].Prompt)
		}

		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")

		switch fieldName {
		case "groupName":
			cnt := c.CheckGroup(valueEntered)
			if cnt != 0 {
				u.PrintRed(fmt.Sprintf("\n\tGiven user %s already exist, aborting...\n\n", valueEntered))
				return false
			}
			u.PrintPurple(fmt.Sprintf("\tUsing Group: %s\n\n", valueEntered))
			c.Group["groupName"] = valueEntered
			c.Group["cn"] = fmt.Sprintf("cn=%s,%s", valueEntered, c.Config.ServerValues.GroupDN)

		case "groupType":
			switch valueEntered {
			case "", "p", "posix":
				valueEntered = "posix"
			case "g", "groupOfNames":
				valueEntered = "groupOfNames"
			}
			if !u.InList(validTypes, valueEntered) {
				u.PrintRed(fmt.Sprintf("\tWrong group type (%s) aborting...\n\n", valueEntered))
				return false
			}
			c.Group["groupType"] = valueEntered
		}

		if len(valueEntered) == 0 && v.GroupTemplate[fieldName].NoEmpty == true {
			u.PrintRed("\tNo value was entered aborting...\n\n")
			return false
		}
	}

	if c.Group["groupType"] == "posix" {
		c.Group["objectClass"] = "posixGroup"

		nextGID = c.GetNextGID()
		u.PrintPurple(fmt.Sprintf("\t\tOptional set groups's GID, press enter to use the next GID: %d\n", nextGID))
		fmt.Printf("\t%s: ", v.GroupTemplate["gidNumber"].Prompt)
		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")

		if len(valueEntered) == 0 && v.GroupTemplate["gidNumber"].NoEmpty == false {
			//u.PrintRed("\tNo value was entered aborting...\n\n")
			//return false
			valueEntered = strconv.Itoa(nextGID)
		} else {
			gidEntered, _ := strconv.Atoi(valueEntered)
			if found, groupName := c.CheckGroupID(gidEntered); found == true {
				u.PrintRed(fmt.Sprintf("\n\tGiven group id %v already use by the group %s, aborting...\n",
					valueEntered, groupName))
				return false
			}
		}
		c.Group["gidNumber"] = valueEntered
	}

	if c.Group["groupType"] == "groupOfNames" {
		c.Group["objectClass"] = "groupOfNames"
		c.Group["member"] = v.GroupTemplate["member"].Value
	}
	return true
}

func Create(c *l.Connection) {
	u.PrintHeader(u.Purple, "Create Group", true)
	if createGroupRecord(c) {
		u.PrintLine(u.Purple)
		if !c.AddGroup() {
			u.PrintRed(fmt.Sprintf("\n\tFailed to create the group %s, check the log file\n", c.Group["groupName"]))
		} else {
			u.PrintGreen(fmt.Sprintf("\n\tGroup %s created successfully\n", c.Group["groupName"]))
		}
	}
	u.PrintLine(u.Purple)
}
