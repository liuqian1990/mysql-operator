// Copyright 2018 Oracle and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build all default

package e2e

import (
	"testing"

	"github.com/oracle/mysql-operator/test/e2e/framework"
	e2eutil "github.com/oracle/mysql-operator/test/e2e/util"
)

func TestBackUpRestore(test *testing.T) {
	t := e2eutil.NewT(test)
	f := framework.Global

	t.Log("Creating mysqlcluster...")
	testdb := e2eutil.CreateTestDB(t, "e2e-br-", 1, false, f.DestroyAfterFailure)
	defer testdb.Delete()
	clusterName := testdb.Cluster().Name

	testdb.Populate()
	testdb.Test()

	databaseName := "employees"

	t.Logf("Creating mysqlbackup for mysqlcluster '%s'...", clusterName)
	backupName := e2eutil.Backup(t, f, clusterName, "e2e-br-backup-", databaseName)

	t.Log("Trying connection to container")
	testdb.CheckConnection(t)

	t.Log("Validating database..")
	testdb.Test()

	t.Logf("Deleting the %s database..", databaseName)
	e2eutil.DeleteDatabase(t, f, clusterName, databaseName)

	t.Logf("creating mysqlrestore from mysqlbackup '%s' for mysqlcluster '%s'.", backupName, clusterName)
	e2eutil.Restore(t, f, clusterName, backupName)

	t.Log("trying connection to container")
	testdb.CheckConnection(t)

	t.Log("validating database...")
	testdb.Test()

	t.Report()
}
