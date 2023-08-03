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

func TestLogMsgFormatters(t *testing.T) {
	t.Run(
		"format",
		func(t *testing.T) {
			type logParameters struct {
				level  string
				format string
				args   []any
			}
			type TestScenario struct {
				name       string
				parameters logParameters
				expected   string
			}
			generateScenario := func(name string, level string, format string, args []any) (scenario TestScenario) {
				scenario.name = name
				scenario.parameters = logParameters{
					level:  level,
					format: format,
					args:   args,
				}
				step1 := fmt.Sprintf(format, args...)
				step2 := fmt.Sprintf("%s: %s", level, step1)
				scenario.expected = step2
				return
			}

			testScenarios := []TestScenario{
				generateScenario("alpha", "world", "flush-%s", []any{"apples"}),
				generateScenario("beta", "hello", "a-%s", []any{"b"}),
			}

			for _, scenario := range testScenarios {
				t.Run(
					scenario.name,
					func(t *testing.T) {
						result := msgFormat(
							scenario.parameters.level,
							scenario.parameters.format,
							scenario.parameters.args...,
						)
						if result != scenario.expected {
							t.Errorf("Unexpected message generated: '%s' != '%s'", result, scenario.expected)
						}
					},
				)
			}
		},
	)
	t.Run(
		"unitFormat",
		func(t *testing.T) {
			type logParameters struct {
				unit   string
				level  string
				format string
				args   []any
			}
			type TestScenario struct {
				name       string
				parameters logParameters
				expected   string
			}
			generateScenario := func(name string, unit string, level string, format string, args []any) (scenario TestScenario) {
				scenario.name = name
				scenario.parameters = logParameters{
					unit:   unit,
					level:  level,
					format: format,
					args:   args,
				}
				step1 := fmt.Sprintf(format, args...)
				step2 := fmt.Sprintf("%s(%s): %s", level, unit, step1)
				scenario.expected = step2
				return
			}
			testScenarios := []TestScenario{
				generateScenario("one", "red", "blue", "fish-%s", []any{"one"}),
				generateScenario("two", "old", "new", "fish-%s", []any{"two"}),
			}

			for _, scenario := range testScenarios {
				t.Run(
					scenario.name,
					func(t *testing.T) {
						result := msgUnitFormat(
							scenario.parameters.unit,
							scenario.parameters.level,
							scenario.parameters.format,
							scenario.parameters.args...,
						)
						if result != scenario.expected {
							t.Errorf("Unexpected message generated: '%s' != '%s'", result, scenario.expected)
						}
					},
				)
			}
		},
	)
}
