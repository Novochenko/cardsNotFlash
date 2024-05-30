package store

import "errors"

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrNicknameNotAvialable = errors.New("nickname is not avialable")
	ErrEmptyLK              = errors.New("lk is empty")
)
