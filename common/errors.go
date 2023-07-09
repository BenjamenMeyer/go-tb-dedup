package common

import (
	"errors"
)

var (
	ErrNotOpened                = errors.New("Not Opened")
	ErrNotImplemented           = errors.New("Not Implemented")
	ErrAlreadyOpen              = errors.New("Already open")
	ErrBadState                 = errors.New("Bad State")
	ErrNoFolders                = errors.New("No subfolders")
	ErrNoFiles                  = errors.New("No Files")
	ErrNoDatabaseAvailable      = errors.New("No Database")
	ErrInvalidMessage           = errors.New("Invalid Message Record for Operation")
	ErrUnableToBuildMessageData = errors.New("Unable to reconsitute message data")
	ErrNoMessagesFound          = errors.New("No records found")
	ErrUnableToWrite            = errors.New("Unable to write message")
    ErrMissingHandlerFunc       = errors.New("Missing handler function")
)
