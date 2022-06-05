//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package create

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/epoch"
	"github.com/badassops/packages-go/is"
	"github.com/badassops/packages-go/print"
	"github.com/badassops/packages-go/random"
	//ldapv3 "gopkg.in/ldap.v2"
)

var (
	// not required for create a new user : cn, gidNumber, displayName, gecos
	// homeDirectory, shadowLastChange, shadowLastChange
	// groups is handled seperat

	fields = []string{"uid", "givenName", "sn",
		"uidNumber", "departmentNumber",
		"mail", "loginShell", "userPassword",
		"shadowWarning", "shadowMax",
		"sshPublicKey"}

	// construct base on FirstName + LastName
	userFullname = []string{"cn", "displayName", "gecos"}

	// given field value
	email       string
	passWord    string
	shells      string
	departments string
	nextUID     int
	shadowMax   int

	e = epoch.New()
	i = is.New()
	p = print.New()
)

func joinGroup(c *l.Connection) int {
	var errored int = 0
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
		if !c.AddToGroup() {
			errored++
		}
	}
	return errored
}

func createUserRecord(c *l.Connection) bool {
	r := random.New(c.Config.DefaultValues.PassComplex, c.Config.DefaultValues.PassLenght)
	usersName := c.GetAllUsers()
	usersNameUid := c.GetUsersUID()
	reader := bufio.NewReader(os.Stdin)

	for _, fieldName := range fields {
		// these will be valid once the field was filled since they depends
		// on some of the fields value
		switch fieldName {
		case "uid":
			p.PrintPurple(
				fmt.Sprintf("\tThe userid / login name is case sensitive, it will be made all lowercase\n"))

		case "uidNumber":
			nextUID = c.GetNextUID()
			p.PrintPurple(fmt.Sprintf("\t\tOptional set user's UID, press enter to use the next UID: %d\n",
				nextUID))

		case "departmentNumber":
			p.PrintYellow(fmt.Sprintf("\t\tValid departments: %s\n",
				strings.Join(c.Config.GroupValues.Groups[:], ", ")))

		case "mail":
			email = fmt.Sprintf("%s.%s@%s",
				strings.ToLower(v.WorkRecord.Fields["givenName"]),
				strings.ToLower(v.WorkRecord.Fields["sn"]),
				c.Config.ServerValues.EmailDomain)
			p.PrintCyan(fmt.Sprintf("\tDefault email: %s\n", email))

		case "loginShell":
			p.PrintYellow(fmt.Sprintf("\t\tValid shells: %s\n",
				strings.Join(c.Config.DefaultValues.ValidShells[:], ", ")))

		case "userPassword":
			passWord = r.Generate()
			p.PrintPurple("\t\tPress Enter to accept the suggested password\n")
			p.PrintCyan(fmt.Sprintf("\tSuggested password: %s\n", passWord))

		case "shadowMax":
			p.PrintYellow(
				fmt.Sprintf("\t\tMin %d days and max %d days\n",
					c.Config.DefaultValues.ShadowMin,
					c.Config.DefaultValues.ShadowMax))
		}

		if v.Template[fieldName].Value != "" {
			p.PrintPurple(fmt.Sprintf("\t ** Default to: %s **\n", v.Template[fieldName].Value))
		}

		fmt.Printf("\t%s: ", v.Template[fieldName].Prompt)

		reader = bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")

		switch fieldName {
		case "uid":
			if i.IsInList(usersName, valueEntered) {
				p.PrintRed(fmt.Sprintf("\n\tGiven user %s already exist, aborting...\n\n", valueEntered))
				return false
			}
			fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
			v.WorkRecord.Fields["uid"] = valueEntered
			v.WorkRecord.ID = valueEntered
			p.PrintPurple(fmt.Sprintf("\tUsing user: %s\n", valueEntered))

		case "givenName", "sn":
			valueEntered = strings.Title(valueEntered)

		case "uidNumber":
			if len(valueEntered) > 0 {
				if userName, found := usersNameUid[valueEntered]; found {
					p.PrintRed(fmt.Sprintf("\n\tGiven uid id %s already use by the user %s , aborting...\n",
						valueEntered, userName))
					return false
				}
				valueEntered = valueEntered
			} else {
				valueEntered = strconv.Itoa(nextUID)
			}

		case "departmentNumber":
			if len(valueEntered) > 0 {
				if !i.IsInList(c.Config.GroupValues.Groups, valueEntered) {
					p.PrintRed(fmt.Sprintf("\n\tGiven departmentNumber %s is not valid, aborting...\n\n",
						valueEntered))
					return false
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
			if len(valueEntered) == 0 {
				valueEntered = strings.ToUpper(c.Config.DefaultValues.GroupName)
				v.WorkRecord.Fields["gidNumber"] = strconv.Itoa(c.Config.DefaultValues.GroupId)
			}

		case "mail":
			if len(valueEntered) == 0 {
				valueEntered = email
			}

		case "loginShell":
			if len(valueEntered) > 0 {
				if !i.IsInList(c.Config.DefaultValues.ValidShells, valueEntered) {
					p.PrintRed(fmt.Sprintf("\n\tGiven shell %s is not valid, aborting...\n\n", valueEntered))
					return false
				}
				valueEntered = "/bin/" + valueEntered
			}
			if len(valueEntered) == 0 {
				valueEntered = "/bin/" + c.Config.DefaultValues.Shell
			}

		case "userPassword":
			if len(valueEntered) > 0 {
				v.WorkRecord.Fields["userPassword"] = valueEntered
			}
			if len(valueEntered) == 0 {
				valueEntered = passWord
			}

		case "shadowWarning":
			if len(valueEntered) == 0 {
				valueEntered = strconv.Itoa(c.Config.DefaultValues.ShadowWarning)
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
			}
			if len(valueEntered) == 0 {
				valueEntered = strconv.Itoa(c.Config.DefaultValues.ShadowAge)
			}

		default:
			if len(valueEntered) == 0 {
				valueEntered = v.Template[fieldName].Value
			}

		}

		if len(valueEntered) == 0 && v.Template[fieldName].NoEmpty == true {
			p.PrintRed("\tNo value was entered aborting...\n\n")
			return false
		}
		// set the default values
		if len(valueEntered) == 0 {
			v.WorkRecord.Fields[fieldName] = v.Template[fieldName].Value
		}
		// update the user record so it can be submitted
		v.WorkRecord.Fields[fieldName] = valueEntered
	}

	for idx, _ := range userFullname {
		v.WorkRecord.Fields[userFullname[idx]] =
			v.WorkRecord.Fields["givenName"] + " " + v.WorkRecord.Fields["sn"]
	}

	// dn is create base on given uid and user DN
	v.WorkRecord.DN = fmt.Sprintf("uid=%s,%s", v.WorkRecord.Fields["uid"], c.Config.ServerValues.UserDN)

	// this is always /home + userlogin
	v.WorkRecord.Fields["homeDirectory"] = "/home/" + v.WorkRecord.Fields["uid"]

	// initialized to be today's epoch days
	v.WorkRecord.Fields["shadowExpire"] = v.Template["shadowExpire"].Value
	v.WorkRecord.Fields["shadowLastChange"] = v.Template["shadowLastChange"].Value

	// setup the groups for the user
	fmt.Printf("\n\t%s\n", p.PrintLine(v.Purple, 50))
	fmt.Printf("\n")
	reg, _ := regexp.Compile("^cn=|,ou=groups,.*")
	availableGroups := c.GetAllGroups()
	for _, joinGroup := range availableGroups {
		groupToJoin := reg.ReplaceAllString(joinGroup, "")
		if v.WorkRecord.Fields["departmentNumber"] != strings.ToUpper(groupToJoin) {
			fmt.Printf("\t%sJoin%s the group %s%s%s? default not to join group, [Y/N]: ",
				v.Green, v.Off, v.Green, groupToJoin, v.Off)
			valueEntered, _ := reader.ReadString('\n')
			valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
			switch valueEntered {
			case "y", "yes":
				c.Record.GroupAddList = append(c.Record.GroupAddList, joinGroup)
			}
		} else {
			c.Record.GroupAddList = append(c.Record.GroupAddList, groupToJoin)
		}
	}
	return c.CreateUser()
}

func Create(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Create User", 18, true))
	if !createUserRecord(c) {
		p.PrintRed(fmt.Sprintf("\n\tFailed adding the user %s, check the log file\n", v.WorkRecord.Fields["uid"]))
	} else {
		p.PrintGreen(fmt.Sprintf("\n\tUser %s added successfully\n", v.WorkRecord.Fields["uid"]))
		if len(c.Record.GroupAddList) > 0 {
			if joinGroup(c) != 0 {
				p.PrintRed(fmt.Sprintf("\n\tFailed adding the user %s to groups, check the log file\n",
					v.WorkRecord.Fields["uid"]))
			}
		}
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
