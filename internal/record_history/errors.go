package recordhistory

import "errors"

var (
	ErrRecordIDCollision = errors.New("record ID collision detected")
)
