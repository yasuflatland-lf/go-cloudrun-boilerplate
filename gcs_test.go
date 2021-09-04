package main

import (
	"context"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	_ "github.com/fsouza/fake-gcs-server/fakestorage"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGCSHandlerSmoke(t *testing.T) {

	runServersTest(t, nil, func(t *testing.T, server *fakestorage.Server) {
		t.Run("Write and IsExist Smoke Test", func(t *testing.T) {
			const content = "some nice content"
			const bucketName = "some-bucket"
			const objectName = "other/interesting/object.txt"

			server.CreateBucketWithOpts(fakestorage.CreateBucketOpts{Name: bucketName})
			client := server.Client()

			ctx := context.Background()
			gcs := NewGCS(ctx, client)

			_, err := gcs.Write(ctx, bucketName, objectName, []byte(content), ContentTypeText)
			assert.Nil(t, err)

			isExist := gcs.IsExist(ctx, bucketName, objectName)

			obj, err := server.GetObject(bucketName, objectName)
			assert.Nil(t, err)
			assert.Equal(t, true, isExist)
			assert.Equal(t, content, string(obj.Content))

		})
	})
}
