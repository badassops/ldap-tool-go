// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package modify

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	cu "badassops.ldap/cmds/common/user"
	l "badassops.ldap/ldap"
	u "badassops.ldap/utils"
	v "badassops.ldap/vars"
)

var (
	logRecord string

	fields = []string{"uidNumber", "givenName", "sn", "departmentNumber",
		"mail", "loginShell", "userPassword",
		"shadowMax", "shadowExpire", "sshPublicKey"}

	// given fiels value
	shells        string
	departments   string
	currShadowMax string
	userGroupList []string

	// list of the groups to be added or deleted
	addList []string
	delList []string

	// to show field name
	prefix string = ""

	// input
	valueEntered string

	// keep track if password was changed
	shadowMaxChanged bool = false
)

func createModifyUserRecord(c *l.Connection) {
	u.PrintPurple(fmt.Sprintf("\tUsing user: %s\n", c.User.Field["uid"]))
	u.PrintYellow(fmt.Sprintf("\tPress enter to leave the value unchanged\n"))
	u.PrintLine(u.Purple)

	for _, fieldName := range fields {
		// these will be valid once the field was filled since they depends
		// on some of the fields value
		switch fieldName {
		case "uidNumber":
			fmt.Printf("\t%s\n", u.DangerZone)
			u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", c.User.Field[fieldName]))

		case "givenName":
			u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", c.User.Field[fieldName]))

		case "sn":
			u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", c.User.Field[fieldName]))

		case "mail":
			u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", c.User.Field[fieldName]))

		case "departmentNumber":
			u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", c.User.Field[fieldName]))
			for _, value := range c.Config.GroupValues.Groups {
				departments = departments + " " + value
			}
			u.PrintPurple(fmt.Sprintf("\t\tValid departments:%s\n", departments))

		case "loginShell":
			u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", c.User.Field[fieldName]))
			for _, value := range c.Config.DefaultValues.ValidShells {
				shells = shells + " " + value
			}
			u.PrintPurple(fmt.Sprintf("\t\tValid shells:%s\n", shells))

		case "userPassword":
			passWord := u.GenerateRandom(
				c.Config.DefaultValues.PassComplex,
				c.Config.DefaultValues.PassLenght)
			u.PrintCyan(fmt.Sprintf("\tCurrent value (encrypted!): %s\n", c.User.Field[fieldName]))
			u.PrintYellow(fmt.Sprintf("\t\tsuggested password: %s\n", passWord))

		case "shadowMax":
			u.PrintCyan(fmt.Sprintf("\tCurrent max password age: %s\n", c.User.Field[fieldName]))
			u.PrintPurple(
				fmt.Sprintf("\t\tMin %d days and max %d days\n",
					c.Config.DefaultValues.ShadowMin,
					c.Config.DefaultValues.ShadowMax))

		case "shadowExpire":
			value, _ := strconv.ParseInt(c.User.Field["shadowExpire"], 10, 64)
			_, passExpired := u.GetReadableEpoch(value * 86400)
			u.PrintCyan(fmt.Sprintf("\tCurrent password will expire on: %s\n", passExpired))

		case "sshPublicKey":
			u.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", c.User.Field[fieldName]))

		}

		if c.Config.Debug {
			prefix = fmt.Sprintf("(%s) - ", fieldName)
		}

		if fieldName == "shadowExpire" {
			fmt.Printf("\t%sReset password expired to (%s days from now) Y/N: ",
				prefix, v.ModRecord.Field["shadowMax"])
		} else {
			fmt.Printf("\t%s%s: ", prefix, v.Template[fieldName].Prompt)
		}

		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		switch fieldName {
		case "givenName", "sn":
			valueEntered = strings.Title(valueEntered)

		case "mail":
			valueEntered = strings.ToLower(valueEntered)
			valueEntered = strings.ToLower(valueEntered)

		case "departmentNumber":
			if len(valueEntered) != 0 {
				if cnt := c.CheckGroup(valueEntered); cnt == 0 {
					u.PrintRed(fmt.Sprintf("\t\tGiven departmentNumber %s is not valid, given value ignored...\n",
						valueEntered))
					valueEntered = ""
				}
				for _, mapValues := range c.Config.GroupValues.GroupsMap {
					if mapValues.Name == valueEntered {
						v.ModRecord.Field["gidNumber"] = strconv.Itoa(mapValues.Gid)
					}
				}
				valueEntered = strings.ToUpper(valueEntered)
			}

		case "shadowMax":
			if len(valueEntered) != 0 {
				shadowMax, _ := strconv.Atoi(valueEntered)
				if shadowMax < c.Config.DefaultValues.ShadowMin ||
					shadowMax > c.Config.DefaultValues.ShadowMax {
					u.PrintRed(fmt.Sprintf("\t\tGiven value %d, is out or range, is set to %d\n",
						shadowMax, c.Config.DefaultValues.ShadowAge))
					valueEntered = strconv.Itoa(c.Config.DefaultValues.ShadowAge)
				}
				shadowMaxChanged = true
			}

		case "shadowExpire":
			if len(valueEntered) == 0 {
				u.PrintCyan(fmt.Sprintf("\tPassword expiration date will not be changed\n"))
			} else {
				// calculate when it will be expired based on default value if shadowMax
				// otherwise it will be today + new shadowMax value
				currShadowMax = c.User.Field["shadowMax"]
				if shadowMaxChanged == true {
					currShadowMax = v.ModRecord.Field["shadowMax"]
				}

				// set last changed to now
				v.ModRecord.Field["shadowLastChange"] = v.Template["shadowLastChange"].Value
				// calculate the new shadowExpire
				shadowLastChange, _ := strconv.ParseInt(v.ModRecord.Field["shadowLastChange"], 10, 64)
				shadowMax, _ := strconv.ParseInt(currShadowMax, 10, 64)
				_, passExpired := u.GetReadableEpoch((shadowLastChange + shadowMax) * 86400)
				u.PrintCyan(fmt.Sprintf("\tCurrent password will now expire on: %s\n", passExpired))
				// replace the 'Y' with the correct value
				valueEntered = strconv.FormatInt((shadowLastChange + shadowMax), 10)
			}
		}

		if len(valueEntered) != 0 {
			v.ModRecord.Field[fieldName] = valueEntered
		}
	}

	// we only handle groupOfNames type of group
	fmt.Printf("\n")
	// set the user's groups
	c.User.Groups = c.GetUserGroups("groupOfNames")

	// need only the group name
	reg, _ := regexp.Compile("^cn=|,ou=groups,.*")
	for _, userGroup := range c.User.Groups {
		userGroupList = append(userGroupList, fmt.Sprintf("%s", reg.ReplaceAllString(userGroup, "")))
	}

	availableGroups := c.GetGroupsNameByType("groupOfNames")
	u.PrintPurple(fmt.Sprintf("\t\tAvailable groups: %s\n", strings.Join(availableGroups[:], " ")))

	u.PrintPurple(fmt.Sprintf("\t\tUser %s groups: %s\n",
		c.User.Field["uid"], strings.Join(userGroupList, " ")))

	for _, leaveGroup := range userGroupList {
		u.PrintRed(fmt.Sprintf("\tRemove the group %s? default to not remove group, [Y/N]:  ", leaveGroup))
		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
		switch valueEntered {
		case "y", "Y", "yes", "YES", "d", "del", "D", "DEL":
			delList = append(delList, leaveGroup)
		}
	}

	for _, joinGroup := range availableGroups {
		if u.InList(userGroupList, joinGroup) == false {
			u.PrintGreen(fmt.Sprintf("\tJoin the group %s? default to not join group, [Y/N]:  ", joinGroup))
			reader := bufio.NewReader(os.Stdin)
			valueEntered, _ := reader.ReadString('\n')
			valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
			switch valueEntered {
			case "y", "Y", "yes", "YES", "a", "add", "A", "ADD":
				addList = append(addList, joinGroup)
			}
		}
	}
	v.ModRecord.AddList = addList
	v.ModRecord.DelList = delList
}

func Modify(c *l.Connection) {
	u.PrintHeader(u.Purple, "Modify User", true)
	if cu.User(c, true, false) {
		//if len(c.User.Field["uid"]) != 0 {
		createModifyUserRecord(c)
		if len(v.ModRecord.Field) == 0 &&
			len(v.ModRecord.AddList) == 0 &&
			len(v.ModRecord.DelList) == 0 {
			u.PrintBlue(fmt.Sprintf("\n\tNo field were changed, no modification made to the user %s\n",
				c.User.Field["uid"]))
		} else {
			if !c.ModifyUser() {
				u.PrintRed(fmt.Sprintf("\n\tFailed modify the user %s, check the log file\n", c.User.Field["uid"]))
			} else {
				u.PrintGreen(fmt.Sprintf("\n\tUser %s modified successfully\n", c.User.Field["uid"]))
			}
		}
		//}
	}
	u.PrintLine(u.Purple)
}
