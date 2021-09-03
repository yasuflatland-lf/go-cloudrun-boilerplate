package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"golang.org/x/xerrors"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type (
	Secret interface {
		GetSecretName(secretName string) string
		GetSecret(name string) (string, error)
		DeleteSecret(secretId string) error
		CreateSecret(secretId string, plainText interface{}) (*secretmanagerpb.Secret, error)
	}

	secret struct {
		ctx         context.Context
		ProjectUuid string
		ProjectId   string
	}
)

// Google APIs Secret Manager Library
// https://pkg.go.dev/cloud.google.com/go/secretmanager@v0.1.0/apiv1
func NewSecret(ctx context.Context, projectId string, projectUuid string) *secret {
	return &secret{
		ctx, projectId, projectUuid,
	}
}

// Get Secret from Google Secret Manager
func (s *secret) GetSecret(secretId string) (string, error) {

	// Get Client
	client, err := secretmanager.NewClient(s.ctx)
	if err != nil {
		return "", xerrors.Errorf("failed to fetch secret manager : %w\n", err)
	}
	defer client.Close()

	// Get Secret Path
	secretPath := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", s.ProjectId, secretId)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	}

	result, err := client.AccessSecretVersion(s.ctx, req)
	if err != nil {
		return "", xerrors.Errorf("failed to access secret version : %w : %w\n", secretId, err)
	}

	return string(result.Payload.Data), nil
}

func (s *secret) DeleteSecret(secretId string) error {
	// Get Client
	client, err := secretmanager.NewClient(s.ctx)
	if err != nil {
		return xerrors.Errorf("failed to fetch secret manager : %w\n", err)
	}
	defer client.Close()

	// Get Secret Path
	secretPath := fmt.Sprintf("projects/%s/secrets/%s", s.ProjectId, secretId)

	req := &secretmanagerpb.DeleteSecretRequest{
		Name: secretPath,
	}

	err = client.DeleteSecret(s.ctx, req)
	if err != nil {
		return xerrors.Errorf("failed to delete secret : %w : %w\n", secretId, err)
	}

	return nil
}

func (s *secret) CreateSecret(secretId string, plainText string) (*secretmanagerpb.Secret, error) {
	// Storing data
	bytes := []byte(plainText)

	// Get Client
	client, err := secretmanager.NewClient(s.ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to fetch secret manager : %w\n", err)
	}
	defer client.Close()

	// Create Secret
	secret, err := client.CreateSecret(s.ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", s.ProjectUuid),
		SecretId: secretId,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	})

	// When secret already exists
	if err != nil {
		return secret, xerrors.Errorf("Secret has already been created. : %w\n", err)
	}

	// Write data into the latest version.
	_, err = client.AddSecretVersion(s.ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: bytes,
		},
	})

	return secret, nil
}
