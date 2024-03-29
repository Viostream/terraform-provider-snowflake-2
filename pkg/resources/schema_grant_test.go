package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/viostream/terraform-provider-snowflake/pkg/provider"
	"github.com/viostream/terraform-provider-snowflake/pkg/resources"
	. "github.com/viostream/terraform-provider-snowflake/pkg/testhelpers"
)

func TestSchemaGrant(t *testing.T) {
	r := require.New(t)
	err := resources.SchemaGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestSchemaGrantCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"schema_name":   "test-schema",
		"database_name": "test-db",
		"privilege":     "USAGE",
		"roles":         []string{"test-role-1", "test-role-2"},
		"shares":        []string{"test-share-1", "test-share-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.SchemaGrant().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT USAGE ON SCHEMA "test-db"."test-schema" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON SCHEMA "test-db"."test-schema" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON SCHEMA "test-db"."test-schema" TO SHARE "test-share-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON SCHEMA "test-db"."test-schema" TO SHARE "test-share-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadSchemaGrant(mock)
		err := resources.CreateSchemaGrant(d, db)
		a.NoError(err)
	})
}

func expectReadSchemaGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "SCHEMA", "test-schema", "ROLE", "test-role-1", false, "bob").AddRow(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "VIEW", "test-view", "ROLE", "test-role-2", false, "bob").AddRow(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "VIEW", "test-view", "SHARE", "test-share-1", false, "bob").AddRow(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "VIEW", "test-view", "SHARE", "test-share-2", false, "bob")
	mock.ExpectQuery(`^SHOW GRANTS ON SCHEMA "test-db"."test-schema"$`).WillReturnRows(rows)
}
