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
package main

import (
	"flag"
	"fmt"
    golog "log"
    "os"
	"os/signal"
    "sync"
	"syscall"

	"github.com/BenjamenMeyer/go-tb-dedup/dedup"
	"github.com/BenjamenMeyer/go-tb-dedup/log"
	"github.com/BenjamenMeyer/go-tb-dedup/storage"
)

func main() {
	var baseLocation = flag.String("location", "", "Location of Mbox folders to search")
	var storageLocation = flag.String("hash-storage", "", "Optional location to store data")
	var storageQueueDepth = flag.Uint("queue-depth", 1000000, "Queue depth for relaying messages to storage")
	var outputMbox = flag.String("output-mbox", "", "File to write the MBox data to")
	var logFileName = flag.String("log-file", ".go-tb-dedup.log", "File to write log data to")

	flag.Parse()

	// check for required parameters
	if len(*baseLocation) == 0 {
		fmt.Printf("location is a required parameter.\n\n")
		flag.Usage()
		return
	}

    logFile, logFileErr := os.OpenFile(*logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
    if logFileErr != nil {
        fmt.Printf("failed to open log for writing: %#v\n\n", logFileErr)
        return
    }
    defer logFile.Close()
    golog.SetOutput(logFile)

	logUnit := "MAIN"

	mailDataStorage := storage.NewSqliteStorage(*storageLocation)
	mailDataStorageErr := mailDataStorage.Open()
	if mailDataStorageErr != nil {
		log.UnitError(logUnit, "Unable to create hash storage (backing: %s): %#v", *storageLocation, mailDataStorageErr)
		return
	}
	defer mailDataStorage.Close()

	log.UnitInfo(logUnit, "Using location: %s", *baseLocation)
	log.UnitInfo(logUnit, "Using Storage Queue Depth: %d", *storageQueueDepth)

	// capture CTRL+C, CTRL+BRK, and SIG TERM
    userCapture := make(chan os.Signal)
    defer close(userCapture)

	signal.Notify(userCapture, os.Interrupt, syscall.SIGTERM)

    var processing sync.WaitGroup
    processing.Add(1)
	dp := dedup.NewDataProcessor(mailDataStorage, *storageQueueDepth, userCapture)
    go dp.Run(*baseLocation, *outputMbox, &processing)
    processing.Wait()
}
