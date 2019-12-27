package rpc

import (
	"context"
	//"fmt"
	"io"
	"io/ioutil"
	"os"
	//"mime/multipart"
	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type Uploader struct {
	logger *logrus.Entry
	bucket string
	gscli  *storage.Client
	bh     *storage.BucketHandle
}

func NewUploader(bucket string, logger *logrus.Entry) (*Uploader, error) {
	gscli, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	bh := gscli.Bucket(bucket)
	_, err = bh.Attrs(context.Background())
	if err != nil {
		return nil, err
	}

	return &Uploader{
		logger: logger,
		gscli:  gscli,
		bh:     bh,
		bucket: bucket,
	}, nil
}

func (u *Uploader) uploadSegments(streamID string, localMediaDir string) error {

	logger := u.logger.WithFields(logrus.Fields{
		"stream_id": streamID,
		"bucket":    u.bucket,
	})

	logger.Info("uploading segments")

	files, err := ioutil.ReadDir(localMediaDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		emptyCtx := context.Background()
		filePath := filepath.Join(localMediaDir, f.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		obj := u.bh.Object(filePath)
		w := obj.NewWriter(emptyCtx)
		w.CacheControl = "no-cache"
		//w.ContentType = ct  TODO:!

		if _, err := io.Copy(w, file); err != nil {
			return err
		}

		if err := w.Close(); err != nil {
			return err
		}

		if err := obj.ACL().Set(emptyCtx, storage.AllUsers, storage.RoleReader); err != nil {
			return err
		}

		_, err = obj.Attrs(emptyCtx)
		if err != nil {
			return err
		}

		logger.Info("segment has been uploaded successfully")
	}

	return err
}
