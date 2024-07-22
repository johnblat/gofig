package main

import (
	"fmt"

	"github.com/ippontech/gofig/example/example2/config"
)

func sendAudit() {
	fmt.Println("Sent audit")
}

func doSomething() {
	fmt.Println("Did something")

	enableAudit, err := config.App.Get(config.EnableAuditGfId)
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

func main() {
	config.Load()

	dbHost, err := config.App.GetString(config.DatabaseHostGfId)
	if err != nil {
		panic(err)
	}
	dbPort, err := config.App.GetString(config.DatabasePortGfId)
	if err != nil {
		panic(err)
	}
	dbUser, err := config.App.GetString(config.DatabaseUserGfId)
	if err != nil {
		panic(err)
	}
	dbPassword, err := config.App.GetString(config.DatabasePasswordGfId)
	if err != nil {
		panic(err)
	}
	dbName, err := config.App.GetString(config.DatabaseNameGfId)
	if err != nil {
		panic(err)
	}

	setUpDbConn(dbHost, dbPort, dbUser, dbPassword, dbName)
	doSomething()
}
