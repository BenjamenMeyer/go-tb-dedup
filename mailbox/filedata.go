/*
	Copyright (C) 2023 Benjamen R. Meyer

	https://github.com/BenjamenMeyer/go-tb-dedup

	Licensed to the Apache Software Foundation (ASF) under one
	or more contributor license agreements.  See the NOTICE file
	distributed with this work for additional information
	regarding copyright ownership.  The ASF licenses this file
	to you under the Apache License, Version 2.0 (the
	"License"); you may not use this file except in compliance
	with the License.  You may obtain a copy of the License at

	  http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing,
	software distributed under the License is distributed on an
	"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
	KIND, either express or implied.  See the License for the
	specific language governing permissions and limitations
	under the License.
*/
package mailbox

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	//"net/mail"
	"crypto/sha256"
	"os"
	"regexp"
	"strings"

	"github.com/tvanriper/mbox"

	"github.com/BenjamenMeyer/go-tb-dedup/common"
	"github.com/BenjamenMeyer/go-tb-dedup/log"
	"github.com/BenjamenMeyer/go-tb-dedup/storage"
)

const (
	// 16K processing buffer
	DetectRecordBufferLength = 16 * 1024
	// 1K read buffer
	ReadBuffer = 1024
)

var (
	mboxRdMatch = regexp.MustCompile(`^>\s?From - `)
	mboxOMatch  = regexp.MustCompile(`^From - `)
	mboxCLMatch = regexp.MustCompile(`^Content-Length:`)
)

type fileData struct {
	filename Filepath
	mboxMode int
    recordCount int

	fileReader *os.File
	reader     *mbox.MboxReader

	fileWriter *os.File
	writer     *mbox.MboxWriter
}

func NewMailbox() Mailbox {
	return &fileData{}
}

func (fd *fileData) isReader() (result bool) {
	if fd.isOpened() && fd.fileReader != nil {
		result = true
	}
	return
}

func (fd *fileData) isWriter() (result bool) {
	if fd.isOpened() && fd.fileWriter != nil {
		result = true
	}
	return
}

func (fd *fileData) isOpened() (result bool) {
	if len(fd.filename) > 0 || fd.fileReader != nil || fd.reader != nil || fd.fileWriter != nil || fd.writer != nil {
		result = true
	}
	return
}

func (fd *fileData) detectMboxType() (mboxType int, err error) {
	if len(fd.filename) == 0 {
		err = fmt.Errorf("%w: no filename to check", common.ErrBadState)
		return
	}

	// mboxo (default)
	// mboxrd: mboxro + prepends '>' to the `From` used to identify messages
	// mboxcl: mboxrd + content length
	// mboxcl2: mboxro + content length - that is, doesn't add `>` to `From`

	fileReader, fileReaderErr := os.Open(string(fd.filename))
	if fileReaderErr != nil {
		err = fmt.Errorf("%w: Unable read to file '%s' to detect MBox Type", fileReaderErr, fd.filename)
		return
	}
	defer fileReader.Close()

	fileScanner := bufio.NewScanner(fileReader)
	fileScanner.Split(bufio.ScanLines)

	var inMessage bool
	var fromIsPrepended bool
	var hasContentLength bool
	for fileScanner.Scan() {
		fileData := fileScanner.Text()
		isRdMatch := mboxRdMatch.MatchString(fileData)
		isOMatch := mboxOMatch.MatchString(fileData)

		if isRdMatch || isOMatch {
			inMessage = true
		}
		if isRdMatch {
			log.Debug("MBOXRD Detected using line: %s", fileData)
			fromIsPrepended = true
			continue
		}
		if isOMatch {
			log.Debug("MBOXO  or MBOXCL2 Detected using line: %s", fileData)
			fromIsPrepended = false
			continue
		}

		// note: if there's a way to detect the end of message other than the next FROM detection
		//  then do so here to reset inMessage

		if inMessage {
			isClMatch := mboxCLMatch.MatchString(fileData)
			if isClMatch {
				log.Debug("MBOXCL or MBOXCL2 Detected using line: %s", fileData)
				hasContentLength = true
			}
		}
	}

	log.Info("IsPrepended: %t", fromIsPrepended)
	log.Info("HasContentLength: %t", hasContentLength)
	if fromIsPrepended && hasContentLength {
		log.Info("Detected MBOXCL")
		mboxType = mbox.MBOXCL
	} else {
		if fromIsPrepended && !hasContentLength {
			log.Info("Detected MBOXRD")
			mboxType = mbox.MBOXRD
		} else {
			if !fromIsPrepended && hasContentLength {
				log.Info("Detected MBOXCL2")
				mboxType = mbox.MBOXCL2
			} else {
				// !fromIsPrepended && !hasContentLength
				log.Info("Detected MBOXO")
				mboxType = mbox.MBOXO
			}
		}
	}
	return
}

func (fd *fileData) Open(filename Filepath) (err error) {
	if fd.isOpened() {
		err = fmt.Errorf("%w: Already working on %s", common.ErrAlreadyOpen, fd.filename)
		return
	}
	fd.filename = filename

	// unfortunately the mbox package doesn't provide a detector
	// so this does it's own based on the raw file
	fd.mboxMode, err = fd.detectMboxType()
	if err != nil {
		// detectMboxType already sets a good error message
		return
	}

	fileReader, fileReaderErr := os.Open(string(fd.filename))
	if fileReaderErr != nil {
		err = fmt.Errorf("%w: Unable to open file '%s' for reading", fileReaderErr, fd.filename)
		return
	}

	fd.fileReader = fileReader
	fd.reader = mbox.NewReader(fd.fileReader)
	//fd.reader.Type = fd.mboxMode
	return
}

func (fd *fileData) OpenForWrite(filename Filepath, mboxMode int) (err error) {
	if fd.isOpened() {
		err = fmt.Errorf("%w: Already working on %s", common.ErrAlreadyOpen, fd.filename)
		return
	}
	fd.filename = filename

	fileWriter, fileWriterErr := os.OpenFile(string(fd.filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0600)
	if fileWriterErr != nil {
		err = fmt.Errorf("%w: Unable open file '%s' for writing", fileWriterErr, fd.filename)
		return
	}
	fd.fileWriter = fileWriter
	fd.writer = mbox.NewWriter(fd.fileWriter)
	fd.writer.Type = fd.mboxMode
	return
}

func (fd *fileData) Close() (err error) {
	if fd.reader != nil {
		fd.reader = nil
	}
	if fd.fileReader != nil {
		fd.fileReader.Close()
		fd.fileReader = nil
	}
	if fd.writer != nil {
		fd.writer = nil
	}
	if fd.fileWriter != nil {
		fd.fileWriter.Close()
		fd.fileWriter = nil
	}
	fd.filename = ""
	return
}

func makeMessageId(location string, index int) (messageId string) {
	// if running on Windows then there may be `:` at the start - e.g c:\
	// for data safety remove the colons so it's all clean; this will essentially
	// be a NOP on other platforms
	removeColons := strings.ReplaceAll(location, ":", "-")
	// Remove the slashes (everything)
	removeForwardSlashes := strings.ReplaceAll(removeColons, "/", "-")
	// Remove the backslashes (Windows only)
	removeBackwardSlashes := strings.ReplaceAll(removeForwardSlashes, "", "-")
	// now build a string out of it
	// zero-pad it for logging consistency
	messageId = fmt.Sprintf("%s_%010d", removeBackwardSlashes, index)
	return
}

func (fd *fileData) getMessages(msgHandler func(filename string, index int, msgId string, fromStr string, hash []byte, msgData []byte)) (err error) {
	if !fd.isReader() {
		err = fmt.Errorf("%w: Attempting to read when reader is not opened", common.ErrBadState)
		return
	}

    if msgHandler == nil {
        err = fmt.Errorf("%w: Missing message handler", common.ErrMissingHandlerFunc)
        return
    }

	mailBytes := bytes.NewBuffer([]byte{})
	var index int
	var fromStr string

	fromStr, err = fd.reader.NextMessage(mailBytes)
	for err == nil {
		/*
		   // if we need to read the actual message then do the following:
		   msg, e := mail.ReadMessage(bytes.NewBuffer(mailBytes.Bytes()))
		   if e != nil {
		       log.Error("ERROR: Error reading message - %v", e)
		       continue
		   }
		*/
		if mailBytes.Len() > 0 {
			// generate the message hash
			rawMsgData := mailBytes.Bytes()
			msgData := append([]byte(nil), rawMsgData...) // until supported by copy() in Golang 1.20+
			hashData := sha256.Sum256(rawMsgData)

			// get around golang limitations since [32]byte doesn't convert to []byte
			var msgHash []byte = make([]byte, len(hashData))
			copy(msgHash[:], hashData[:])

            msgHandler(
                string(fd.filename),
                // MBOX format doesn't provide an index value so just make one up that can be utilized
                index,
                makeMessageId(string(fd.filename), index),
                fromStr,
                msgHash,
                msgData,
            ) 

		}

		mailBytes.Reset()
		index += 1

		fromStr, err = fd.reader.NextMessage(mailBytes)
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// GetMessages is useful to read the data and get it back as an array
// This does not work well on large files though
func (fd *fileData) GetMessages() (msgs []storage.MailData, err error) {
    // using the common message reader generate a storage.MailData object 
    msgHandler := func(filename string, index int, msgId string, fromStr string, msgHash []byte, msgData []byte) {
        msg := storage.NewMessageData()
        msg.SetHash(msgHash)
        msg.SetLocation(filename)
        msg.SetLocationIndex(index)
        msg.SetFrom(fromStr)
        msg.SetData(msgData)
        // MBOX format doesn't provide an index value so just make one up that can be utilized
        msg.SetMessageId(msgId)
        msgs = append(msgs, msg)
    }

    // common msg reader
    err = fd.getMessages(msgHandler)
    return
}

// GetMessagesReferences is useful to determine how many records there are
// and get their hashes for quick deduplication
func (fd *fileData) GetMessagesReferences() (msgs []storage.MailRecord, err error) {
    // using the common message reader generate a storage.MailRecord object 
    msgHandler := func(filename string, index int, msgId string, fromStr string, msgHash []byte, msgData []byte) {
        // ignores msgData and fromStr
        msg := MailRecord{
            Index: index,
            Hash: msgHash,
            MsgId: msgId,
            Filename: filename,
        }
        msgs = append(msgs, msg)
    }

    // common msg reader
    err = fd.getMessages(msgHandler)
    return
}

func (fd *fileData) AddMessage(message storage.MailData) (err error) {
	if !fd.isWriter() {
		err = fmt.Errorf("%w: Attempting to read when writer is not opened", common.ErrBadState)
		return
	}

	fromStr, fromErr := message.GetFrom()
	if fromErr != nil || len(fromStr) == 0 {
		err = fmt.Errorf("%w: Unable to retrieve message 'from' line", fromErr)
		return
	}

	rawData, rawDataErr := message.GetData()
	if rawDataErr != nil || len(rawData) == 0 {
		err = fmt.Errorf("%w: Unable to retrieve message data", rawDataErr)
		return
	}
	// copy and add a terminating \n as that is what the mbox writer expects
	dataCopy := append([]byte(nil), rawData...)
	dataCopy = append(dataCopy, '\n')
	msgReader := bytes.NewReader(dataCopy)

	writeErr := fd.writer.WriteMail(fromStr, msgReader)
	if writeErr != nil {
		err = fmt.Errorf("%w: error writing message - %v", common.ErrUnableToWrite, writeErr)
	}
	return
}

var _ Mailbox = &fileData{}
