package recordhistory

import (
	"log/slog"

	"github.com/cespare/xxhash"
)

type RecordMode uint8

type RecordUnity struct {
	MsgHash uint64
	Mode    RecordMode
}

type RecordHistory struct {
	Records map[uint64]*RecordUnity
}

type IRecordHistory interface {
	AddRecord(id uint64, msg string, mode RecordMode) error
	RemoveRecord(id uint64)
	GetRecord(id uint64) *RecordUnity
}

const (
	logOnce RecordMode = iota
	logOnChange
)

func NewRecordHistory() IRecordHistory {
	return &RecordHistory{
		Records: make(map[uint64]*RecordUnity),
	}
}

func GenHash(s string) uint64 {
	return xxhash.Sum64String(s)
}

func (r *RecordHistory) AddRecord(id uint64, msg string, mode RecordMode) error {
	if r.GetRecord(id) != nil {
		return ErrRecordIDCollision
	}
	r.Records[id] = &RecordUnity{
		MsgHash: GenHash(msg),
		Mode:    mode,
	}
	return nil
}

func (r *RecordHistory) RemoveRecord(id uint64) {
	if r.GetRecord(id) == nil {
		slog.Debug("Trying to remove non-existing record")
	}
	delete(r.Records, id)
}

func (r *RecordHistory) GetRecord(id uint64) *RecordUnity {
	return r.Records[id]
}
