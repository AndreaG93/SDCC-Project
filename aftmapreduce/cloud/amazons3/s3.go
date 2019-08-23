package amazons3

import (
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

const (
	AmazonAWSRegion    = "us-east-1"
	AmazonS3BucketName = "graziani-filestorage"
)

type S3Client struct {
	session *session.Session
}

func New() *S3Client {

	output := new(S3Client)

	(*output).session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(AmazonAWSRegion)},
	))

	return output
}

func (obj *S3Client) Upload(inputFile *os.File, key string) {

	var err error

	amazonAWSS3Uploader := s3manager.NewUploader((*obj).session)

	_, err = amazonAWSS3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AmazonS3BucketName),
		Key:    aws.String(key),
		Body:   inputFile,
	})
	utility.CheckError(err)
}

func (obj *S3Client) Download(key string, outputFile *os.File) {

	amazonAWSS3Downloader := s3manager.NewDownloader((*obj).session)

	_, err := amazonAWSS3Downloader.Download(outputFile, &s3.GetObjectInput{
		Bucket: aws.String(AmazonS3BucketName),
		Key:    aws.String(key),
	})
	utility.CheckError(err)
}

func (obj *S3Client) Delete(key string) {

	var err error
	amazonS3Service := s3.New((*obj).session)

	_, err = amazonS3Service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(AmazonS3BucketName),
		Key:    aws.String(key),
	})
	if err != nil {

		if err.Error() == s3.ErrCodeNoSuchKey {
			fmt.Println("Specified Key not found - OK")
		} else {
			utility.CheckError(err)
		}
	}
	err = amazonS3Service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(AmazonS3BucketName),
		Key:    aws.String(key),
	})
	utility.CheckError(err)
}
