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
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
)

// Service error constants.
var (
	ErrorPreferenceNotFound = serviceerror.ServiceError{
		Code:             "PREF-4001",
		Type:             serviceerror.ClientErrorType,
		Error:            "Preference not found",
		ErrorDescription: "The requested preference does not exist",
	}
	ErrorInvalidPreferenceKey = serviceerror.ServiceError{
		Code:             "PREF-4002",
		Type:             serviceerror.ClientErrorType,
		Error:            "Invalid preference key",
		ErrorDescription: "The preference key is invalid or exceeds maximum length",
	}
	ErrorInvalidPreferenceValue = serviceerror.ServiceError{
		Code:             "PREF-4003",
		Type:             serviceerror.ClientErrorType,
		Error:            "Invalid preference value",
		ErrorDescription: "The preference value exceeds maximum length",
	}
	ErrorAuthenticationFailed = serviceerror.ServiceError{
		Code:             "PREF-4000",
		Type:             serviceerror.ClientErrorType,
		Error:            "Authentication failed",
		ErrorDescription: "User authentication is required",
	}
	ErrorInvalidRequest = serviceerror.ServiceError{
		Code:             "PREF-4004",
		Type:             serviceerror.ClientErrorType,
		Error:            "Invalid request",
		ErrorDescription: "The request is invalid or missing required fields",
	}
	ErrorInternalServerError = serviceerror.ServiceError{
		Code:             "PREF-5000",
		Type:             serviceerror.ServerErrorType,
		Error:            "Internal server error",
		ErrorDescription: "An unexpected error occurred while processing the request",
	}
)
