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

	enableAudit, err := config.GF.Get(config.EnableAuditGfId)
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
	err := config.Load()
	if err != nil {
		panic(err)
	}

	dbHostAny, err := config.GF.Get(config.DatabaseHostGfId)
	dbHost := dbHostAny.(string)
	if err != nil {
		panic(err)
	}
	dbPort, err := config.GF.GetString(config.DatabasePortGfId)
	if err != nil {
		panic(err)
	}
	dbUser, err := config.GF.GetString(config.DatabaseUserGfId)
	if err != nil {
		panic(err)
	}
	dbPassword, err := config.GF.GetString(config.DatabasePasswordGfId)
	if err != nil {
		panic(err)
	}
	dbName, err := config.GF.GetString(config.DatabaseNameGfId)
	if err != nil {
		panic(err)
	}

	setUpDbConn(dbHost, dbPort, dbUser, dbPassword, dbName)
	doSomething()

}
