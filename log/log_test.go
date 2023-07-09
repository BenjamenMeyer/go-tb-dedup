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
