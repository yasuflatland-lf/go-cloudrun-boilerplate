package main

import (
	"cloud.google.com/go/storage"
	"context"
)

type (
	Repository interface {
		CloudSQL() CloudSQL
		GCS() GCS
	}
	repository struct {
		cloudSQL CloudSQL
		gcs      GCS
	}
)

func NewRepository(ctx context.Context, client *storage.Client) Repository {
	return &repository{
		cloudSQL: NewCloudSQL(ctx),
		gcs:      NewGCS(ctx, client),
	}
}

func (r *repository) CloudSQL() CloudSQL {
	return r.cloudSQL
}

func (r *repository) GCS() GCS {
	return r.gcs
}
