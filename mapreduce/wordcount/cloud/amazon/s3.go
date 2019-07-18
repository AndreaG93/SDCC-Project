package amazon

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system"
	"SDCC-Project-WorkerNode/utility"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type S3Client struct {
	session *session.Session
}

func New() *S3Client {

	output := new(S3Client)

	(*output).session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(system.AmazonAWSRegion)},
	))

	return output
}

func (obj *S3Client) Upload(inputFilename string, key string) {

	var inputFile *os.File
	var err error

	amazonAWSS3Uploader := s3manager.NewUploader((*obj).session)

	inputFile, err = os.Open(inputFilename)
	utility.CheckError(err)

	defer func() {
		utility.CheckError(inputFile.Close())
	}()

	_, err = amazonAWSS3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(system.AmazonS3BucketName),
		Key:    aws.String(key),
		Body:   inputFile,
	})
	utility.CheckError(err)
}

func (obj *S3Client) Download(key string, outputFilename string) {

	var outputFile *os.File
	var err error

	amazonAWSS3Downloader := s3manager.NewDownloader((*obj).session)

	outputFile, err = os.Create(outputFilename)
	utility.CheckError(err)

	_, err = amazonAWSS3Downloader.Download(outputFile, &s3.GetObjectInput{
		Bucket: aws.String(system.AmazonS3BucketName),
		Key:    aws.String(key),
	})
	utility.CheckError(err)
}

func (obj *S3Client) Delete(key string) {

	var err error
	amazonS3Service := s3.New((*obj).session)

	_, err = amazonS3Service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(system.AmazonS3BucketName),
		Key:    aws.String(key),
	})
	utility.CheckError(err)

	err = amazonS3Service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(system.AmazonS3BucketName),
		Key:    aws.String(key),
	})
	utility.CheckError(err)
}
