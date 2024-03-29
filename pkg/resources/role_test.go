package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viostream/terraform-provider-snowflake/pkg/provider"
	"github.com/viostream/terraform-provider-snowflake/pkg/resources"
	. "github.com/viostream/terraform-provider-snowflake/pkg/testhelpers"
)

func TestRole(t *testing.T) {
	r := require.New(t)
	err := resources.Role().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestRoleCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"name":    "good_name",
		"comment": "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Role().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE ROLE "good_name" COMMENT='great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadRole(mock)
		err := resources.CreateRole(d, db)
		a.NoError(err)
	})
}

func expectReadRole(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "is_default", "is_current", "is_inherited", "assigned_to_users", "granted_to_roles", "granted_roles", "owner", "comment",
	},
	).AddRow("created_on", "role name", "is_default", "is_current", "is_inherited", "assigned_to_users", "granted_to_roles", "granted_roles", "owner", "mock comment")
	mock.ExpectQuery(`SHOW ROLES LIKE 'good_name'`).WillReturnRows(rows)
}

func TestRoleRead(t *testing.T) {
	a := assert.New(t)

	d := role(t, "good_name", map[string]interface{}{"name": "good_name"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadRole(mock)
		err := resources.ReadRole(d, db)
		a.NoError(err)
		a.Equal("mock comment", d.Get("comment").(string))
		a.Equal("role name", d.Get("name").(string))
	})
}

func TestRoleDelete(t *testing.T) {
	a := assert.New(t)

	d := role(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP ROLE "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteRole(d, db)
		a.NoError(err)
	})
}
