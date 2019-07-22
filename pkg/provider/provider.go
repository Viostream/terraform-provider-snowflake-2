package provider

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/snowflakedb/gosnowflake"
	"github.com/viostream/terraform-provider-snowflake/pkg/db"
	"github.com/viostream/terraform-provider-snowflake/pkg/resources"
)

// Provider is a provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_PASSWORD", nil),
				Sensitive:     true,
				ConflictsWith: []string{"browser_auth"},
			},
			"browser_auth": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("SNOWFLAKE_USE_BROWSER_AUTH", nil),
				Sensitive:     false,
				ConflictsWith: []string{"password"},
			},
			"role": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_ROLE", nil),
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SNOWFLAKE_REGION", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"snowflake_database":         resources.Database(),
			"snowflake_database_grant":   resources.DatabaseGrant(),
			"snowflake_managed_account":  resources.ManagedAccount(),
			"snowflake_resource_monitor": resources.ResourceMonitor(),
			"snowflake_role":             resources.Role(),
			"snowflake_role_grants":      resources.RoleGrants(),
			"snowflake_schema":           resources.Schema(),
			"snowflake_schema_grant":     resources.SchemaGrant(),
			"snowflake_share":            resources.Share(),
			"snowflake_user":             resources.User(),
			"snowflake_view":             resources.View(),
			"snowflake_view_grant":       resources.ViewGrant(),
			"snowflake_warehouse":        resources.Warehouse(),
			"snowflake_warehouse_grant":  resources.WarehouseGrant(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  configureProvider,
	}
}

// configureProvider sets up the connection
func configureProvider(s *schema.ResourceData) (interface{}, error) {
	dsn, err := DSN(s)
	if err != nil {
		return nil, errors.Wrap(err, "could not build dsn for snowflake connection")
	}

	log.Printf("[DEBUG] connecting to %s", dsn)
	db, err := db.Open(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open snowflake database.")
	}

	return db, nil
}

// DSN returns the connection string for Snowflake
func DSN(s *schema.ResourceData) (string, error) {
	account := s.Get("account").(string)
	username := s.Get("username").(string)
	password := s.Get("password").(string)
	browserAuth := s.Get("browser_auth").(bool)
	region := s.Get("region").(string)
	role := s.Get("role").(string)

	var auth gosnowflake.AuthType

	if browserAuth {
		auth = gosnowflake.AuthTypeExternalBrowser
	}

	return gosnowflake.DSN(&gosnowflake.Config{
		Account:       account,
		User:          username,
		Region:        region,
		Password:      password,
		Role:          role,
		Authenticator: auth,
	})
}
