package cloud_storage

type CloudStorage interface {
	DownloadFile() error
	UploadFile() float64
}
