package filesystem

import (
	"bytes"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/minio/minio-go/v6"
	"os"
	"path/filepath"
)

type MinioDriver struct {
	client *minio.Client
	bucket string
}

func (d MinioDriver) Exists(path string) bool {
	_, err := d.client.StatObject(d.bucket, path, minio.StatObjectOptions{})
	if err != nil {
		return false
	}
	return true
}

func (d MinioDriver) NotExists(path string) bool {
	return !d.Exists(path)
}

func (d MinioDriver) Put(path, content string) error {
	reader := bytes.NewReader([]byte(content))
	ext := filepath.Ext(path)
	mime := filetype.GetType(ext).MIME.Type
	_, err := d.client.PutObject(d.bucket, path, reader, reader.Size(), minio.PutObjectOptions{ContentType: mime})
	if err != nil {
		return err
	}
	return nil
}

func (d MinioDriver) Get(path string) (string, error) {
	reader, err := d.client.GetObject(d.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}

	defer func() {
		err := reader.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(reader); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (d MinioDriver) MakeDir(path string, perm os.FileMode) error {
	reader := bytes.NewReader(make([]byte, 0))
	_, err := d.client.PutObject(d.bucket, path+"/", reader, reader.Size(), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func NewMinioDriver() (*MinioDriver, error) {
	endpoint := "minio:9000"
	accessKeyID := "MqLknlISEZ1wPOp6Ra1T"
	secretAccessKey := "m3HHJSLgmRkbv186exPqQC4mHu5dUWnmgPIS23zc"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)

	if err != nil {
		return nil, err
	}
	doneCh := make(chan struct{})
	defer close(doneCh)
	for object := range minioClient.ListObjects("template", "ahhi/", true, doneCh) {
		if object.Err != nil {
			fmt.Println(object.Err)
			return nil, nil
		}
		fmt.Println(object.Key)
	}

	return &MinioDriver{
		client: minioClient,
		bucket: "template",
	}, nil
}
