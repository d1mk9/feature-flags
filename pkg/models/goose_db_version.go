package models

//go:generate reform

// GooseDbVersion represents a row in goose_db_version table.
//
//reform:goose_db_version
type GooseDbVersion struct {
	ID        int32  `reform:"id,pk"`
	VersionID int64  `reform:"version_id"`
	IsApplied bool   `reform:"is_applied"`
	Tstamp    []byte `reform:"tstamp"` // FIXME unhandled database type "timestamp without time zone"
}
