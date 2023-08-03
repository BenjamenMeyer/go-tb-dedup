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
package dedup

import (
	"errors"
	"fmt"
	"os"
    "sync"
	"time"

    "github.com/tvanriper/mbox"

	"github.com/BenjamenMeyer/go-tb-dedup/common"
	"github.com/BenjamenMeyer/go-tb-dedup/log"
	"github.com/BenjamenMeyer/go-tb-dedup/mailbox"
	"github.com/BenjamenMeyer/go-tb-dedup/storage"
)

var (
    logUnit = "DEDUP"
)

type DataProcessor struct {
	mailStore   storage.MailStorage
	mailData    chan storage.MailData
    mailRecords chan storage.MailRecord
	mailMode    int
	mailCounter chan int
	done        chan bool
	errs        chan error
    completionChan chan string
    userCapture chan os.Signal
    baseLogUnit string
    baseLocation string
    outputMbox string
}

func NewDataProcessor(mailStore storage.MailStorage, storageQueueDepth uint, userCapture chan os.Signal) *DataProcessor {
    dp := &DataProcessor{
		mailStore:   mailStore,
		mailData:    make(chan storage.MailData, storageQueueDepth),
        mailRecords: make(chan storage.MailRecord, storageQueueDepth),
	    mailCounter: make(chan int, storageQueueDepth),
		done:        make(chan bool, 0),
		errs:        make(chan error, storageQueueDepth),
        completionChan:  make(chan string, 0),
        userCapture: userCapture,
        // skip: baseLogUnit
        // skip: baseLocation
        // skip: outputMBox
    }
    return dp
}

func (dp *DataProcessor) doProcessFolder() {
	log.UnitInfo(logUnit, "Starting to process the mailbox")
	go dp.processFolder("0", mailbox.Folderpath(dp.baseLocation), dp.completionChan)
	// wait for the completion
ProcessingLoop:
	for {
		select {
		case <-dp.completionChan:
			log.UnitInfo(dp.baseLogUnit, "Received completion notice.")
			log.UnitInfo(dp.baseLogUnit, "Waiting for input processing to finish")
			log.UnitInfo(dp.baseLogUnit, "Finished processing incoming requests")
			log.UnitInfo(dp.baseLogUnit, "Processing finished")

			// process the results and generate the final set, and what is to be deleted
			// clean out duplicated data
			log.UnitInfo(dp.baseLogUnit, "All complete")
			break ProcessingLoop

		case <-dp.userCapture:
			log.UnitInfo(dp.baseLogUnit, "Caught user termination")
			log.UnitInfo(dp.baseLogUnit, "Signaling termination")
			close(dp.done)
			log.UnitInfo(dp.baseLogUnit, "Waiting for terminal completion")
			// close doneChan and wait for completion to return
			return
		}
	}
}

func (dp *DataProcessor) doOutput(outputMbox string) {
	// taking a random guess
	dp.mailMode = mbox.MBOXCL
	go dp.writeUniqueMessages(
		mailbox.Filepath(outputMbox),
		dp.completionChan,
	)
WriterLoop:
	for {
		select {
		case <-dp.completionChan:
			log.UnitInfo(dp.baseLogUnit, "Received completion notice.")
			log.UnitInfo(dp.baseLogUnit, "Waiting 30 seconds before closing")
			time.Sleep(30 * time.Second)
			log.UnitInfo(dp.baseLogUnit, "Finished processing incoming requests")
			log.UnitInfo(dp.baseLogUnit, "Processing finished")

			// process the results and generate the final set, and what is to be deleted
			// clean out duplicated data
			log.UnitInfo(dp.baseLogUnit, "All complete")
			break WriterLoop

		case <-dp.userCapture:
			log.UnitInfo(dp.baseLogUnit, "Caught user termination")
			log.UnitInfo(dp.baseLogUnit, "Signaling termination")
			close(dp.done)
			log.UnitInfo(dp.baseLogUnit, "Waiting for terminal completion")
			// close doneChan and wait for completion to return
			return
		}
	}
}

func (dp *DataProcessor) mailBoxHandler() {
	logUnit := "MAILBOX HANDLER"
	log.UnitInfo(logUnit, "Starting")
	var msgCount int
	var expectedCount int
	for {
		select {
		case count := <-dp.mailCounter:
			// new func context to recover errors and prevent crashes
			expectedCount += count
			log.UnitInfo(logUnit, "Expecting %d more records for a total of %d records", count, expectedCount)
		case msg := <-dp.mailData:
			msgCount += 1
			// new func context to recover errors and prevent crashes
			func() {
				defer func() {
					if r := recover(); r != nil {
						dp.errs <- fmt.Errorf(
							"%s: Recovered from error processing message %#v. %#v",
							logUnit,
							msg,
							r,
						)
					}
				}()
				dp.mailStore.AddMessage(msg)
			}()
        case msg := <-dp.mailRecords:
            msgCount += 1
            func() {
				defer func() {
					if r := recover(); r != nil {
						dp.errs <- fmt.Errorf(
							"%s: Recovered from error processing message %#v. %#v",
							logUnit,
							msg,
							r,
						)
					}
				}()
                dp.mailStore.AddRecord(msg)
            }()
		case <-dp.done:
			log.UnitInfo(logUnit, "Received %d of %d records", msgCount, expectedCount)
			log.UnitInfo(logUnit, "Finished processing incoming requests")
			log.UnitInfo(logUnit, "Terminating")
			return
		}
	}
}

func (dp *DataProcessor) errHandler() {
	logUnit := "ERROR HANDLER"
	log.UnitInfo(logUnit, "Starting")
	for {
		select {
		case err := <-dp.errs:
			log.UnitError(logUnit, "%#v", err)
		case <-dp.done:
			log.UnitInfo(logUnit, "Finished processing incoming requests")
			log.UnitInfo(logUnit, "Terminating")
			return
		}
	}
}

func (dp *DataProcessor) processFile(id string, location mailbox.Filepath, completion chan string, inputGroup *sync.WaitGroup) {
    inputGroup.Done()
	logUnit := fmt.Sprintf("FILE HANDLER:%s", string(location))
	// split a goroutine off for each message
	log.UnitInfo(logUnit, "Opening mailbox file")
	mailFile := mailbox.NewMailbox()
	mailFileErr := mailFile.Open(location)
	if mailFileErr != nil {
		dp.errs <- fmt.Errorf("%w: Unable to process file '%s'", mailFileErr, string(location))
		completion <- id
		return
	}
	defer mailFile.Close()

	// note: the GetMessages() generates the hashes on all the messages it reads
	// therefore there is no point in splitting off another goroutine to handle
	// each message as there is nothing left to process after this point other than
	// recording it to storage which is already being done in another goroutine
	// the alternative would be to copy and expose the RAW message data to this call
	// which could then be split to another goroutine to generate the hash and send the
	// final message to storage
	log.UnitInfo(logUnit, "Reading messages")
    // TODO: Convert this to generate an index, iterate over the indexes, and read one message at a time
	//msgs, msgsErr := mailFile.GetMessages()
	msgs, msgsErr := mailFile.GetMessagesReferences()
	if msgsErr != nil {
		dp.errs <- fmt.Errorf("%w: Unable to retrieve messages from location '%s'", msgsErr, string(location))
		completion <- id
		return
	}

	log.UnitInfo(logUnit, "Sending messages to storage")
	var msgCount int
	var counter int
	for _, msg := range msgs {
        //dp.mailData <- msg
        dp.mailRecords <- msg
		msgCount += 1
		// as this is a tight loop, let's not overload things and lock the system
		// therefore yield to the scheduler
		if counter > 5 {
			time.Sleep(0)
			counter = 0
		} else {
			counter++
		}
	}
	dp.mailCounter <- msgCount

	log.UnitInfo(logUnit, "Found %d messages", msgCount)
	log.UnitInfo(logUnit, "completed processing messages")
	completion <- id
	log.UnitInfo(logUnit, "notified parent; terminating")
	return
}

func (dp *DataProcessor) processFolder(id string, location mailbox.Folderpath, completion chan string) {
	logUnit := fmt.Sprintf("FOLDER HANDLER:%s", string(location))
	log.UnitInfo(logUnit, "Opening mailbox folder")
	mailFolder := mailbox.NewMailFolder()
	mailFolderErr := mailFolder.Open(location)
	if mailFolderErr != nil {
		dp.errs <- fmt.Errorf("%w: Unable to process folder '%s'", mailFolderErr, string(location))
		completion <- id
		return
	}
	defer mailFolder.Close()

	log.UnitInfo(logUnit, "accessing content")
	var totalCount int = 10000
	fileCount, folderCount, counterErr := mailFolder.GetCounts()
	if counterErr != nil {
		dp.errs <- fmt.Errorf("%w: Unable to get file and folder counts. Defaulting to %d", counterErr, totalCount)
	} else {
		totalCount = fileCount + folderCount
	}

	// it is known how many entries to watch for
	subCompletion := make(chan string, totalCount)

	// split a goroutine off for each mbox in the folder
	log.UnitInfo(logUnit, "generating works for files")
	files, filesErr := mailFolder.GetMailFiles()
	if filesErr != nil {
		if !errors.Is(filesErr, common.ErrNoFiles) {
			dp.errs <- fmt.Errorf("%w: Error processing mail box files", filesErr)
		}
	} else {
		// do the loop in another go routine
		go func() {
            var inputGroup  sync.WaitGroup
			for fileEntryIndex, fileEntry := range files {
				// there is no recursion in this, so add file to denote the chain difference
				// and then use the index to denote which file
				fileEntryId := fmt.Sprintf("%s:file:%d", id, fileEntryIndex)
                inputGroup.Add(1)
				go dp.processFile(fileEntryId, fileEntry, subCompletion, &inputGroup)
			}
            log.UnitInfo(logUnit, "Waiting on files to finish processing")
            inputGroup.Wait()
            log.UnitInfo(logUnit, "Files completed processing")
		}()
	}

	// split a goroutine off for each subfolder in the folder
	log.UnitInfo(logUnit, "generating works for sub-folders")
	subdirs, subdirsErr := mailFolder.GetSubfolders()
	if subdirsErr != nil {
		if !errors.Is(subdirsErr, common.ErrNoFolders) {
			dp.errs <- fmt.Errorf("%w: Error processing mail box subfolders", subdirsErr)
		}
	} else {
		// do the loop in another go routine
		go func() {
			for subdirEntryIndex, subdirEntry := range subdirs {
				// there may be many folders, and it doesn't make sense to have `folder` repeated
				// many times as it is recursive so just use the id and index
				subdirEntryId := fmt.Sprintf("%s:%d", id, subdirEntryIndex)
				go dp.processFolder(subdirEntryId, subdirEntry, subCompletion)
			}
		}()
	}

	log.UnitInfo(logUnit, "waiting on processing of files and sub-folders")

	// collect the results of the spun off goroutines above and return only when all children
	// have returned or the system has terminated
	var completedRoutines []string
	watcher := time.NewTimer(500 * time.Millisecond) // every half-second
	for {
		select {
		case response := <-subCompletion:
			completedRoutines = append(completedRoutines, response)
			if len(completedRoutines) >= totalCount {
				log.UnitInfo(logUnit, "Finished processing incoming requests")
				completion <- id
				log.UnitInfo(logUnit, "notified parent; terminating")
				return
			}
		case <-dp.done:
			dp.errs <- fmt.Errorf("%s: Terminating before all requests finished. Only %d of %d finished", logUnit, len(completedRoutines), totalCount)
			completion <- id
			log.UnitInfo(logUnit, "notified parent; terminating")
			return
		case <-watcher.C:
			log.UnitInfo(logUnit, "%d of %d completed", len(completedRoutines), totalCount)
		}
	}
}

func (dp *DataProcessor) writeUniqueMessages(targetLocation mailbox.Filepath, completion chan string) {
	//
	logUnit := "UNIQUE WRITER"
	defer func() {
		completion <- "writer"
	}()
	uniqueMsgs, getUniqueMsgErr := dp.mailStore.GetUniqueMessages()
	if getUniqueMsgErr != nil {
		log.UnitError(logUnit, "Error retrieving the : %#v", getUniqueMsgErr)
		return
	}

	// Open the mail box storage
	mailFile := mailbox.NewMailbox()
	mailFileErr := mailFile.OpenForWrite(targetLocation, dp.mailMode)
	if mailFileErr != nil {
		dp.errs <- fmt.Errorf("%w: Unable to open file '%s' for writing", mailFileErr, string(targetLocation))
		return
	}

	for _, msg := range uniqueMsgs {
		// write the unique message
		recordErr := mailFile.AddMessage(msg)
		if recordErr != nil {
			dp.errs <- fmt.Errorf("%w: Error writing record to target location (%s) - %v", recordErr, string(targetLocation), msg)
		}
	}
}

func (dp *DataProcessor) Run(baseLocation string, outputMbox string, processingGroup *sync.WaitGroup) {
    defer processingGroup.Done()
    dp.baseLocation = baseLocation
    dp.outputMbox = outputMbox
	dp.baseLogUnit = fmt.Sprintf("BASE HANDLER:%s", dp.baseLocation)


	// ensure things get cleaned up on application exit
	cleanup := func() {
		close(dp.done)
		time.Sleep(1 * time.Second)
		close(dp.mailData)
	}
	defer cleanup()

	// spin off a go-routine to async handle the mail storage recording
	go dp.mailBoxHandler()

	// spin off a go-routine to async handle error logging
	go dp.errHandler()

    dp.doProcessFolder()

	if len(outputMbox) == 0 {
		fmt.Printf("No output location specified. Terminating.")
		return
	}

    dp.doOutput(outputMbox)
}
