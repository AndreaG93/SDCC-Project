package s3

import (
	"SDCC-Project-WorkerNode/utility"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type Client struct {
	session *session.Session
}

func New() *Client {

	output := new(Client)

	(*output).session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(utility.AmazonAWSRegion)},
	))

	return output
}

func (obj *Client) Upload(inputFilename string, outputFileName string) error {

	var inputFile *os.File
	var err error

	amazonAWSS3Uploader := s3manager.NewUploader((*obj).session)

	if inputFile, err = os.Open(inputFilename); err != nil {
		return err
	}
	defer func() {
		utility.CheckError(inputFile.Close())
	}()

	if _, err = amazonAWSS3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(utility.AmazonS3BucketName),
		Key:    aws.String(outputFileName),
		Body:   inputFile,
	}); err != nil {
		return err
	}
	return nil
}

func (obj *Client) Download(filename string, outputFilename string) error {

	var outputFile *os.File
	var err error

	amazonAWSS3Downloader := s3manager.NewDownloader((*obj).session)

	if outputFile, err = os.Create(outputFilename); err != nil {
		return err
	}

	if _, err = amazonAWSS3Downloader.Download(outputFile, &s3.GetObjectInput{
		Bucket: aws.String(utility.AmazonS3BucketName),
		Key:    aws.String(filename),
	}); err != nil {
		return err
	}

	return nil
}

func (obj *Client) Delete(objectKey string) error {

	var err error
	amazonS3Service := s3.New((*obj).session)

	if _, err = amazonS3Service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(utility.AmazonS3BucketName),
		Key:    aws.String(objectKey),
	}); err != nil {
		return err
	}

	if err = amazonS3Service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(utility.AmazonS3BucketName),
		Key:    aws.String(objectKey),
	}); err != nil {
		return err
	}

	return nil
}
