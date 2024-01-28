package in

type MultipartFileDTO struct {
	FileContent []byte
	Filename    string
	Size        int64
	Host        string
	Path        string
	FileID      int64
}
