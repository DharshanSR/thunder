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
	"encoding/json"
	"net/http"
	"strings"

	apierror "github.com/asgardeo/thunder/internal/system/error/apierror"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/log"
	"github.com/asgardeo/thunder/internal/system/security"
)

// preferenceHandler handles HTTP requests for preference operations.
type preferenceHandler struct {
	service PreferenceServiceInterface
}

// newPreferenceHandler creates a new instance of preferenceHandler.
func newPreferenceHandler(service PreferenceServiceInterface) *preferenceHandler {
	return &preferenceHandler{
		service: service,
	}
}

// handleGetPreferences handles GET /users/me/preferences - retrieves all preferences for the authenticated user.
func (ph *preferenceHandler) handleGetPreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.GetLogger()

	// Extract user ID from authentication context
	userID := security.GetUserID(ctx)
	if strings.TrimSpace(userID) == "" {
		handleError(w, &ErrorAuthenticationFailed)
		return
	}

	preferences, svcErr := ph.service.GetPreferencesByUserID(ctx, userID)
	if svcErr != nil {
		handleError(w, svcErr)
		return
	}

	response := GetPreferencesResponse{
		Preferences: preferences,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", log.Error(err))
	}
}

// handleGetPreferenceByKey handles GET /users/me/preferences/{key} - retrieves a specific preference.
func (ph *preferenceHandler) handleGetPreferenceByKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.GetLogger()

	// Extract user ID from authentication context
	userID := security.GetUserID(ctx)
	if strings.TrimSpace(userID) == "" {
		handleError(w, &ErrorAuthenticationFailed)
		return
	}

	// Extract preference key from path
	key := strings.TrimPrefix(r.URL.Path, "/users/me/preferences/")
	if strings.TrimSpace(key) == "" {
		handleError(w, &ErrorInvalidRequest)
		return
	}

	preference, svcErr := ph.service.GetPreferenceByKey(ctx, userID, key)
	if svcErr != nil {
		handleError(w, svcErr)
		return
	}

	response := GetPreferenceResponse{
		Key:   preference.Key,
		Value: preference.Value,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", log.Error(err))
	}
}

// handleUpsertPreferences handles PUT /users/me/preferences - creates or updates preferences.
func (ph *preferenceHandler) handleUpsertPreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.GetLogger()

	// Extract user ID from authentication context
	userID := security.GetUserID(ctx)
	if strings.TrimSpace(userID) == "" {
		handleError(w, &ErrorAuthenticationFailed)
		return
	}

	// Parse request body
	var req UpsertPreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request", log.Error(err))
		handleError(w, &ErrorInvalidRequest)
		return
	}

	if req.Preferences == nil || len(req.Preferences) == 0 {
		handleError(w, &ErrorInvalidRequest)
		return
	}

	updatedKeys, svcErr := ph.service.UpsertPreferences(ctx, userID, req.Preferences)
	if svcErr != nil {
		handleError(w, svcErr)
		return
	}

	response := UpsertPreferencesResponse{
		UpdatedKeys: updatedKeys,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", log.Error(err))
	}
}

// handleDeletePreference handles DELETE /users/me/preferences/{key} - deletes a specific preference.
func (ph *preferenceHandler) handleDeletePreference(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.GetLogger()

	// Extract user ID from authentication context
	userID := security.GetUserID(ctx)
	if strings.TrimSpace(userID) == "" {
		handleError(w, &ErrorAuthenticationFailed)
		return
	}

	// Extract preference key from path
	key := strings.TrimPrefix(r.URL.Path, "/users/me/preferences/")
	if strings.TrimSpace(key) == "" {
		handleError(w, &ErrorInvalidRequest)
		return
	}

	svcErr := ph.service.DeletePreference(ctx, userID, key)
	if svcErr != nil {
		handleError(w, svcErr)
		return
	}

	response := DeletePreferenceResponse{
		Message: "Preference deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", log.Error(err))
	}
}

// handleError writes an error response based on service error.
func handleError(w http.ResponseWriter, svcErr *serviceerror.ServiceError) {
	var statusCode int
	if svcErr.Type == serviceerror.ClientErrorType {
		switch svcErr.Code {
		case ErrorPreferenceNotFound.Code:
			statusCode = http.StatusNotFound
		case ErrorAuthenticationFailed.Code:
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusBadRequest
		}
	} else {
		statusCode = http.StatusInternalServerError
	}

	errResp := apierror.ErrorResponse{
		Code:        svcErr.Code,
		Message:     svcErr.Error,
		Description: svcErr.ErrorDescription,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errResp)
}
