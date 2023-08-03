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
package log

import (
	golog "log"
)

//go: cover ignore
func Info(format string, v ...any) {
	golog.Print(msgFormat(msgInfo, format, v...))
}

//go: cover ignore
func Warning(format string, v ...any) {
	golog.Print(msgFormat(msgWarning, format, v...))
}

//go: cover ignore
func Debug(format string, v ...any) {
	golog.Print(msgFormat(msgDebug, format, v...))
}

//go: cover ignore
func Notice(format string, v ...any) {
	golog.Print(msgFormat(msgNotice, format, v...))
}

//go: cover ignore
func Error(format string, v ...any) {
	golog.Print(msgFormat(msgError, format, v...))
}

//go: cover ignore
func Emergency(format string, v ...any) {
	golog.Print(msgFormat(msgEmergency, format, v...))
}

//go: cover ignore
func Critical(format string, v ...any) {
	golog.Print(msgFormat(msgCritical, format, v...))
}

//go: cover ignore
func Alert(format string, v ...any) {
	golog.Print(msgFormat(msgAlert, format, v...))
}
