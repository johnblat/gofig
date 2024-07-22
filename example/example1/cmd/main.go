package main

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

func sendAudit() {
	fmt.Println("Sent audit")
}

func doSomething() {
	fmt.Println("Did something")

	enableAudit, err := gf.Get(EnableAuditGfId)
	if err != nil {
		panic(err)
	}

	if enableAudit == true {
		sendAudit()
	}
}

func setUpDbConn(dbHost, dbPort, dbUser, dbPassword, dbName string) {
	fmt.Printf("Connecting to %s:%s as %s. Set up Db Conn\n", dbHost, dbPort, dbUser)
}

var gf gofig.Gofig

func main() {

	// setup gofig
	var err error
	gf, err = gofig.Init([]gofig.InitOpt{
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
	})
	if err != nil {
		panic(err)
	}

	// actual stuff
	dbHost, err := gf.GetString(DatabaseHostGfId)
	if err != nil {
		panic(err)
	}
	dbPort, err := gf.GetString(DatabasePortGfId)
	if err != nil {
		panic(err)
	}
	dbUser, err := gf.GetString(DatabaseUserGfId)
	if err != nil {
		panic(err)
	}
	dbPassword, err := gf.GetString(DatabasePasswordGfId)
	if err != nil {
		panic(err)
	}
	dbName, err := gf.GetString(DatabaseNameGfId)
	if err != nil {
		panic(err)
	}

	setUpDbConn(dbHost, dbPort, dbUser, dbPassword, dbName)
	doSomething()

}
