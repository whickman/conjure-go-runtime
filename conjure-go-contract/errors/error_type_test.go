// Copyright (c) 2018 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"fmt"
	"testing"

	"github.com/palantir/conjure-go-runtime/conjure-go-contract/errors"
	"github.com/stretchr/testify/assert"
)

var invalidErrorNames = map[string]string{
	"no namespace":                       ":PascalCase",
	"no cause":                           "PascalCase:",
	"no namespace nor cause":             ":",
	"no colon":                           "PascalCase",
	"namespace with invalid case":        "notPascalCase:PascalCase",
	"cause with invalid case":            "PascalCase:notPascalCase",
	"name with three parts":              "PascalCase:PascalCase:PascalCase",
	"default namespace and custom cause": "Default:CustomError",
}

var customValidErrorNames = []string{
	"MyApplication:MyCustomError",
	"Aa:Bb",
	"A1:B2",
	"A1A1:B2B2",
	"Aa1Aa1:Bb2Bb2",
	"A1aA1a:B2bB2b",
	"MyApplication:MyCustomError",
}

var defaultErrorTypes = []errors.ErrorType{
	errors.DefaultPermissionDenied,
	errors.DefaultInvalidArgument,
	errors.DefaultNotFound,
	errors.DefaultConflict,
	errors.DefaultRequestEntityTooLarge,
	errors.DefaultFailedPrecondition,
	errors.DefaultInternal,
	errors.DefaultTimeout,
}

func TestNewErrorType_ForDefaultErrorNames(t *testing.T) {
	for _, errorCode := range validErrorCodes {
		t.Run(errorCode.String(), func(t *testing.T) {
			for i, defaultErrorType := range defaultErrorTypes {
				defaultErrorName := defaultErrorType.Name()
				assignedErrorCode := defaultErrorType.Code()
				t.Run(fmt.Sprintf("default_error_name_%d", i), func(t *testing.T) {
					_, err := errors.NewErrorType(errorCode, defaultErrorName)
					if assignedErrorCode == errorCode {
						assert.NoError(t, err)
					} else {
						assert.EqualError(t, err, "errors: invalid combination of default error name and error code")
					}
				})
			}
		})
	}
}

func TestNewErrorType_ForDefaultErrorTypeSucceeds(t *testing.T) {
	for _, defaultErrorType := range defaultErrorTypes {
		t.Run(defaultErrorType.String(), func(t *testing.T) {
			_, err := errors.NewErrorType(
				defaultErrorType.Code(),
				defaultErrorType.Name(),
			)
			assert.NoError(t, err)
		})
	}
}

func TestNewErrorType_ForCustomValidErrorName(t *testing.T) {
	for _, errorCode := range validErrorCodes {
		t.Run(errorCode.String(), func(t *testing.T) {
			for i, errorName := range customValidErrorNames {
				t.Run(fmt.Sprintf("valid_name_%d", i), func(t *testing.T) {
					_, err := errors.NewErrorType(errorCode, errorName)
					assert.NoError(t, err)
				})
			}
		})
	}
}

func TestNewErrorType_ForInvalidErrorName(t *testing.T) {
	for _, errorCode := range validErrorCodes {
		t.Run(errorCode.String(), func(t *testing.T) {
			for caseName, errorName := range invalidErrorNames {
				t.Run(caseName, func(t *testing.T) {
					_, err := errors.NewErrorType(errorCode, errorName)
					switch caseName {
					case "default namespace and custom cause":
						assert.EqualError(t, err, "errors: error name with default namespace cannot use custom cause")
					default:
						assert.EqualError(t, err, "errors: error name does not match regexp `^(([A-Z][a-z0-9]+)+):(([A-Z][a-z0-9]+)+)$`")
					}
				})
			}
		})
	}
}
