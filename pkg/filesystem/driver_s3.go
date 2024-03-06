package filesystem

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kurneo/go-template/pkg/filesystem/helper"
	"github.com/kurneo/go-template/pkg/logger"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type S3Driver struct {
	s  *s3.S3
	l  logger.Contract
	b  string
	r  string
	p  helper.S3PathHelper
	pf helper.PathPreFixer
}

func (s S3Driver) FileExists(path string) (bool, error) {
	_, err := s.s.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.b),
		Key:    aws.String(s.pf.PrefixPath(path)),
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "NotFound") {
			return false, nil
		}
		s.l.Error(err)
		return false, err
	}

	return true, nil
}

func (s S3Driver) DirExists(path string) (bool, error) {
	_, err := s.s.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.b),
		Key:    aws.String(s.pf.PrefixDirectoryPath(path)),
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "NotFound") {
			return false, nil
		}
		s.l.Error(err)
		return false, err
	}

	return true, nil
}

func (s S3Driver) Put(path string, content []byte) error {
	_, err := s.s.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.b),
		ACL:           aws.String("public-read"),
		Key:           aws.String(s.pf.PrefixPath(path)),
		Body:          bytes.NewReader(content),
		ContentType:   aws.String(http.DetectContentType(content)),
		ContentLength: aws.Int64(int64(len(content))),
	})

	if err != nil {
		return err
	}
	return nil
}

func (s S3Driver) Get(path string) (string, error) {
	r, err := s.s.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.b),
		Key:    aws.String(s.pf.PrefixPath(path)),
	})

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(r.Body); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s S3Driver) MakeDir(path string, perm os.FileMode) error {
	acl := "private"
	if perm == 777 {
		acl = "public-read"
	}
	_, err := s.s.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.b),
		Key:         aws.String(s.pf.PrefixDirectoryPath(path)),
		ACL:         aws.String(acl),
		ContentType: aws.String("application/x-directory; charset=UTF-8"),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s S3Driver) Delete(path string) error {
	_, err := s.s.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.b),
		Key:    aws.String(s.pf.PrefixPath(path)),
	})

	if err != nil {
		return err
	}
	return nil
}

func (s S3Driver) ListContents(path string) ([]File, []Directory, error) {
	fPath := s.pf.PrefixDirectoryPath(s.p.GetDirectoryPath(path))
	resp, err := s.s.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(s.b),
		Prefix:    aws.String(fPath),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		return nil, nil, err
	}

	files := make([]File, 0)
	directories := make([]Directory, 0)

	for _, v := range resp.CommonPrefixes {
		switch true {
		case *v.Prefix == fPath:
			break
		default:
			directories = append(directories, Directory{
				Path:    s.pf.StripPrefix(s.pf.StripTrailingSeparator(*v.Prefix)),
				Name:    filepath.Base(*v.Prefix),
				ModTime: nil,
			})
		}
	}

	for _, v := range resp.Contents {
		switch true {
		case *v.Key == fPath:
			break
		case strings.HasSuffix(*v.Key, "/"):
			directories = append(directories, Directory{
				Path:    s.pf.StripPrefix(s.pf.StripTrailingSeparator(*v.Key)),
				Name:    filepath.Base(*v.Key),
				ModTime: v.LastModified,
			})
			break
		default:
			e := filepath.Ext(*v.Key)
			files = append(files, File{
				Path:      s.pf.StripPrefix(*v.Key),
				Name:      filepath.Base(*v.Key),
				ModTime:   v.LastModified,
				Size:      v.Size,
				Mime:      nil,
				Extension: &e,
			})
		}
	}

	return files, directories, nil
}

func (s S3Driver) Move(from, to string) error {
	err := s.Copy(from, to)

	if err != nil {
		return err
	}

	err = s.Delete(from)

	if err != nil {
		return err
	}

	return nil
}

func (s S3Driver) Copy(from, to string) error {

	_, err := s.s.CopyObject(&s3.CopyObjectInput{
		Bucket:            aws.String(s.b),
		Key:               aws.String(s.pf.PrefixPath(to)),
		ACL:               aws.String("public-read"),
		MetadataDirective: aws.String("COPY"),
		CopySource:        aws.String(fmt.Sprintf("%s/%s", s.b, s.pf.PrefixPath(from))),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s S3Driver) Rename(from, to string) error {
	return s.Move(from, to)
}

func (s S3Driver) Mime(path string) (string, string, error) {
	resp, err := s.s.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.b),
		Key:    aws.String(s.pf.PrefixPath(path)),
	})

	if err != nil {
		return "", "", err
	}

	return *resp.ContentType, filepath.Ext(path), nil
}

func (s S3Driver) ReadFile(path string) ([]byte, error) {
	r, err := s.s.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.b),
		Key:    aws.String(s.pf.PrefixPath(path)),
	})

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(r.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s S3Driver) Path(path string) string {
	return s.pf.PrefixPath(path)
}

func (s S3Driver) DirPath(path string) string {
	return s.pf.PrefixDirectoryPath(path)
}

func (s S3Driver) IsDir(path string) (bool, error) {
	resp, err := s.s.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.b),
		Key:    aws.String(s.pf.PrefixPath(path)),
	})

	if err != nil {
		return false, err
	}

	ct := strings.Split(*resp.ContentType, ";")

	if len(ct) == 0 {
		return false, nil
	}

	return ct[0] == "application/x-directory", nil
}

func (s S3Driver) PutObject(path string, f *bytes.Reader) error {
	_, err := s.s.PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(s.b),
		Key:                aws.String(s.pf.PrefixPath(path)),
		ACL:                aws.String("public-read"),
		Body:               f,
		ContentDisposition: aws.String("attachment"),
		// ContentLength:      aws.Int64(int64(len(buffer))),
		// ContentType:        aws.String(http.DetectContentType(buffer)),
	})

	if err != nil {
		return err
	}

	return nil
}

func NewS3Driver(r, b string, l logger.Contract, p helper.S3PathHelper, pf helper.PathPreFixer) *S3Driver {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(r),
	}))

	svc := s3.New(sess)

	return &S3Driver{
		s:  svc,
		l:  l,
		b:  b,
		r:  r,
		p:  p,
		pf: pf,
	}
}
