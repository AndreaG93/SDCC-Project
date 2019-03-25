package utility

const (
	DefaultNetwork                      = "tcp"
	AmazonAWSRegion                     = "us-east-1"
	AmazonS3BucketName                  = "graziani-filestorage"
	DefaultArbitraryFaultToleranceLevel = 3
)

func GetNumberOfActiveWorkerNodes() uint {
	return 3
}
