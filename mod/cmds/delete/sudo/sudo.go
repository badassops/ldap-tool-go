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

func Delete(c *l.Connection) {
	fmt.Printf("\t%s\n", p.PrintHeader(v.Blue, v.Purple, "Delete Sudo rules", 18, true))
	v.SearchResultData.WildCardSearchBase = v.SudoWildCardSearchBase
	v.SearchResultData.RecordSearchbase = v.SudoWildCardSearchBase
	v.SearchResultData.DisplayFieldID = v.SudoDisplayFieldID
	if common.GetObjectRecord(c, true, "sudo rules") {
		common.DeleteObjectRecord(c, v.SearchResultData.SearchResult, "sudo rules")
	}
	fmt.Printf("\t%s\n", p.PrintLine(v.Purple, 50))
}
