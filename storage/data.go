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

type messageData struct {
	id            string
	location      string
	locationIndex int
	hash          []byte
	from          string
	data          []byte
}

func NewMessageData() MailData {
	return &messageData{}
}

func (md *messageData) SetMessageId(id string) (err error) {
	md.id = id
	return
}
func (md *messageData) GetMessageId() (id string, err error) {
	id = md.id
	return
}

func (md *messageData) SetLocation(location string) (err error) {
	md.location = location
	return
}
func (md *messageData) GetLocation() (location string, err error) {
	location = md.location
	return
}

func (md *messageData) SetLocationIndex(index int) (err error) {
	md.locationIndex = index
	return
}
func (md *messageData) GetLocationIndex() (index int, err error) {
	index = md.locationIndex
	return
}

func (md *messageData) SetHash(hash []byte) (err error) {
	md.hash = hash
	return
}
func (md *messageData) GetHash() (hash []byte, err error) {
	hash = md.hash
	return
}

func (md *messageData) SetFrom(from string) (err error) {
	md.from = from
	return
}
func (md *messageData) GetFrom() (from string, err error) {
	from = md.from
	return
}

func (md *messageData) SetData(data []byte) (err error) {
	md.data = data
	return
}
func (md *messageData) GetData() (data []byte, err error) {
	data = md.data
	return
}

var _ MailData = &messageData{}
