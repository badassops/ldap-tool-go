//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//
package vars

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

type StringRecord struct {
	Value	string
	Changed	bool
}

type IntRecord struct {
	Value	int
	Changed	bool
}

type UserRecord struct {
	DN		string
	Strings	map[string]StringRecord
	Ints	map[string]IntRecord
	Groups	[]string
}

type RecordField struct {
	FieldName	string
	Default		string
	Prompt		string
	NoEmpty		bool
	UseDefault	bool
}

type Log struct {
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
	myCopyright = "Copyright 2019 - " + strconv.Itoa(now.Year()) + " ©Badassops LLC"
	myLicense	= "License 3-Clause BSD, https://opensource.org/licenses/BSD-3-Clause ♥"
	myEmail		= "<luc@badassops.com>"
	MyInfo = fmt.Sprintf("%s (version %s)\n%s\n%s\nWritten by %s %s\n",
		MyProgname, MyVersion, myCopyright, myLicense, myAuthor, myEmail)
	MyDescription = "Simple script to manage LDAP users"

	// users ldap record
	User	UserRecord

	// ldap logs
	Logs	Log

	// the ldap fields
	Fields	[]string

	// for user record (create and modify)
	RecordFields	map[int]RecordField

	// we sets these under variable
	LogsDir			= "/var/log/ldap-go"
	LogFile			= fmt.Sprintf("%s.log", MyProgname)
	LogMaxSize		= 128 // megabytes
	LogMaxBackups	= 14  // 14 files
	LogMaxAge		= 14  // 14 days
)
