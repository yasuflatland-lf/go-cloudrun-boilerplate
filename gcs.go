package main

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/glassonion1/logz"
	"golang.org/x/xerrors"
	"io/ioutil"
)

type (
	GCS interface {
		Object(bucketName string, objectName string) *storage.ObjectHandle
		Write(ctx context.Context, bucketName string, objectName string, writeStr []byte, contentType string) (*storage.Writer, error)
		Read(ctx context.Context, bucketName string, objectName string) ([]byte, error)
		IsExist(ctx context.Context, bucketName string, objectName string) bool
	}

	gcs struct {
		credentialFilePath string
		storageClient      *storage.Client
	}
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

func NewGCS(ctx context.Context, client *storage.Client) GCS {
	g := &gcs{}

	if client == nil {
		// Production should be passed client is null, then create the new client
		newClient, err := storage.NewClient(ctx)
		if err != nil {
			logz.Criticalf(ctx, "%+v\n", xerrors.Errorf(": %w", err))
		}
		g.storageClient = newClient
	} else {
		// For testing, test server client should be passed here.
		g.storageClient = client
	}

	return g
}

// GetDomains Object
func (g *gcs) Object(bucketName string, objectName string) *storage.ObjectHandle {
	return g.storageClient.Bucket(bucketName).Object(objectName)
}

func (g *gcs) Write(ctx context.Context, bucketName string, objectName string, writeStr []byte, contentType string) (*storage.Writer, error) {
	writer := g.Object(bucketName, objectName).NewWriter(ctx)
	writer.ContentType = contentType

	if _, err := writer.Write(writeStr); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return writer, nil
}

func (g *gcs) Read(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	reader, err := g.Object(bucketName, objectName).NewReader(ctx)

	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, xerrors.Errorf("ioutil.ReadAll: %+v\n", err)
	}

	return data, nil
}

func (g *gcs) IsExist(ctx context.Context, bucketName string, objectName string) bool {
	reader, err := g.Object(bucketName, objectName).NewReader(ctx)
	if err != nil {
		logz.Infof(ctx, "No corresponding bucketName: %+v\n", xerrors.Errorf(": %w", err))
		return false
	}

	if err := reader.Close(); err != nil {
		logz.Errorf(ctx, "%+v\n", xerrors.Errorf(": %w", err))
		return false
	}

	return true
}
