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
