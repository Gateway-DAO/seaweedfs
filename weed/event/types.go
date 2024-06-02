package event

/**
 * MARK: Generic Events
 */

type EventInterface interface {
	Id() string

	Metadata() interface{}

	SetTimestamps(t int64)
	CreatedAt() int64
	LastUpdated() int64
	LastTouched() int64
}

type Event struct {
	LtTimestamp int64 `json:"lt_timestamp"`
	LuTimestamp int64 `json:"lu_timestamp"`
	CaTimestamp int64 `json:"ca_timestamp"`
}

func (e Event) LastTouched() int64 { return e.LtTimestamp }
func (e Event) LastUpdated() int64 { return e.LuTimestamp }
func (e Event) CreatedAt() int64   { return e.CaTimestamp }
func (e Event) SetTimestamps(t int64) {
	e.LtTimestamp = t
	e.LuTimestamp = t
	e.CaTimestamp = t
}

/**
 *	MARK: FileEvents
 */

type FileStatus int32

const (
	ASSIGNED FileStatus = iota
	UPLOADED
	DELETED
)

type FileEvent struct {
	Event

	Fid    string     `json:"fid"`
	Hash   string     `json:"hash"`
	Status FileStatus `json:"status"`
}

func (fe FileEvent) Id() string { return fe.Fid }

func (e FileEvent) Metadata() interface{} {
	return struct {
		Hash   string     `json:"hash"`
		Status FileStatus `json:"status"`
	}{
		Hash:   e.Hash,
		Status: e.Status,
	}
}
