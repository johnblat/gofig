package config

import "github.com/ippontech/gofig"

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

var App gofig.Gofig

func Load() error {
	err := App.Init(
		gofig.InitOpt{
			Name:        "DATABASE_HOST",
			Description: "The database host",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &DatabaseHostGfId,
		},
		gofig.InitOpt{
			Name:        "DATABASE_PORT",
			Description: "The database port.",
			Type:        gofig.TypeString,
			Required:    false,
			Default:     "5432",
			IdPtr:       &DatabasePortGfId,
		},
		gofig.InitOpt{
			Name:        "DATABASE_USER",
			Description: "The username for the database",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &DatabaseUserGfId,
		},
		gofig.InitOpt{
			Name:        "DATABASE_PASSWORD",
			Description: "The password for the database",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &DatabasePasswordGfId,
		},
		gofig.InitOpt{
			Name:        "DATABASE_NAME",
			Description: "The name of the database",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &DatabaseNameGfId,
		},
		gofig.InitOpt{
			Name:        "ENABLE_AUDIT",
			Description: "Enable audit logging",
			Type:        gofig.TypeBool,
			Required:    false,
			Default:     false,
			IdPtr:       &EnableAuditGfId,
		},
		gofig.InitOpt{
			Name:        "ENABLE_VERBOSE_LOGGING",
			Description: "Enable verbose logging",
			Type:        gofig.TypeBool,
			Required:    false,
			Default:     false,
			IdPtr:       &EnableVerboseLoggingGfId,
		},
		gofig.InitOpt{
			Name:        "ENVIRONMENT",
			Description: "The environment the application is running in",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &EnvironmentGfId,
		},
	)

	return err
}
