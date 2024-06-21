package store

import "errors"

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrNicknameNotAvialable = errors.New("nickname is not avialable")
	ErrEmptyLK              = errors.New("lk is empty")
	ErrMaxSizeAttained      = errors.New("image max size attained")
	ErrNotMatch             = errors.New("images and buffer doesn`t match")
	ErrMaxImagesAttained    = errors.New("too much images")
	ErrNilFile              = errors.New("nil file")
)
