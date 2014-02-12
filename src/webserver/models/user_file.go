package models

type UserFile struct {
	Title  string
	Path   string
	UserId int64
	Size   int64
	Tags   []string
}
