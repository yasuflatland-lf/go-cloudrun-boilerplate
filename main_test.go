package main

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
)

// Common Test Setting
func TestMain(m *testing.M) {

	// Place all MySQL related tests as sum test of this parent test
	// so that only one Instance up and test against it
	_, mysqlTerm := initMySQLContainer()
	defer mysqlTerm()

	// Run tests
	os.Exit(m.Run())
}

// Generate dsn and close function for test use.
// *** DO NOT USE FOR PRODUCTION ***
func initMySQLContainer() (string, func()) {
	ctx := context.Background()
	username := "root"
	password := "password"
	seedDataPath, err := os.Getwd()
	mysqlPort, _ := nat.NewPort("tcp", "3306")
	req := testcontainers.ContainerRequest{
		Image:        "mysql:5.7",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
		},
		BindMounts: map[string]string{
			seedDataPath + "/test/db/mysql_init": "/docker-entrypoint-initdb.d",
			seedDataPath + "/test/db/my.cnf":     "/etc/mysql/conf.d/my.cnf",
		},
		WaitingFor: wait.ForListeningPort(mysqlPort),
	}
	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	ip, err := mysqlC.Host(ctx)
	if err != nil {
		panic(err)
	}

	port, err := mysqlC.MappedPort(ctx, "3306")
	if err != nil {
		panic(err)
	}

	// cloudSQL service fetch MySQL connection data from environment valuables.
	// Set here dummy server information for test purpose.
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_USERNAME", username)
	os.Setenv("DB_PASSWORD", password)
	os.Setenv("DB_IP", ip)
	os.Setenv("DB_PORT", port.Port())

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/test", username, password, ip, port.Int())
	fmt.Println(dataSourceName)
	cTerm := func() {
		defer mysqlC.Terminate(ctx)
	}

	return dataSourceName, cTerm
}

// Test gcs Server wrapper
func runServersTest(t *testing.T, objs []fakestorage.Object, fn func(*testing.T, *fakestorage.Server)) {

	t.Run("tcp listener", func(t *testing.T) {
		t.Parallel()
		tcpServer, err := fakestorage.NewServerWithOptions(fakestorage.Options{NoListener: false, InitialObjects: objs})
		if err != nil {
			t.Fatal(err)
		}
		defer tcpServer.Stop()
		fn(t, tcpServer)
	})
	t.Run("no listener", func(t *testing.T) {
		t.Parallel()
		noListenerServer, err := fakestorage.NewServerWithOptions(fakestorage.Options{NoListener: true, InitialObjects: objs})
		if err != nil {
			t.Fatal(err)
		}
		fn(t, noListenerServer)
	})
}
