package main

import (
	"os"

	"github.com/ippontech/gofig/example/example3/config"
	cdv "github.com/ippontech/gofig/example/example3/config_derived_values"
)

func main() {

	os.Setenv("DATABASE_ENGINE", "postgres")
	os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("DATABASE_USER", "user")
	os.Setenv("DATABASE_PASSWORD", "password")
	os.Setenv("DATABASE_NAME", "dbname")
	os.Setenv("ENABLE_AUDIT", "true")
	os.Setenv("ENABLE_VERBOSE_LOGGING", "true")

	err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := cdv.DeriveDbConnFromConfig()
	if err != nil {
		panic(err)
	}

	httpClient, err := cdv.DeriveHttpClientFromConfig()
	if err != nil {
		panic(err)
	}

	db.Close()
	httpClient.Do(nil)

}
