//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package ldap

import (
	"github.com/badassops/packages-go/exit"
	"github.com/badassops/packages-go/lock"
	ldapv3 "gopkg.in/ldap.v2"
)

// search the ldap database
func (c *Connection) Search() (*ldapv3.SearchResult, int) {
	searchRecords := ldapv3.NewSearchRequest(
		c.Config.ServerValues.BaseDN,
		ldapv3.ScopeWholeSubtree,
		ldapv3.NeverDerefAliases, 0, 0, false,
		c.SearchInfo.SearchBase,
		c.SearchInfo.SearchAttribute,
		nil,
	)
	searchResult, err := c.Conn.Search(searchRecords)
	if err != nil {
		c.Conn.Close()
		l := lock.New(c.Config.DefaultValues.LockFile)
		e := exit.New("ldap search", 1)
		l.LockRelease()
		e.ExitError(err)
	}
	return searchResult, len(searchResult.Entries)
}
