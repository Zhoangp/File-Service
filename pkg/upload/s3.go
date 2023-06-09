package upload

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Zhoangp/File-Service/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type UploadProvider interface {
	UploadFile(ctx context.Context, data []byte, dst string) (string, error)
	DeleteFile(objectUrl string) error
}

type s3Provider struct {
	bucket  string
	region  string
	apiKey  string
	secret  string
	domain  string
	session *session.Session
}

func NewS3Provider(cfg *config.Config) *s3Provider {
	provider := &s3Provider{
		bucket: cfg.AWS.S3Bucket,
		region: cfg.AWS.Region,
		apiKey: cfg.AWS.APIKey,
		secret: cfg.AWS.SecretKey,
		domain: cfg.AWS.S3Domain,
	}

	s3Session, err := session.NewSession(&aws.Config{
		Region: aws.String(provider.region),
		Credentials: credentials.NewStaticCredentials(
			provider.apiKey, // Access key ID
			provider.secret, // Secret access key
			"",              // Token có thể bỏ qua
		),
	})

	if err != nil {

		log.Fatalln(err)
	}

	provider.session = s3Session

	return provider
}
func (p *s3Provider) DeleteFile(objectUrl string) error {
	u, err := url.Parse(objectUrl)
	if err != nil {
		return err
	}
	objectKey := strings.TrimPrefix(u.Path, "/")
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(objectKey),
	}

	_, err = s3.New(p.session).DeleteObject(input)
	if err != nil {
		return err
	}
	return nil
}
func (p *s3Provider) UploadFile(ctx context.Context, data []byte, dst string) (string, error) {
	fileBytes := bytes.NewReader(data)
	fileType := http.DetectContentType(data)
	_, err := s3.New(p.session).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(dst),
		ACL:         aws.String("private"),
		ContentType: aws.String(fileType),
		Body:        fileBytes,
	})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return fmt.Sprintf("%s/%s", p.domain, dst), nil
}
