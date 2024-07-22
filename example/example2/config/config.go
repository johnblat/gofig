package config

import (
	"fmt"

	"github.com/ippontech/gofig"
)

var (
	// Config is the configuration for the application
	DatabaseHostGfId         gofig.Id
	DatabasePortGfId         gofig.Id
	DatabaseUserGfId         gofig.Id
	DatabasePasswordGfId     gofig.Id
	DatabaseNameGfId         gofig.Id
	EnableAuditGfId          gofig.Id
	EnableVerboseLoggingGfId gofig.Id
	EnvironmentGfId          gofig.Id
)

var GF gofig.Gofig

// Init
// Get-family functions

// auto-add to readme
// lookup id based on name given

func Load() error {
	var err error
	initOpts := []gofig.InitOpt{
		{
			Name:        "DATABASE_HOST",
			Description: "The database host",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &DatabaseHostGfId,
		},
		{
			Name:        "DATABASE_PORT",
			Description: "The database port.",
			Type:        gofig.TypeString,
			Required:    false,
			Default:     "5432",
			IdPtr:       &DatabasePortGfId,
		},
		{
			Name:        "DATABASE_USER",
			Description: "The username for the database",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &DatabaseUserGfId,
		},
		{
			Name:        "DATABASE_PASSWORD",
			Description: "The password for the database",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &DatabasePasswordGfId,
		},
		{
			Name:        "DATABASE_NAME",
			Description: "The name of the database",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &DatabaseNameGfId,
		},
		{
			Name:        "ENABLE_AUDIT",
			Description: "Enable audit logging",
			Type:        gofig.TypeBool,
			Required:    false,
			Default:     false,
			IdPtr:       &EnableAuditGfId,
		},
		{
			Name:        "ENABLE_VERBOSE_LOGGING",
			Description: "Enable verbose logging",
			Type:        gofig.TypeBool,
			Required:    false,
			Default:     false,
			IdPtr:       &EnableVerboseLoggingGfId,
		},
		{
			Name:        "ENVIRONMENT",
			Description: "The environment the application is running in",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &EnvironmentGfId,
		},
	}
	docStr, err := gofig.DocString(initOpts)
	if err != nil {
		return err
	}
	fmt.Println(docStr)

	GF, err = gofig.Init(initOpts)
	if err != nil {
		return err
	}

	return err
}
