package entity

import "time"

type File struct {
	ID        int       `gorm:"primaryKey" json:"fileid"`
	Name      string    `gorm:"unique" json:"name"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *File) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ID":        f.ID,
		"Name":      f.Name,
		"Size":      f.Size,
		"MimeType":  f.MimeType,
		"CreatedAt": f.CreatedAt,
	}
}
