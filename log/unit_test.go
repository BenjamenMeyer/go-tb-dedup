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
	"fmt"
	"testing"
)

func TestUnitInterface(t *testing.T) {
	t.Run(
		"format",
		func(t *testing.T) {
			type logParameters struct {
				unit   string
				format string
				args   []any
			}
			type TestScenario struct {
				name       string
				unitFunc   func(unit string, format string, v ...any)
				parameters logParameters
				expected   string
			}
			generateScenario := func(name string, level string, fn func(unit string, format string, v ...any), unit string, format string, args []any) (scenario TestScenario) {
				scenario.name = name
				scenario.unitFunc = fn
				scenario.parameters = logParameters{
					unit:   unit,
					format: format,
					args:   args,
				}
				scenario.expected = msgUnitFormat(unit, level, format, args...)
				return
			}

			testScenarios := []TestScenario{
				generateScenario("info", msgInfo, UnitInfo, "world", "flush-%s", []any{"apples"}),
				generateScenario("warning", msgWarning, UnitWarning, "hello", "a-%s", []any{"b"}),
				generateScenario("debug", msgDebug, UnitDebug, "foo", "oof-%s", []any{"b33f"}),
				generateScenario("notice", msgNotice, UnitNotice, "bar", "rab-%s", []any{"d34d"}),
				generateScenario("error", msgError, UnitError, "b33f", "foo-%s", []any{"oof"}),
				generateScenario("emergency", msgEmergency, UnitEmergency, "d34d", "bar", []any{"rab"}),
				generateScenario("critical", msgCritical, UnitCritical, "oof", "b33f-%s", []any{"foo"}),
				generateScenario("alert", msgAlert, UnitAlert, "rab", "d34d", []any{"bar"}),
			}

			for _, scenario := range testScenarios {
				t.Run(
					scenario.name,
					func(t *testing.T) {
						var result string
						orig := logOutput
						logOutput = func(a ...any) {
							result = fmt.Sprint(a...)
							return
						}
						scenario.unitFunc(
							scenario.parameters.unit,
							scenario.parameters.format,
							scenario.parameters.args...,
						)
						if result != scenario.expected {
							t.Errorf("Unexpected message generated: '%s' != '%s'", result, scenario.expected)
						}

						logOutput = orig
					},
				)
			}
		},
	)
}
