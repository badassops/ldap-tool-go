//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
package constants

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	Off			= "\033[0m"		// Text Reset
	Black		= "\033[1;30m"	// Black
	Red			= "\033[1;31m"	// Red
	Green		= "\033[1;32m"	// Green
	Yellow		= "\033[1;33m"	// Yellow
	Blue		= "\033[1;34m"	// Blue
	Purple		= "\033[1;35m"	// Purple
	Cyan		= "\033[1;36m"	// Cyan
	White		= "\033[1;37m"	// White
	CyanBase	= "\033[0;36m"	// Cyan no highlighted

	OK			= 0
	WARNING		= 1
	CRITICAL	= 2
	UNKNOWN		= 3

	ShadowMax		= 90
	ShadowWarning	= 14
)

type UserData struct {
	LoginName			string
	FirstName			string
	LastName			string
	UidNumber			string
	GidNumber			string
	Email				string
	Department			string
	Shell				string
	Password			string
	ShadowMax			int
	ShadowExpired		string
	ShadowWarning		int
	ShadowLastChange	string
	SSHPublicKey		string
	AdminGroups			[]string
	VPNGroups			[]string
	HomeDirectory		string
}

type LdapLog struct {
    LogsDir			string
    LogFile			string
    LogMaxSize		int
    LogMaxBackups	int
    LogMaxAge		int
}


var (
	MyVersion	= "0.1"
	now			= time.Now()
	MyProgname	= path.Base(os.Args[0])
	myAuthor	= "Luc Suryo"
	myCopyright = "Copyright 2008 - " + strconv.Itoa(now.Year()) + " ©Badassops LLC"
	myLicense	= "License 3-Clause BSD, https://opensource.org/licenses/BSD-3-Clause ♥"
	myEmail		= "<luc@badassops.com>"
	MyInfo = fmt.Sprintf("%s (version %s)\n%s\n%s\nWritten by %s %s\n",
		MyProgname, MyVersion, myCopyright, myLicense, myAuthor, myEmail)
	MyDescription = "Simple script to manage LDAP users"


	// users ldap record
	UserLdapData	UserData

	// ldap logs
	ServerLog		LdapLog

	// we sets these under variable
	LogsDir       = "/var/log/ldap-go"
	LogFile       = fmt.Sprintf("%s.log", MyProgname)
	LogMaxSize    = 128 // megabytes
	LogMaxBackups = 14  // 14 files
	LogMaxAge     = 14  // 14 days
)
