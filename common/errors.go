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
