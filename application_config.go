package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/glassonion1/logz"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"sync"
)

var once sync.Once

const (
	ENV_TEST        = "test"
	ENV_DEVELOPMENT = "development"
	ENV_PRODUCTION  = "production"
	APP_ENV         = "APP_ENV"
)

type (
	ApplicationConfig interface {
		IsTest() bool
		IsDevelopment() bool
		IsProduction() bool
		GetSecret(ctx context.Context, name string) ([]byte, error)
	}

	applicationConfig struct {
		// Environment
		Environment string `required:"true" envconfig:"APP_ENV" default:"test"`

		// Instance related
		TimeOut     int    `required:"false" envconfig:"TIMEOUT" default:"1200"`
		ImageName   string `required:"false" envconfig:"IMAGE_NAME" default:"go-cloudrun-boilerplate"`
		ProjectUuid string `required:"false" envconfig:"PROJECT_UUID" default:""`

		// Secrets
		UserName         string `required:"true" envconfig:"DB_USERNAME" default:"root"`
		Password         string `required:"true" envconfig:"DB_PASSWORD" default:"admin"`
		IP               string `required:"false" envconfig:"DB_IP" default:"127.0.0.1"`
		Port             int64  `required:"true" envconfig:"DB_PORT" default:"3306"`
		Name             string `required:"true" envconfig:"DB_NAME" default:"test"`
		MaxIdleConns     int    `required:"false" envconfig:"DB_MAX_IDLE_CONNS" default:"10"`
		MaxOpenConns     int    `required:"false" envconfig:"DB_MAX_OPEN_CONNS" default:"100"`
		LoadDataLimit    int64  `required:"false" envconfig:"LOAD_DATA_LIMIT" default:"1000"`
		BucketName       string `required:"false" envconfig:"BUCKET_NAME" default:"go-cloudrun-boilerplate-us-central1-data"`
		ObjectName       string `required:"false" envconfig:"OBJECT_NAME" default:"test.json"`
		CloudSQLInstance string `required:"false" envconfig:"CLOUDSQL_INSTANCES" default:""`
	}
)

func (conf *applicationConfig) GetInstance(ctx context.Context) *applicationConfig {
	// Google APIs Secret Manager Library
	// https://pkg.go.dev/cloud.google.com/go/secretmanager@v0.1.0/apiv1
	once.Do(func() {
		// Bind Environment valuable to the structure automatically.
		// The default values are used when no environment valuables are defined
		err := envconfig.Process("", &conf)
		if err != nil {
			logz.Criticalf(ctx, "Required environment values are not defined properly. Please check required values. : %+v", err)
		}

		if conf.IsProduction() {
			// Cloud SQL Instance
			conf.CloudSQLInstance, err = conf.GetSecret(ctx, conf.ImageName+"-CLOUDSQL_INSTANCES")
			if err != nil {
				logz.Criticalf(ctx, "%+v\n", xerrors.Errorf("CLOUDSQL_INSTANCES : %w\n", err))
			}
			conf.Name, err = conf.GetSecret(ctx, conf.ImageName+"-DB_NAME")
			if err != nil {
				logz.Criticalf(ctx, "%+v\n", xerrors.Errorf("DB_NAME : %w\n", err))
			}
			conf.UserName, err = conf.GetSecret(ctx, conf.ImageName+"-DB_USERNAME")
			if err != nil {
				logz.Criticalf(ctx, "%+v\n", xerrors.Errorf("DB_USERNAME : %w\n", err))
			}
			conf.Password, err = conf.GetSecret(ctx, conf.ImageName+"-DB_PASSWORD")
			if err != nil {
				logz.Criticalf(ctx, "%+v\n", xerrors.Errorf("DB_PASSWORD : %w\n", err))
			}
		}
	})
	return conf
}

// Get Secret path for fetching Secret from Secret Manager
// For ProjectUuid, please see TF_VAR_PROJECT_UUID.
func (conf *applicationConfig) GetSecretName(projectUuid string, secretName string) string {
	return fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectUuid, secretName)
}

// Get Secret from Google Secret Manager
func (conf *applicationConfig) GetSecret(ctx context.Context, name string) (string, error) {

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", xerrors.Errorf("failed to create secret manager : %w\n", err)
	}

	// Get Secret Path
	secretPath := conf.GetSecretName(conf.ProjectUuid, name)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", xerrors.Errorf("failed to access secret version : %w : %w\n", name, err)
	}

	return string(result.Payload.Data), nil
}

// Test if the environment is test
func (conf *applicationConfig) IsTest() bool {
	if conf.Environment == ENV_TEST {
		return true
	}
	return false
}

// Test if the environment is development
func (conf *applicationConfig) IsDevelopment() bool {
	if conf.Environment == ENV_DEVELOPMENT {
		return true
	}
	return false
}

// Test if the environment is production
func (conf *applicationConfig) IsProduction() bool {
	if conf.Environment == ENV_PRODUCTION {
		return true
	}
	return false
}
