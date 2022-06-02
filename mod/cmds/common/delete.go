// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version    :  0.1
//

package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/readinput"
	ldapv3 "gopkg.in/ldap.v2"
)

func DeleteObjectRecord(c *l.Connection, records *ldapv3.SearchResult, objectType string) {
	reader := bufio.NewReader(os.Stdin)

	switch objectType {
	case "user":
		objectID = records.Entries[0].GetAttributeValue("uid")
		protectedList = c.Config.GroupValues.Groups
		v.WorkRecord.DN = fmt.Sprintf("uid=%s,%s", objectID, c.Config.ServerValues.UserDN)
		v.WorkRecord.ID = objectID

	case "group":
		objectID = records.Entries[0].GetAttributeValue("cn")
		protectedList = append(c.Config.GroupValues.Groups, c.Config.GroupValues.SpecialGroups...)
		v.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", objectID, c.Config.ServerValues.GroupDN)
		v.WorkRecord.ID = objectID

	case "sudo rules":
		objectID = records.Entries[0].GetAttributeValue("cn")
		protectedList = c.Config.SudoValues.ExcludeSudo
		v.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", objectID, c.Config.SudoValues.SudoersBase)
		v.WorkRecord.ID = objectID
	}

	if objectType != "user" {
		if i.IsInList(protectedList, objectID) {
			p.PrintRed(fmt.Sprintf("\n\tGiven %s %s is protected and can not be deleted, aborting...\n\n",
				objectType, objectID))
			return
		}
	}

	p.PrintRed(fmt.Sprintf("\n\tGiven %s %s will be delete, this can not be undo!\n", objectType, objectID))
	p.PrintYellow(fmt.Sprintf("\tContinue (default to N)? [y/n]: "))
	continueDelete, _ := reader.ReadString('\n')
	continueDelete = strings.ToLower(strings.TrimSuffix(continueDelete, "\n"))
	if readinput.ReadYN(continueDelete, false) == true {
		if !c.Delete(objectID, objectType) {
			p.PrintRed(fmt.Sprintf("\n\tFailed to delete the %s %s, check the log file\n", objectType, objectID))
		} else {
			p.PrintGreen(fmt.Sprintf("\n\tGiven %s %s has been deleted\n", objectType, objectID))
		}
	} else {
		p.PrintBlue(fmt.Sprintf("\n\tDeletion of the %s %s cancelled\n", objectType, objectID))
	}
	return
}
