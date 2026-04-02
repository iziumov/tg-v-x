package dto

// const (
// 	StatusPending     = "pending"
// 	StatusDownloading = "downloading"
// 	StatusDone        = "done"
// 	StatusError       = "failed"
// )

type CreateJob struct {
	TgId int64
	Url  string
	// Status   string
	Platform string
}

type UpdateJob struct {
	ID       int
	Status   string
	ErrorMsg *string
}
