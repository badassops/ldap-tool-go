// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version		:	0.1
//

package ldap

import (
	"fmt"
	"net"
	"strings"
	"strconv"
	"time"

	"badassops.ldap/utils"
	"badassops.ldap/configurator"

	ldapv3 "gopkg.in/ldap.v2"
)

type (
	Record struct {
		Data	string
		Changed	bool
	}

	MemberOf struct {
		Data	[]string
		Changed	bool
	}

	UserRecord struct {
		DN					string
		UserName			string
		FirstName			Record
		LastName			Record
		Email				Record
		UID					int
		GID					Record
		GroupName			string
		Groups				[]string
		Shell				Record
		HomeDir				string
		Password			Record
		SSHPublicKey		Record
		AdminGroups			MemberOf
		VPNGroups			MemberOf
		ShadowMax			Record
		ShadowLastChange	Record
		ShadowExpire		Record
		ShadowWarning		Record
	}

	Connection struct {
		Conn		*ldapv3.Conn
		User		UserRecord
		Config		*configurator.Config
		LockFile	string
		LockPid		int
	}
)

// function to initialize a user record
func New(config *configurator.Config) *Connection {
	// set variable for the ldap connection
	var ppolicy *ldapv3.ControlBeheraPasswordPolicy

	// check if we can search the server
	timeout := time.Second
	dialConn, err := net.DialTimeout("tcp", net.JoinHostPort(config.Server, "389"), timeout)
	if err != nil {
		utils.ReleaseIT(config.LockFile, config.LockPID)
		utils.ExitWithMesssage(err.Error())
	}
	dialConn.Close()

	// set to expire by default as today + ShadowMax
	currExpired := strconv.FormatInt(utils.GetEpoch("days") + int64(config.ShadowMax), 10)

	ServerConn, err := ldapv3.Dial("tcp", fmt.Sprintf("%s:%d", config.Server ,389))
	if err != nil {
		utils.ReleaseIT(config.LockFile, config.LockPID)
		utils.ExitWithMesssage(err.Error())
	}

	// now we need to reconnect with TLS
	// if config.TLS {
	// 	err := ServerConn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	// 	if err != nil {
	//		utils.ReleaseIT(config.LockFile, config.LockPID)
	// 		utils.ExitIfError(err)
	// 	}
	// }

	// setup control
	controls := []ldapv3.Control{}
	controls = append(controls, ldapv3.NewControlBeheraPasswordPolicy())

	// bind to the ldap server
	bindRequest := ldapv3.NewSimpleBindRequest(config.Admin, config.AdminPass, controls)
	request, err := ServerConn.SimpleBind(bindRequest)
	ppolicyControl := ldapv3.FindControl(request.Controls, ldapv3.ControlTypeBeheraPasswordPolicy)
	if ppolicyControl != nil {
		ppolicy = ppolicyControl.(*ldapv3.ControlBeheraPasswordPolicy)
	 }
	if err != nil {
		errStr := "ERROR: Cannot bind: " + err.Error()
		if ppolicy != nil && ppolicy.Error >= 0 {
			errStr += ":" + ppolicy.ErrorString
		}
		utils.ReleaseIT(config.LockFile, config.LockPID)
		utils.ExitWithMesssage(errStr)
	}

	// the rest of the values will be filled during the process
    return &Connection {
		Conn:			ServerConn,
		User: UserRecord{
			Shell:			Record{ Data: config.Shell,							Changed: false },
			GID:			Record{ Data: strconv.Itoa(config.GroupId),			Changed: false },
			GroupName:		config.GroupName,
			ShadowMax:		Record{ Data: strconv.Itoa(config.ShadowMax),		Changed: false },
			ShadowWarning:	Record{ Data: strconv.Itoa(config.ShadowWarning),	Changed: false },
			ShadowExpire:	Record{ Data: currExpired,							Changed: false },
		},
		Config:				config,
		LockFile:			config.LockFile,
        LockPid:			config.LockPID,
	}
}

func (c *Connection) CheckUser(wildCard bool, user string) (*ldapv3.SearchResult, int) {
	attributes := []string{}

	orginalUser := user
	// wildcard should only return the DN value
	if wildCard {
		user = "*" + user + "*"
		//attributes = append(attributes, "uid")
	}

	searchBase :=fmt.Sprintf("(&(objectClass=inetOrgPerson)(uid=%s))", user)
	records, cnt := c.search(searchBase, attributes)

	if cnt == 0  {
		return records, 0
	}
	if cnt == 1 {
		UID, _ := strconv.Atoi(records.Entries[0].GetAttributeValue("uidNumber"))
		c.User.DN				= records.Entries[0].DN
		c.User.UserName			= records.Entries[0].GetAttributeValue("uid")
		c.User.FirstName		= Record{records.Entries[0].GetAttributeValue("givenName"),			false}
		c.User.LastName			= Record{records.Entries[0].GetAttributeValue("sn"),				false}
		c.User.Email			= Record{records.Entries[0].GetAttributeValue("mail"),				false}
		c.User.UID				= UID
		c.User.GID				= Record{records.Entries[0].GetAttributeValue("gidNumber"),			false}
		c.User.GroupName		= records.Entries[0].GetAttributeValue("departmentNumber")
		c.User.Shell			= Record{records.Entries[0].GetAttributeValue("loginShell"),		false}
		c.User.HomeDir			= records.Entries[0].GetAttributeValue("homeDirectory")
		c.User.Password			= Record{records.Entries[0].GetAttributeValue("userPassword"),		false}
		c.User.SSHPublicKey		= Record{records.Entries[0].GetAttributeValue("sshPublicKey"),		false}
		c.User.ShadowMax		= Record{records.Entries[0].GetAttributeValue("shadowMax"),			false}
		c.User.ShadowLastChange	= Record{records.Entries[0].GetAttributeValue("shadowLastChange"),	false}
		c.User.ShadowExpire		= Record{records.Entries[0].GetAttributeValue("shadowExpire"),		false}
		c.User.ShadowWarning	= Record{records.Entries[0].GetAttributeValue("shadowWarning"),		false}
		c.userGroups(orginalUser)
	 }
	return records, cnt
}

func (c *Connection) SearchUsers() {
	attributes := []string{}
	searchBase := fmt.Sprintf("(objectClass=person)")
	records, cnt := c.search(searchBase, attributes)
	if cnt > 0 {
		for idx, entry := range records.Entries {
			utils.PrintColor(utils.Blue, fmt.Sprintf("\tdn: %s\n", entry.DN))
			utils.PrintColor(utils.Green, fmt.Sprintf("\t\tFull Name: %s %s\n",
				records.Entries[idx].GetAttributeValue("givenName"),
				records.Entries[idx].GetAttributeValue("sn")))
			utils.PrintColor(utils.Green, fmt.Sprintf("\t\tdepartmentNumber: %s \n",
				records.Entries[idx].GetAttributeValue("departmentNumber")))
		}
		utils.PrintColor(utils.Yellow, fmt.Sprintf("\n\tTotal records: %d \n", cnt))
	}
}

func (c *Connection) SearchGroup(group string, all bool) {
	var records *ldapv3.SearchResult
	var cnt int
	var searchBase string
	if all {
		searchBase = fmt.Sprintf("(objectClass=posixGroup)")
	} else {
		searchBase = fmt.Sprintf("(&(objectClass=posixGroup)(cn=%s))", group)
	}
	records, cnt = c.search(searchBase, []string{"cn", "gidNumber", "memberUid"})
	if cnt > 0 {
		for _, entry := range records.Entries {
			fmt.Printf("dn: %s \n", entry.DN)
			if ! all {
				fmt.Printf("cn: %s \n", entry.GetAttributeValue("cn"))
				fmt.Printf("gidNumber: %s \n", entry.GetAttributeValue("gidNumber"))
			}
			for _, member := range entry.GetAttributeValues("memberUid") {
				utils.PrintColor(utils.Blue, fmt.Sprintf("memberUid: %s \n", member))
			}
			fmt.Printf("\n")
		}
		if all == true {
			utils.PrintColor(utils.Yellow, fmt.Sprintf("\n\tTotal records: %d \n", cnt))
		}
	}
}

func (c *Connection) userGroups(user string) {

	searchBase := fmt.Sprintf("(&(objectClass=groupOfNames)(member=uid=%s,%s))", user, c.Config.UserDN)
	records, cnt := c.search(searchBase, []string{"dn"})
	if cnt > 0 {
		var adminGroups []string
		var vpnGroups []string
		for _, entry := range records.Entries {
			if strings.Contains(entry.DN, "vpn") {
				vpnGroups = append(vpnGroups, entry.DN)
			} else {
				adminGroups = append(adminGroups, entry.DN)
			}
		}
		c.User.AdminGroups = MemberOf{adminGroups, false}
		c.User.VPNGroups = MemberOf{vpnGroups, false}
	}
}

func (c *Connection) search(searchBase string, searchAttribute []string) (*ldapv3.SearchResult, int) {
	searchRecords := ldapv3.NewSearchRequest(
		c.Config.BaseDN,
		ldapv3.ScopeWholeSubtree,
		ldapv3.NeverDerefAliases, 0, 0, false,
		searchBase,
		searchAttribute,
		nil,
	)
	sr, err := c.Conn.Search(searchRecords)
	if err != nil {
		c.Conn.Close()
		utils.ReleaseIT(c.LockFile, c.LockPid)
		utils.ExitWithMesssage(err.Error())
	}
	if len(sr.Entries) > 0 {
		return sr, len(sr.Entries)
	}
	return sr, 0
}
