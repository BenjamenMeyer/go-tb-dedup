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
