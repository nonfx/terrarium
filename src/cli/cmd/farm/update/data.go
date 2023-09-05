// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package update

import (
	"regexp"
	"strings"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/rotisserie/eris"
)

func executeSQLStatement(db db.DB, sqlStatement string) error {
	return db.ExecuteSQLStatement(sqlStatement)
}

func seedDatabase(dumpContent string) error {
	dump := cleanup(dumpContent)
	g, err := config.DBConnect()
	if err != nil {
		return eris.Wrap(err, "failed to connect DB")
	}
	err = executeSQLStatement(g, dump)
	if err != nil {
		return eris.Wrap(err, "error executing dump file")
	}
	return nil
}

func cleanup(dump string) string {
	dump = strings.ReplaceAll(dump, "public.", "")
	setPattern := `(?m)^SET.*$`
	re := regexp.MustCompile(setPattern)
	dump = re.ReplaceAllString(dump, "")

	selectPattern := `(?m)^SELECT.*$`
	res := regexp.MustCompile(selectPattern)
	return res.ReplaceAllString(dump, "")
}
