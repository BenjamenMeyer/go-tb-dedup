package storage

// Data storage
type MailData interface {
	SetMessageId(id string) error
	GetMessageId() (string, error)

	SetLocation(location string) error
	GetLocation() (string, error)

	// NOTE: Location Index is only guaranteed to be accurate while data is not modified
	// If a toold is modifying data then it should process based on the reverse index values
	SetLocationIndex(index int) error
	GetLocationIndex() (int, error)

	SetHash(hash []byte) error
	GetHash() ([]byte, error)

	SetFrom(from string) error
	GetFrom() (string, error)

	SetData(data []byte) error
	GetData() ([]byte, error)
}

type MailStorage interface {
	Open() error
	Close() error

	AddMessage(message MailData) error
	GetMessage(id string) (MailData, error)

	GetMessagesByLocation(location string) ([]MailData, error)
	GetMessagesByHash(hash []byte) ([]MailData, error)
	GetUniqueMessages() (msgs []MailData, err error)

	AddRecord(message MailRecord) error
	GetRecord(id string) (MailRecord, error)

	GetRecordsByLocation(location string) ([]MailRecord, error)
	GetRecordsByHash(hash []byte) ([]MailRecord, error)
	GetUniqueRecords() (msgs []MailRecord, err error)
}

type MailRecord struct {
    Index int
    Hash []byte
    //FileOffset int
    MsgId string
    Filename string
}
