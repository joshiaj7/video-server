package response

import "time"

type File struct {
	ID        string    `json:"fileid"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
