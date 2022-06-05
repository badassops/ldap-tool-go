//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package modify

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"badassops.ldap/cmds/common"
	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/epoch"
	"github.com/badassops/packages-go/is"
	"github.com/badassops/packages-go/print"
	"github.com/badassops/packages-go/random"
	ldapv3 "gopkg.in/ldap.v2"
)

var (
	fields = []string{"uidNumber", "givenName", "sn", "departmentNumber",
		"mail", "loginShell", "userPassword",
		"shadowMax", "shadowExpire", "sshPublicKey"}

	// input
	valueEntered string

	// user's groups
	userGroupList      []string
	availableGroupList []string
	// need to strip the full dn
	displayUserGroupList      []string
	displayAvailableGroupList []string

	// keep track if password was changed
	shadowMaxChanged bool = false

	e = epoch.New()
	i = is.New()
	p = print.New()
)

func leaveGroup(c *l.Connection) {
	groupList := c.GetGroupType()
	for _, groupName := range c.Record.GroupDelList {
		if i.IsInList(groupList["posixGroup"], groupName) {
			v.WorkRecord.MemberType = "posixGroup"
			v.WorkRecord.MemberType = "memberUid"
			v.WorkRecord.ID = v.WorkRecord.ID
		}
		if i.IsInList(groupList["groupOfNames"], groupName) {
			v.WorkRecord.MemberType = "groupOfNames"
			v.WorkRecord.MemberType = "member"
			v.WorkRecord.ID = fmt.Sprintf("uid=%s,%s", v.WorkRecord.ID, c.Config.ServerValues.UserDN)
		}
		v.WorkRecord.DN = groupName
		c.RemoveFromGroups()
	}
}

func joinGroup(c *l.Connection) {
	groupList := c.GetGroupType()
	for _, groupName := range c.Record.GroupAddList {
		if i.IsInList(groupList["posixGroup"], groupName) {
			v.WorkRecord.MemberType = "posixGroup"
			v.WorkRecord.MemberType = "memberUid"
			v.WorkRecord.ID = v.WorkRecord.ID
		}
		if i.IsInList(groupList["groupOfNames"], groupName) {
			v.WorkRecord.MemberType = "groupOfNames"
			v.WorkRecord.MemberType = "member"
			v.WorkRecord.ID = fmt.Sprintf("uid=%s,%s", v.WorkRecord.ID, c.Config.ServerValues.UserDN)
		}
		v.WorkRecord.DN = groupName
		c.AddToGroup()
	}
}

func createModifyUserRecord(c *l.Connection, records *ldapv3.SearchResult) int {
	r := random.New(c.Config.DefaultValues.PassComplex, c.Config.DefaultValues.PassLenght)
	reader := bufio.NewReader(os.Stdin)

	v.WorkRecord.DN = fmt.Sprintf("uid=%s,%s", v.WorkRecord.ID, c.Config.ServerValues.UserDN)

	p.PrintPurple(fmt.Sprintf("\tUsing user: %s\n", v.WorkRecord.ID))
	p.PrintYellow(fmt.Sprintf("\tPress enter to leave the value unchanged\n"))
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))

	for _, fieldName := range fields {
		// these will be valid once the field was filled since they depends
		// on some of the fields value
		switch fieldName {
		case "uidNumber":
			fmt.Printf("\t%s\n", v.DangerZone)
			p.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				v.Green, records.Entries[0].GetAttributeValue(fieldName), v.Off))

		case "givenName":
			p.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				v.Green, records.Entries[0].GetAttributeValue(fieldName), v.Off))

		case "sn":
			p.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				v.Green, records.Entries[0].GetAttributeValue(fieldName), v.Off))

		case "mail":
			p.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				v.Green, records.Entries[0].GetAttributeValue(fieldName), v.Off))

		case "departmentNumber":
			p.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				v.Green, records.Entries[0].GetAttributeValue(fieldName), v.Off))
			p.PrintYellow(fmt.Sprintf("\t\tValid departments: %s\n",
				strings.Join(c.Config.GroupValues.Groups[:], ", ")))

		case "loginShell":
			p.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				v.Green, records.Entries[0].GetAttributeValue(fieldName), v.Off))
			p.PrintYellow(fmt.Sprintf("\t\tValid shells: %s\n",
				strings.Join(c.Config.DefaultValues.ValidShells[:], ", ")))

		case "userPassword":
			passWord := r.Generate()
			p.PrintCyan(fmt.Sprintf("\tCurrent value (encrypted!): %s%s%s\n",
				v.Green, records.Entries[0].GetAttributeValue(fieldName), v.Off))
			p.PrintYellow(fmt.Sprintf("\t\tsuggested password: %s\n", passWord))

		case "shadowMax":
			p.PrintCyan(fmt.Sprintf("\tCurrent max password age: %s%s%s\n",
				v.Green, records.Entries[0].GetAttributeValue(fieldName), v.Off))
			p.PrintYellow(
				fmt.Sprintf("\t\tMin %d days and max %d days\n",
					c.Config.DefaultValues.ShadowMin,
					c.Config.DefaultValues.ShadowMax))

		case "shadowExpire":
			value, _ := strconv.ParseInt(records.Entries[0].GetAttributeValue(fieldName), 10, 64)
			passExpired := e.ReadableEpoch(value * 86400)
			p.PrintCyan(fmt.Sprintf("\tCurrent password will expire on: %s%s%s\n",
				v.Green, passExpired, v.Off))

		case "sshPublicKey":
			p.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				v.Green, records.Entries[0].GetAttributeValue(fieldName), v.Off))
		}

		p.PrintPurple(fmt.Sprintf("\t%s: ", v.Template[fieldName].Prompt))

		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		switch fieldName {
		case "givenName", "sn":
			valueEntered = strings.Title(valueEntered)

		case "mail":
			valueEntered = strings.ToLower(valueEntered)
			valueEntered = strings.ToLower(valueEntered)

		case "loginShell":
			if len(valueEntered) > 0 {
				if !i.IsInList(c.Config.DefaultValues.ValidShells, valueEntered) {
					p.PrintRed("\t\tInvalid shell was given, it will be ignored\n")
					valueEntered = ""
				}
			}

		case "departmentNumber":
			if len(valueEntered) != 0 {
				if !i.IsInList(c.Config.GroupValues.Groups, valueEntered) {
					p.PrintRed("\t\tInvalid departments was given, it will be ignored\n")
					valueEntered = ""
				} else {
					for _, mapValues := range c.Config.GroupValues.GroupsMap {
						if mapValues.Name == valueEntered {
							v.WorkRecord.Fields["gidNumber"] = strconv.Itoa(mapValues.Gid)
							break
						}
					}
					valueEntered = strings.ToUpper(valueEntered)
				}
			}

		case "shadowMax":
			if len(valueEntered) != 0 {
				shadowMax, _ := strconv.Atoi(valueEntered)
				if shadowMax < c.Config.DefaultValues.ShadowMin ||
					shadowMax > c.Config.DefaultValues.ShadowMax {
					p.PrintRed(fmt.Sprintf("\t\tGiven value %d, is out or range, is set to %d\n",
						shadowMax, c.Config.DefaultValues.ShadowAge))
					valueEntered = strconv.Itoa(c.Config.DefaultValues.ShadowAge)
				}
				shadowMaxChanged = true
			}

		case "shadowExpire":
			if len(valueEntered) == 0 {
				p.PrintCyan(fmt.Sprintf("\tPassword expiration date will not be changed\n"))
			} else {
				// calculate when it will be expired based on default value if shadowMax
				// otherwise it will be today + new shadowMax value
				currShadowMax := records.Entries[0].GetAttributeValue("shadowMax")
				if shadowMaxChanged == true {
					currShadowMax = v.WorkRecord.Fields["shadowMax"]
					p.PrintYellow(fmt.Sprintf("\t\tCalculate from new given max password age\n"))
				}

				// set last changed to now
				v.WorkRecord.Fields["shadowLastChange"] = v.Template["shadowLastChange"].Value
				// calculate the new shadowExpire
				shadowLastChange, _ := strconv.ParseInt(v.WorkRecord.Fields["shadowLastChange"], 10, 64)
				shadowMax, _ := strconv.ParseInt(currShadowMax, 10, 64)
				passExpired := e.ReadableEpoch((shadowLastChange + shadowMax) * 86400)
				p.PrintCyan(fmt.Sprintf("\tCurrent password will now expire on: %s\n", passExpired))
				// replace the 'Y' with the correct value
				valueEntered = strconv.FormatInt((shadowLastChange + shadowMax), 10)
			}
		}

		if len(valueEntered) != 0 {
			v.WorkRecord.Fields[fieldName] = valueEntered
		}
	}

	fmt.Printf("\n")
	userID := v.WorkRecord.ID
	userDN := records.Entries[0].DN
	c.GetUserGroups(userID, userDN)
	c.GetAvailableGroups(userID, userDN)

	fmt.Printf("\n\t%s\n", p.PrintLine(v.Purple, 50))
	reg, _ := regexp.Compile("^cn=|,ou=groups,.*")

	for _, userGroup := range c.Record.UserGroups {
		userGroupList = append(userGroupList, userGroup)
		displayUserGroupList = append(displayUserGroupList, reg.ReplaceAllString(userGroup, " "))
	}
	for _, availableGroup := range c.Record.AvailableGroups {
		availableGroupList = append(availableGroupList, availableGroup)
		displayAvailableGroupList = append(displayAvailableGroupList, reg.ReplaceAllString(availableGroup, " "))
	}

	p.PrintPurple(fmt.Sprintf("\tUser %s groups: %s\n", v.WorkRecord.ID,
		strings.Join(displayUserGroupList[:], "")))

	p.PrintGreen(fmt.Sprintf("\tAvailable groups: %s\n",
		strings.Join(displayAvailableGroupList[:], " ")))

	fmt.Printf("\n\t%s\n", p.PrintLine(v.Purple, 50))
	for _, leaveGroup := range userGroupList {
		fmt.Printf("\t%sRemove%s the group %s%s%s? default to not remove group, [Y/N]: ",
			v.Red, v.Off, v.Red, reg.ReplaceAllString(leaveGroup, ""), v.Off)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
		switch valueEntered {
		case "y", "yes", "d", "del":
			c.Record.GroupDelList = append(c.Record.GroupDelList, leaveGroup)
		}
	}

	fmt.Printf("\n")
	for _, joinGroup := range availableGroupList {
		fmt.Printf("\t%sJoin%s the group %s%s%s? default not to join group, [Y/N]: ",
			v.Green, v.Off, v.Green, reg.ReplaceAllString(joinGroup, ""), v.Off)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
		switch valueEntered {
		case "y", "yes" :
			c.Record.GroupAddList = append(c.Record.GroupAddList, joinGroup)
		}
	}

	if len(c.Record.GroupDelList) == 0 && len(c.Record.GroupAddList) == 0 && len(v.WorkRecord.Fields) == 0 {
		fmt.Printf("\n\t%s\n", p.PrintLine(v.Purple, 50))
		p.PrintBlue(fmt.Sprintf("\n\tNo field were changed, no modification was made for the user %s\n",
			v.WorkRecord.ID))
		return 0
	}
	return 1
}

func Modify(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Modify user", 18, true))
	v.SearchResultData.WildCardSearchBase = v.UserWildCardSearchBase
	v.SearchResultData.RecordSearchbase = v.UserWildCardSearchBase
	v.SearchResultData.DisplayFieldID = v.UserDisplayFieldID
	if common.GetObjectRecord(c, true, "user") {
		if createModifyUserRecord(c, v.SearchResultData.SearchResult) > 0 {
			if len(v.WorkRecord.Fields) > 0 {
				if !c.ModifyUser() {
					p.PrintRed(fmt.Sprintf("\n\tFailed modify the user %s, check the log file\n", v.WorkRecord.ID))
				} else {
					p.PrintGreen(fmt.Sprintf("\n\tUser %s modified successfully\n", v.WorkRecord.ID))
				}
			}
			if len(c.Record.GroupDelList) > 0 {
				leaveGroup(c)
			}
			if len(c.Record.GroupAddList) > 0 {
				joinGroup(c)
			}
		}
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
