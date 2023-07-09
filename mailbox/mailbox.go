package mailbox

import (
	"github.com/BenjamenMeyer/go-tb-dedup/storage"
)

type Folderpath string
type Filepath string

type Mailbox interface {
	Open(filename Filepath) error
	OpenForWrite(filename Filepath, mboxMode int) error
	Close() error

	GetMessages() ([]storage.MailData, error)
	AddMessage(message storage.MailData) error

    GetMessagesReferences() (msgs []storage.MailRecord, err error)
}

type Mailfolder interface {
	Open(folder Folderpath) error
	Close() error

	GetCounts() (fileCount int, folderCount int, err error)
	GetMailFiles() ([]Filepath, error)
	GetSubfolders() ([]Folderpath, error)
}
