package nt

import (
	"fmt"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewNotificationTemplate(t *testing.T) {
	aTest := tester.New(t)
	var nt *NotificationTemplate
	var err error

	type TestCase struct {
		templateName       string
		args               []any
		IsErrorExpected    bool
		ExpectedComponents []Component
	}
	var tests = []TestCase{
		{
			templateName:    "",
			args:            []any{},
			IsErrorExpected: true,
		},
		{
			templateName:    "R",
			args:            []any{1, 2, 3, 4, 5, 6},
			IsErrorExpected: true,
		},
		{
			templateName:    "X",
			args:            []any{1},
			IsErrorExpected: true,
		},
		{
			templateName:    "UUU",
			args:            []any{1, 2, 3},
			IsErrorExpected: true,
		},
		{
			templateName:    "R",
			args:            []any{"text"},
			IsErrorExpected: true,
		},
		{
			templateName:    "F",
			args:            []any{123},
			IsErrorExpected: true,
		},
		{
			templateName:    "TF",
			args:            []any{5005, "Thread: {T}."},
			IsErrorExpected: true,
		},
		{
			templateName:    "FT",
			args:            []any{"{T}-{T}", 5005},
			IsErrorExpected: true,
		},
		{
			templateName:    "FTU",
			args:            []any{"Thread: {T}, User: {X}.", 5005, 1001},
			IsErrorExpected: true,
		},
		{
			templateName:    "FTR",
			args:            []any{"Thread: {T}, Not Record: {U}.", 5005, 1001},
			IsErrorExpected: true,
		},
		{
			templateName:    "FTU",
			args:            []any{"Thread: {T}, User: {U}.", 5005, 1001},
			IsErrorExpected: false,
			ExpectedComponents: []Component{
				{
					Name:  "F",
					Value: "Thread: {T}, User: {U}.",
				},
				{
					Name:  "T",
					Value: 5005,
				},
				{
					Name:  "U",
					Value: 1001,
				},
			},
		},
	}

	for i, test := range tests {
		fmt.Print(fmt.Sprintf("[%d]", i+1))
		nt, err = NewNotificationTemplate(test.templateName, test.args)
		if test.IsErrorExpected {
			aTest.MustBeAnError(err)
		} else {
			aTest.MustBeNoError(err)

			aTest.MustBeEqual(nt.components, test.ExpectedComponents)
		}
	}
	fmt.Println()
}
