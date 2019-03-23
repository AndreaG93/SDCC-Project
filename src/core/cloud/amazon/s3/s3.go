package s3

import (
	"SDCC-Project-WorkerNode/src/core/utility"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type Client struct {
	session    *session.Session
	bucketName string
	region     string
}

func New(region string, bucketName string) *Client {

	output := new(Client)

	(*output).session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION)},
	))
	(*output).region = region
	(*output).bucketName = bucketName

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
		Bucket: aws.String((*obj).bucketName),
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
		Bucket: aws.String((*obj).bucketName),
		Key:    aws.String(filename),
	}); err != nil {
		return err
	}

	return nil
}

// ******************************************************

const (
	BUCKET_NAME = "graziani-filestorage"
	REGION      = "us-east-1"
)

func UploadFileToCloudStorage(inputFilename string, objectKey string) error {

	var file *os.File
	var err error

	amazonAWSSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION)},
	))
	amazonAWSS3Uploader := s3manager.NewUploader(amazonAWSSession)

	if file, err = os.Open(inputFilename); err != nil {
		return err
	}
	defer func() {
		utility.CheckError(file.Close())
	}()

	if _, err = amazonAWSS3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(objectKey),
		Body:   file,
	}); err != nil {
		return err
	}
	return nil
}

func DownloadFileFromCloudStorage(outputFilename string, objectKey string) error {

	var outputFile *os.File
	var err error

	amazonAWSSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION)},
	))
	amazonAWSS3Downloader := s3manager.NewDownloader(amazonAWSSession)

	if outputFile, err = os.Create(outputFilename); err != nil {
		return err
	}

	if _, err = amazonAWSS3Downloader.Download(outputFile, &s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(objectKey),
	}); err != nil {
		return err
	}

	return nil
}

func DeleteFileFromCloudStorage(objectKey string) error {

	var err error

	amazonAWSSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION)},
	))
	amazonS3Service := s3.New(amazonAWSSession)

	if _, err = amazonS3Service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(objectKey),
	}); err != nil {
		return err
	}

	if err = amazonS3Service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(objectKey),
	}); err != nil {
		return err
	}

	return nil
}
