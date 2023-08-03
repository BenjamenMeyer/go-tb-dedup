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
	"fmt"
    "io/fs"
	"os"
	"path"

	"github.com/BenjamenMeyer/go-tb-dedup/common"
	"github.com/BenjamenMeyer/go-tb-dedup/log"
)

type folderData struct {
	folder     Folderpath
	folderData []os.DirEntry

	files   []Filepath
	folders []Folderpath
}

func NewMailFolder() Mailfolder {
	return &folderData{}
}

func (fd *folderData) isOpened() (result bool) {
	if len(fd.folder) > 0 {
		result = true
	}
	return
}

func (fd *folderData) Open(folder Folderpath) (err error) {
	if fd.isOpened() {
		err = fmt.Errorf("%w: Already working on %s", common.ErrAlreadyOpen, fd.folder)
		return
	}

	if len(fd.folderData) > 0 {
		err = fmt.Errorf("%w: data set still exists from previous operation", common.ErrBadState)
		return
	}

	fd.folder = folder
	fd.folderData, err = os.ReadDir(string(fd.folder))
	if err != nil {
		fmt.Errorf("%w: failed to read directory data for location %s", err, fd.folder)
		return
	}

	fd.files = make([]Filepath, 0)
	fd.folders = make([]Folderpath, 0)
	for _, fileData := range fd.folderData {
		fileMode := fileData.Type()

		fullPath := path.Join(string(fd.folder), fileData.Name())

		// save sub-directories off
		if fileMode.IsDir() {
			fd.folders = append(fd.folders, Folderpath(fullPath))
			continue
		}

		// save files off
		// skip any symlinks, etc; only look at local, regular files
		//if fileMode.IsRegular() {
		//if fileMode.IsFile() {
        fileType := fileMode.Type()
        supportedModes := make(map[fs.FileMode]string)
        supportedModes[0] = "Regular"
        supportedModes[fs.ModeSymlink] = "Symlink"

        if _, ok := supportedModes[fileType]; ok { 
			fd.files = append(fd.files, Filepath(fullPath))
			continue
		}

		// log anything else just so we have a record of it
		log.Warning("Skipping %s - File Type %s", fileData.Name(), fileMode.String())
	}

	return
}

func (fd *folderData) Close() (err error) {
	// reset the data
	fd.folderData = nil
	fd.files = nil
	fd.folders = nil
	fd.folder = ""
	return
}

func (fd *folderData) GetCounts() (fileCount int, folderCount int, err error) {
	if !fd.isOpened() {
		err = fmt.Errorf("%w: No data to process", common.ErrNotOpened)
		return
	}

	fileCount = len(fd.files)
	folderCount = len(fd.folders)
	return
}

func (fd *folderData) GetMailFiles() (files []Filepath, err error) {
	if !fd.isOpened() {
		err = fmt.Errorf("%w: No data to process", common.ErrNotOpened)
		return
	}

	if fd.files == nil {
		err = fmt.Errorf("%w: no files found at %s", common.ErrNoFiles, fd.folder)
		return
	}

	files = fd.files
	return
}

func (fd *folderData) GetSubfolders() (folders []Folderpath, err error) {
	if !fd.isOpened() {
		err = fmt.Errorf("%w: No data to process", common.ErrNotOpened)
		return
	}

	if fd.folders == nil {
		err = fmt.Errorf("%w: no sub-folders found at %s", common.ErrNoFolders, fd.folder)
		return
	}

	folders = fd.folders
	return
}

var _ Mailfolder = &folderData{}
