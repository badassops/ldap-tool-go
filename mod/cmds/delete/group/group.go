//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package delete

import (
	"fmt"

	"badassops.ldap/cmds/common"
	l "badassops.ldap/ldap"
	v "badassops.ldap/vars"
	"github.com/badassops/packages-go/print"
)

var (
	p = print.New()
)

// delete a ldap record
func Delete(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Delete Group", 18, true))
	v.SearchResultData.WildCardSearchBase = v.GroupWildCardSearchBase
	v.SearchResultData.RecordSearchbase = v.GroupWildCardSearchBase
	v.SearchResultData.DisplayFieldID = v.GroupDisplayFieldID
	if common.GetObjectRecord(c, true, "group") {
		common.DeleteObjectRecord(c, v.SearchResultData.SearchResult, "group")
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
