package main

import (
	"context"
	"github.com/glassonion1/logz"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
	"sync"
)

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
	}

	applicationConfig struct {
		// Environment
		Environment string `required:"true" envconfig:"APP_ENV" default:"test"`
		ProjectUuid string `required:"true" envconfig:"PROJECT_UUID" default:""`
		ProjectId   string `required:"true" envconfig:"PROJECT_ID" default:""`
		ImageName   string `required:"true" envconfig:"IMAGE_NAME" default:"go-cloudrun-boilerplate"`
		Port        int    `required:"false" envconfig:"PORT" default:1323`

		// Database
		DBIP         string `required:"false" envconfig:"DB_IP" default:"127.0.0.1"`
		DBPort       int64  `required:"true" envconfig:"DB_PORT" default:"3306"`
		MaxIdleConns int    `required:"false" envconfig:"DB_MAX_IDLE_CONNS" default:"10"`
		MaxOpenConns int    `required:"false" envconfig:"DB_MAX_OPEN_CONNS" default:"100"`

		// GCS
		BucketName string `required:"false" envconfig:"BUCKET_NAME" default:"go-cloudrun-boilerplate-us-central1-data"`
		ObjectName string `required:"false" envconfig:"OBJECT_NAME" default:"test.json"`

		// Instance related
		TimeOut int `required:"false" envconfig:"TIMEOUT" default:"1200"`

		// Secrets
		UserName         string `required:"true" envconfig:"DB_USERNAME" default:"root"`
		Password         string `required:"true" envconfig:"DB_PASSWORD" default:"admin"`
		CloudSQLInstance string `required:"false" envconfig:"CLOUDSQL_INSTANCES" default:""`
		Name             string `required:"true" envconfig:"DB_NAME" default:"test"`
	}
)

var (
	appConfig *applicationConfig
	once      sync.Once
)

func GetApplicationConfig(ctx context.Context) *applicationConfig {
	// Google APIs Secret Manager Library
	// https://pkg.go.dev/cloud.google.com/go/secretmanager@v0.1.0/apiv1
	once.Do(func() {
		appConfig = &applicationConfig{}

		// Bind Environment valuable to the structure automatically.
		// The default values are used when no environment valuables are defined
		err := envconfig.Process("", appConfig)
		if err != nil {
			logz.Criticalf(ctx, "Required environment values are not defined properly. Please check required values. : %+v", err)
		}

		// Only for Production
		if appConfig.IsProduction() {
			// Create instance
			secret := NewSecret(ctx, appConfig.ProjectId, appConfig.ProjectUuid)

			// Cloud SQL Instance
			appConfig.CloudSQLInstance, err = secret.GetSecret(appConfig.ImageName + "-CLOUDSQL_INSTANCES")
			if err != nil {
				logz.Criticalf(ctx, "%+v\n", xerrors.Errorf("CLOUDSQL_INSTANCES : %w\n", err))
			}

			appConfig.Name, err = secret.GetSecret(appConfig.ImageName + "-DB_NAME")
			if err != nil {
				logz.Criticalf(ctx, "%+v\n", xerrors.Errorf("DB_NAME : %w\n", err))
			}

			appConfig.UserName, err = secret.GetSecret(appConfig.ImageName + "-DB_USERNAME")
			if err != nil {
				logz.Criticalf(ctx, "%+v\n", xerrors.Errorf("DB_USERNAME : %w\n", err))
			}

			appConfig.Password, err = secret.GetSecret(appConfig.ImageName + "-DB_PASSWORD")
			if err != nil {
				logz.Criticalf(ctx, "%+v\n", xerrors.Errorf("DB_PASSWORD : %w\n", err))
			}
		}
	})
	return appConfig
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
