//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package ldap

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	c "badassops.ldap/configurator"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/exit"
	"github.com/badassops/packages-go/lock"
	ldapv3 "gopkg.in/ldap.v2"
)

type (
	Connection struct {
		Conn       *ldapv3.Conn
		Config     *c.Config
		Record     v.LdapRecord
		SearchInfo v.SearchInfo
	}
)

// function to initialize the ldap system
func New(config *c.Config) *Connection {
	e := exit.New("ldap initialize", 1)
	l := lock.New(config.DefaultValues.LockFile)

	// set variable for the ldap connection
	var ppolicy *ldapv3.ControlBeheraPasswordPolicy

	// check if we can search the server, timeout set to 15 seconds
	timeout := 15 * time.Second
	dialConn, err := net.DialTimeout("tcp", net.JoinHostPort(config.ServerValues.Server, "389"), timeout)
	if err != nil {
		l.LockRelease()
		e.ExitError(err)
	}
	dialConn.Close()

	ServerConn, err := ldapv3.Dial("tcp", fmt.Sprintf("%s:%d", config.ServerValues.Server, 389))
	if err != nil {
		l.LockRelease()
		e.ExitError(err)
	}

	// now we need to reconnect with TLS
	if config.ServerValues.TLS {
		err := ServerConn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			l.LockRelease()
			e.ExitError(err)
		}
	}

	// setup control
	controls := []ldapv3.Control{}
	controls = append(controls, ldapv3.NewControlBeheraPasswordPolicy())

	// bind to the ldap server
	bindRequest := ldapv3.NewSimpleBindRequest(config.ServerValues.Admin, config.ServerValues.AdminPass, controls)
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
		l.LockRelease()
		e.ExitError(err)
	}

	// the rest of the values will be filled during the process
	return &Connection{
		Conn:   ServerConn,
		Config: config,
	}
}
