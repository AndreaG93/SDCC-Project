package amazons3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"time"
)

const (
	awsRegion            = "us-east-1"
	bucketName           = "graziani-filestorage"
	preSignedURLValidity = 15 * time.Minute
)

type Client struct {
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func New() *Client {

	output := new(Client)

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	))

	(*output).client = s3.New(awsSession)
	(*output).uploader = s3manager.NewUploader(awsSession)
	(*output).downloader = s3manager.NewDownloader(awsSession)

	return output
}

func (obj *Client) Put(key string, value []byte) error {

	_, err := (*(*obj).uploader).Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(value),
	})

	return err
}

func (obj *Client) Get(key string) ([]byte, error) {

	output := aws.NewWriteAtBuffer([]byte{})

	_, err := (*(*obj).downloader).Download(output, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	return output.Bytes(), err
}

func (obj *Client) RetrieveURLForPutOperation(key string) (string, error) {

	req, _ := (*(*obj).client).PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	return req.Presign(preSignedURLValidity)
}

func (obj *Client) RetrieveURLForGetOperation(key string) (string, error) {

	req, _ := (*(*obj).client).GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	return req.Presign(preSignedURLValidity)
}

func (obj *Client) Remove(key string) error {

	var err error

	_, err = (*(*obj).client).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil && err.Error() == s3.ErrCodeNoSuchKey {
		return nil
	} else {
		return (*(*obj).client).WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
	}
}
