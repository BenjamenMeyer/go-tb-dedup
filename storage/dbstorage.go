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

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/BenjamenMeyer/go-tb-dedup/common"
	"github.com/BenjamenMeyer/go-tb-dedup/log"
)

const (
	SQL_ADD_MESSAGE             = "INSERT INTO messages (messageid, location, locationindex, hash, msgFrom, message) VALUES(?, ?, ?, ?, ?, ?)"
	SQL_GET_MESSAGE_BY_ID       = "SELECT messageid, location, locationindex, hash, msgFrom, message FROM messages WHERE messageid = ?"
	SQL_GET_MESSAGE_BY_LOCATION = "SELECT messageid, location, locationindex, hash, message msgFrom FROM messages WHERE location = ?"
	SQL_GET_MESSAGE_BY_HASH     = "SELECT messageid, location, locationindex, hash, message msgFrom FROM messages WHERE hash = ?"
	SQL_GET_UNIQUE_MESSAGES     = "SELECT messageid, location, locationindex, hash, msgFrom, message FROM messages WHERE hash IN (SELECT DISTINCT hash FROM messages)"

    SQL_ADD_RECORD              = "INSERT INTO records (messageid, location, locationindex, hash) VALUES(?, ?, ?, ?)"
    SQL_GET_RECORD_BY_ID        = "SELECT messageid, location, locationindex, hash FROM records WHERE messageid = ?"
    SQL_GET_RECORD_BY_LOCATION  = "SELECT messageid, location, locationindex, hash FROM records WHERE location = ?"
	SQL_GET_RECORD_BY_HASH      = "SELECT messageid, location, locationindex, hash FROM records WHERE hash = ?"
    SQL_GET_UNIQUE_RECORDS      = "SELECT messageid, location, locationindex, hash FROM records WHERE hash IN (SELECT DISTINCT hash FROM records)"
)

type schemaEntry struct {
	version int
	entries []string
}

var SQL_SCHEMA = []schemaEntry{
	{
		version: 0,
		entries: []string{
			"CREATE TABLE message_schema(version INT)",
			"INSERT INTO message_schema (version) VALUES(0)",
		},
	},
	{
		version: 1,
		entries: []string{
			"CREATE TABLE messages(messageid TEXT NOT NULL, location TEXT NOT NULL, locationindex INT, hash TEXT NOT NULL, msgFrom TEXT NOT NULL, message TEXT NOT NULL, PRIMARY KEY(messageid, hash))",
		},
	},
    {
        version: 2,
        entries: []string{
            "CREATE TABLE records(messageid TEXT NOT NULL, location TEXT NOT NULL, locationindex INT, hash TEXT NOT NULL, PRIMARY KEY(messageid, hash))",
        },
    },
}

type dbStorage struct {
	dbDiskLocation string
	sqlDb          *sql.DB

	sqlAddMessage           *sql.Stmt
	sqlGetMessageById       *sql.Stmt
	sqlGetMessageByLocation *sql.Stmt
	sqlGetMessageByHash     *sql.Stmt
	sqlGetUniqueMessages    *sql.Stmt

	sqlAddRecord           *sql.Stmt
	sqlGetRecordById       *sql.Stmt
	sqlGetRecordByLocation *sql.Stmt
	sqlGetRecordByHash     *sql.Stmt
	sqlGetUniqueRecords    *sql.Stmt
}

// NOTE: provide a location for debug purposes or if in low memory systems
// in general, it should be sufficient to operate in memory entirely
func NewSqliteStorage(diskLocation string) MailStorage {
	if len(diskLocation) == 0 {
		diskLocation = ":memory:"
	}
	return &dbStorage{
		dbDiskLocation: diskLocation,
	}
}

func (db *dbStorage) Open() (err error) {
	if db.sqlDb != nil {
		err = fmt.Errorf("%w: Already processing '%s'", common.ErrAlreadyOpen, db.dbDiskLocation)
		return
	}

	if len(db.dbDiskLocation) == 0 {
		err = fmt.Errorf("%w: Unknown disk location. Bad allocation status", common.ErrBadState)
		return
	}

	sqliteDb, sqliteDbErr := sql.Open("sqlite3", db.dbDiskLocation)
	if sqliteDbErr != nil {
		err = fmt.Errorf("%w: Failed to open '%s' with Sqlite", sqliteDbErr, db.dbDiskLocation)
		return
	}

	// save off the DB
	db.sqlDb = sqliteDb

	// initialize the DB Schema - returns a formatted error
	err = db.migrate()
	if err != nil {
		return
	}

	// create any prepared statements here
	err = db.prepareQueries()
	if err != nil {
		return
	}

	// all done!
	return
}

func (db *dbStorage) Close() (err error) {
	if db.sqlDb != nil {
		// close prepared statements
		if db.sqlAddMessage != nil {
			closeErr := db.sqlAddMessage.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Add Message", closeErr)
				return
			}
		}
		db.sqlAddMessage = nil

		if db.sqlGetMessageById != nil {
			closeErr := db.sqlGetMessageById.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Get Message By ID", closeErr)
				return
			}
		}
		db.sqlGetMessageById = nil

		if db.sqlGetMessageByLocation != nil {
			closeErr := db.sqlGetMessageByLocation.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Get Message By Location", closeErr)
				return
			}
		}
		db.sqlGetMessageByLocation = nil

		if db.sqlGetMessageByHash != nil {
			closeErr := db.sqlGetMessageByHash.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Get Message By Hash", closeErr)
				return
			}
		}
		db.sqlGetMessageByHash = nil

		if db.sqlGetUniqueMessages != nil {
			closeErr := db.sqlGetUniqueMessages.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Get Unique Messages", closeErr)
				return
			}
		}

		if db.sqlAddRecord != nil {
			closeErr := db.sqlAddRecord.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Add Record", closeErr)
				return
			}
		}
		db.sqlAddRecord = nil

		if db.sqlGetRecordById != nil {
			closeErr := db.sqlGetRecordById.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Get Record By ID", closeErr)
				return
			}
		}
		db.sqlGetRecordById = nil

		if db.sqlGetRecordByLocation != nil {
			closeErr := db.sqlGetRecordByLocation.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Get Record By Location", closeErr)
				return
			}
		}
		db.sqlGetRecordByLocation = nil

		if db.sqlGetRecordByHash != nil {
			closeErr := db.sqlGetRecordByHash.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Get Record By Hash", closeErr)
				return
			}
		}
		db.sqlGetRecordByHash = nil

		if db.sqlGetUniqueRecords != nil {
			closeErr := db.sqlGetUniqueRecords.Close()
			if closeErr != nil {
				err = fmt.Errorf("%w: Error while closing the statement - Get Unique Records", closeErr)
				return
			}
		}

		// close database
		closeErr := db.sqlDb.Close()
		if closeErr != nil {
			err = fmt.Errorf("%w: Error while closing the database", closeErr)
			return
		}
		db.sqlDb = nil
	}

	return
}

func (db *dbStorage) migrate() (err error) {
	for _, schema := range SQL_SCHEMA {
		for sqlIndex, sqlStmt := range schema.entries {
			_, sqlStmtErr := db.sqlDb.Exec(sqlStmt)
			if sqlStmtErr != nil {
				err = fmt.Errorf("%w: Failed to migrate schema version (%d) sql index (%d): '%s'", sqlStmtErr, schema.version, sqlIndex, sqlStmt)
				return
			}
		}

		// as long as there isn't an error, then update the schema version
		_, schemaVersionErr := db.sqlDb.Exec(
			"UPDATE message_schema SET version = ?",
			schema.version,
		)
		if schemaVersionErr != nil {
			err = fmt.Errorf("%w: Failed to set schema version to %d", schemaVersionErr, schema.version)
			return
		}
	}
	return
}

func (db *dbStorage) prepareQueries() (err error) {
	makePrepared := func(sqlStatement string) (stmt *sql.Stmt, err error) {
		stmt, stmtErr := db.sqlDb.Prepare(sqlStatement)
		if stmtErr != nil {
			stmt = nil
			err = fmt.Errorf("%w: Unable to create prepared query '%s'", stmtErr, sqlStatement)
		}
		return
	}

	if db.sqlAddMessage, err = makePrepared(SQL_ADD_MESSAGE); err != nil {
		return
	}
	if db.sqlGetMessageById, err = makePrepared(SQL_GET_MESSAGE_BY_ID); err != nil {
		return
	}
	if db.sqlGetMessageByLocation, err = makePrepared(SQL_GET_MESSAGE_BY_LOCATION); err != nil {
		return
	}
	if db.sqlGetMessageByHash, err = makePrepared(SQL_GET_MESSAGE_BY_HASH); err != nil {
		return
	}
	if db.sqlGetUniqueMessages, err = makePrepared(SQL_GET_UNIQUE_MESSAGES); err != nil {
		return
	}

	if db.sqlAddRecord, err = makePrepared(SQL_ADD_RECORD); err != nil {
		return
	}
	if db.sqlGetRecordById, err = makePrepared(SQL_GET_RECORD_BY_ID); err != nil {
		return
	}
	if db.sqlGetRecordByLocation, err = makePrepared(SQL_GET_RECORD_BY_LOCATION); err != nil {
		return
	}
	if db.sqlGetRecordByHash, err = makePrepared(SQL_GET_RECORD_BY_HASH); err != nil {
		return
	}
	if db.sqlGetUniqueRecords, err = makePrepared(SQL_GET_UNIQUE_RECORDS); err != nil {
		return
	}

	// all done
	return
}

/*****************
 *** Mail Data ***
 *****************/
func (db *dbStorage) AddMessage(message MailData) (err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	msgId, msgIdErr := message.GetMessageId()
	location, locationErr := message.GetLocation()
	locationIndex, locationIndexErr := message.GetLocationIndex()
	msgHash, msgHashErr := message.GetHash()
	msgFrom, msgFromErr := message.GetFrom()
	msgData, msgDataErr := message.GetData()
	if msgIdErr != nil || locationErr != nil || locationIndexErr != nil || msgHashErr != nil || msgFromErr != nil || msgDataErr != nil {
		err = fmt.Errorf(
			"%w: unable to get all message parts (id error: %#v, location error: %#v, location index error %#v, hash error: %#v, msgFrom error: %#v, data error: %#v",
			common.ErrInvalidMessage,
			msgIdErr,
			locationErr,
			locationIndexErr,
			msgHashErr,
			msgFromErr,
			msgDataErr,
		)
		return
	}

	tx, txErr := db.sqlDb.Begin()
	if txErr != nil {
		err = fmt.Errorf("%w: Unable to get transaction", txErr)
		return
	}

	txStmt := tx.Stmt(db.sqlAddMessage)

	_, txResultErr := txStmt.Exec(msgId, location, locationIndex, msgHash, msgFrom, msgData)
	if txResultErr != nil {
		err = fmt.Errorf(
			"%w: Unable to record message (id: %v, location: %v, location index: %v, hash: %v, msgFrom: '%s', data: %v)",
			txResultErr,
			msgId,
			location,
			locationIndex,
			msgHash,
			msgFrom,
			msgData,
		)
		tx.Rollback()
		return
	}

	txCommitErr := tx.Commit()
	if txCommitErr != nil {
		err = fmt.Errorf("%w: error while committing transaction", txCommitErr)
	}
	return
}

func (db *dbStorage) GetMessage(id string) (md MailData, err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	row := db.sqlGetMessageById.QueryRow(id)
	var (
		msg_id             string
		msg_location       string
		msg_location_index int
		msg_hash           []byte
		msg_from           string
		msg_data           []byte
	)
	rowErr := row.Scan(msg_id, msg_location, msg_location_index, msg_hash, msg_from, msg_data)
	if rowErr != nil {
		err = fmt.Errorf("%w: Unable to record message (id: %v)", rowErr, id)
		return
	}

	md = NewMessageData()
	msgIdErr := md.SetMessageId(msg_id)
	locErr := md.SetLocation(msg_location)
	locIdxErr := md.SetLocationIndex(msg_location_index)
	hashErr := md.SetHash(msg_hash)
	fromErr := md.SetFrom(msg_from)
	dataErr := md.SetData(msg_data)

	if msgIdErr != nil || locErr != nil || locIdxErr != nil || hashErr != nil || fromErr != nil || dataErr != nil {
		err = fmt.Errorf(
			"%w: Unable rebuild message (id: %v, location: %v, location index: %v, hash: %v, msgFrom: %v, data: %v)",
			common.ErrUnableToBuildMessageData,
			msgIdErr,
			locErr,
			locIdxErr,
			hashErr,
			fromErr,
			dataErr,
		)
		return
	}

	// all successful!
	return
}

func (db *dbStorage) rowsToMailDataArray(rows *sql.Rows) (msgs []MailData) {
	var (
		msg_id             string
		msg_location       string
		msg_location_index int
		msg_hash           []byte
		msg_from           string
		msg_data           []byte
	)
	for rows.Next() {
		rowErr := rows.Scan(msg_id, msg_location, msg_location_index, msg_hash, msg_from, msg_data)
		if rowErr != nil {
			log.Warning("%v: Unable to record message", rowErr)
			continue
		}

		md := NewMessageData()
		msgIdErr := md.SetMessageId(msg_id)
		locErr := md.SetLocation(msg_location)
		locIdxErr := md.SetLocationIndex(msg_location_index)
		hashErr := md.SetHash(msg_hash)
		fromErr := md.SetFrom(msg_from)
		dataErr := md.SetData(msg_data)

		if msgIdErr != nil || locErr != nil || locIdxErr != nil || hashErr != nil || fromErr != nil || dataErr != nil {
			log.Warning(
				"%v: Unable rebuild message (id: %v, location: %v, location index: %v, hash: %v, msgFrom: %v, data: %v)",
				common.ErrUnableToBuildMessageData,
				msgIdErr,
				locErr,
				locIdxErr,
				hashErr,
				fromErr,
				dataErr,
			)
			continue
		}

		msgs = append(msgs, md)
	}
	return
}

func (db *dbStorage) GetMessagesByLocation(location string) (msgs []MailData, err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	rows, rowsErr := db.sqlGetMessageByLocation.Query(location)
	if rowsErr != nil {
		err = fmt.Errorf("%w: Unable to retrieve records for location '%v'", rowsErr, location)
		return
	}

	msgs = db.rowsToMailDataArray(rows)
	if len(msgs) == 0 {
		err = fmt.Errorf("%w: No records found matching location '%v'", common.ErrNoMessagesFound, location)
		// in case there is something there, just clear it out
		msgs = nil
		return
	}

	// all successful!
	return
}

func (db *dbStorage) GetMessagesByHash(hash []byte) (msgs []MailData, err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	rows, rowsErr := db.sqlGetMessageByHash.Query(hash)
	if rowsErr != nil {
		err = fmt.Errorf("%w: Unable to retrieve records for hash '%v'", rowsErr, hash)
		return
	}

	msgs = db.rowsToMailDataArray(rows)
	if len(msgs) == 0 {
		err = fmt.Errorf("%w: No records found matching hash '%v'", common.ErrNoMessagesFound, hash)
		// in case there is something there, just clear it out
		msgs = nil
		return
	}

	// all successful!
	return
}

func (db *dbStorage) GetUniqueMessages() (msgs []MailData, err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	rows, rowsErr := db.sqlGetUniqueMessages.Query()
	if rowsErr != nil {
		err = fmt.Errorf("%w: Unable to retrieve unique records", rowsErr)
		return
	}

	msgs = db.rowsToMailDataArray(rows)
	if len(msgs) == 0 {
		err = fmt.Errorf("%w: No unique records found", common.ErrNoMessagesFound)
		// in case there is something there, just clear it out
		msgs = nil
		return
	}

	// all successful!
	return
}

/********************
 *** Mail Records ***
 ********************/
func (db *dbStorage) AddRecord(message MailRecord) (err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	tx, txErr := db.sqlDb.Begin()
	if txErr != nil {
		err = fmt.Errorf("%w: Unable to get transaction", txErr)
		return
	}

	txStmt := tx.Stmt(db.sqlAddRecord)

	_, txResultErr := txStmt.Exec(
        message.MsgId,
        message.Filename,
        message.Index,
        message.Hash,
    )
	if txResultErr != nil {
		err = fmt.Errorf(
			"%w: Unable to record message (id: %v, location: %v, location index: %v, hash: %v)",
			txResultErr,
			message.MsgId,
			message.Filename,
			message.Index,
			message.Hash,
		)
		tx.Rollback()
		return
	}

	txCommitErr := tx.Commit()
	if txCommitErr != nil {
		err = fmt.Errorf("%w: error while committing transaction", txCommitErr)
	}
	return
}

func (db *dbStorage) GetRecord(id string) (mr MailRecord, err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	row := db.sqlGetRecordById.QueryRow(id)
	var (
		msg_id             string
		msg_location       string
		msg_location_index int
		msg_hash           []byte
	)
	rowErr := row.Scan(msg_id, msg_location, msg_location_index, msg_hash)
	if rowErr != nil {
		err = fmt.Errorf("%w: Unable to record message (id: %v)", rowErr, id)
		return
	}

	mr = MailRecord{
        Index: msg_location_index,
        Hash: msg_hash,
        MsgId: msg_id,
        Filename: msg_location,
    }

	// all successful!
	return
}

func (db *dbStorage) rowsToMailRecordArray(rows *sql.Rows) (msgs []MailRecord) {
	var (
		msg_id             string
		msg_location       string
		msg_location_index int
		msg_hash           []byte
	)
	for rows.Next() {
		rowErr := rows.Scan(msg_id, msg_location, msg_location_index, msg_hash)
		if rowErr != nil {
			log.Warning("%v: Unable to record message", rowErr)
			continue
		}

        mr := MailRecord{
            Index: msg_location_index,
            Hash: msg_hash,
            MsgId: msg_id,
            Filename: msg_location,
        }

		msgs = append(msgs, mr)
	}
	return
}

func (db *dbStorage) GetRecordsByLocation(location string) (msgs []MailRecord, err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	rows, rowsErr := db.sqlGetRecordByLocation.Query(location)
	if rowsErr != nil {
		err = fmt.Errorf("%w: Unable to retrieve records for location '%v'", rowsErr, location)
		return
	}

	msgs = db.rowsToMailRecordArray(rows)
	if len(msgs) == 0 {
		err = fmt.Errorf("%w: No records found matching location '%v'", common.ErrNoMessagesFound, location)
		// in case there is something there, just clear it out
		msgs = nil
		return
	}

	// all successful!
	return
}

func (db *dbStorage) GetRecordsByHash(hash []byte) (msgs []MailRecord, err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	rows, rowsErr := db.sqlGetRecordByHash.Query(hash)
	if rowsErr != nil {
		err = fmt.Errorf("%w: Unable to retrieve records for hash '%v'", rowsErr, hash)
		return
	}

	msgs = db.rowsToMailRecordArray(rows)
	if len(msgs) == 0 {
		err = fmt.Errorf("%w: No records found matching hash '%v'", common.ErrNoMessagesFound, hash)
		// in case there is something there, just clear it out
		msgs = nil
		return
	}

	// all successful!
	return
}

func (db *dbStorage) GetUniqueRecords() (msgs []MailRecord, err error) {
	if db.sqlDb == nil {
		err = fmt.Errorf("%w: Database not available/open", common.ErrNoDatabaseAvailable)
		return
	}

	rows, rowsErr := db.sqlGetUniqueRecords.Query()
	if rowsErr != nil {
		err = fmt.Errorf("%w: Unable to retrieve unique records", rowsErr)
		return
	}

	msgs = db.rowsToMailRecordArray(rows)
	if len(msgs) == 0 {
		err = fmt.Errorf("%w: No unique records found", common.ErrNoMessagesFound)
		// in case there is something there, just clear it out
		msgs = nil
		return
	}

	// all successful!
	return
}

var _ MailStorage = &dbStorage{}
