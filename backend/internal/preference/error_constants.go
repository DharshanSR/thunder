/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package preference

import (
	apierror "github.com/asgardeo/thunder/internal/system/error/apierror"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
)

// Service error constants.
var (
	ErrorPreferenceNotFound = serviceerror.ServiceError{
		Code:  "PREF-4001",
		Type:  serviceerror.ClientErrorType,
		Error: "Preference not found",
	}
	ErrorInvalidPreferenceKey = serviceerror.ServiceError{
		Code:  "PREF-4002",
		Type:  serviceerror.ClientErrorType,
		Error: "Invalid preference key",
	}
	ErrorInvalidPreferenceValue = serviceerror.ServiceError{
		Code:  "PREF-4003",
		Type:  serviceerror.ClientErrorType,
		Error: "Invalid preference value",
	}
	ErrorInternalServerError = serviceerror.ServiceError{
		Code:  "PREF-5000",
		Type:  serviceerror.ServerErrorType,
		Error: "Internal server error",
	}
)

// API error constants.
var (
	ErrorAuthenticationFailed = apierror.ErrorResponse{
		Code:    "PREF-4000",
		Message: "Authentication failed",
		Status:  401,
	}
	ErrorPreferenceNotFoundAPI = apierror.ErrorResponse{
		Code:    "PREF-4001",
		Message: "Preference not found",
		Status:  404,
	}
	ErrorInvalidRequest = apierror.ErrorResponse{
		Code:    "PREF-4002",
		Message: "Invalid request",
		Status:  400,
	}
	ErrorInternalServerErrorAPI = apierror.ErrorResponse{
		Code:    "PREF-5000",
		Message: "Internal server error",
		Status:  500,
	}
)
