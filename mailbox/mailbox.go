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
