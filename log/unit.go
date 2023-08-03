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

type outputInterface func(a ...any) // (n int, err error)

var logOutput outputInterface = golog.Print

func UnitInfo(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgInfo, format, v...))
}

func UnitWarning(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgWarning, format, v...))
}

func UnitDebug(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgDebug, format, v...))
}

func UnitNotice(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgNotice, format, v...))
}

func UnitError(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgError, format, v...))
}

func UnitEmergency(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgEmergency, format, v...))
}

func UnitCritical(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgCritical, format, v...))
}

func UnitAlert(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgAlert, format, v...))
}
