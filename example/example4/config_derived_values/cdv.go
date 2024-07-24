package cdv

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ippontech/gofig/example/example3/config"
)

// demonstrate pattern of config derived values
// for instance a devClient, uatClient, prodClient
// that are all different in various ways
// for example, dev clients will pickup HTTP_PROXY, whereas prod client does not

func makeDevClient() *http.Client {
	// dev client
	return &http.Client{}
}

func makeUatClient() *http.Client {
	// uat client
	return &http.Client{}
}

func makeProdClient() *http.Client {
	// prod client
	return &http.Client{}
}

func makeLocalClient() *http.Client {
	// local client
	return &http.Client{}
}

func DeriveHttpClientFromConfig() (*http.Client, error) {
	env, err := config.GF.GetString(config.EnvironmentGfId)
	if err != nil {
		return nil, err
	}

	switch env {
	case "dev":
		return makeDevClient(), nil
	case "uat":
		return makeUatClient(), nil
	case "prod":
		return makeProdClient(), nil
	case "local":
		return makeLocalClient(), nil
	}

	return nil, fmt.Errorf("invalid environment. must be one of dev, uat, prod, local. got: %s", env)
}

func makePostgresDbConn() sql.Conn {
	// postgres db config
	return sql.Conn{}
}

func makeMysqlDbConn() sql.Conn {
	// mysql db config
	return sql.Conn{}
}

func makeSqliteDbConn() sql.Conn {
	// sqlite db config
	return sql.Conn{}
}

func DeriveDbConnFromConfig() (sql.Conn, error) {
	dbEngine, err := config.GF.GetString(config.DatabaseEngineGfId)
	if err != nil {
		return sql.Conn{}, err
	}

	switch dbEngine {
	case "postgres":
		return makePostgresDbConn(), nil
	case "mysql":
		return makeMysqlDbConn(), nil
	case "sqlite":
		return makeSqliteDbConn(), nil
	}

	return sql.Conn{}, fmt.Errorf("invalid db engine. must be one of postgres, mysql, sqlite. got: %s", dbEngine)
}
