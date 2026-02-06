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
	"net/http"

	"github.com/asgardeo/thunder/internal/system/database/provider"
	"github.com/asgardeo/thunder/internal/system/middleware"
)

// Initialize initializes the preference service and registers HTTP routes.
func Initialize(mux *http.ServeMux) (PreferenceServiceInterface, error) {
	// Create store
	store, err := newPreferenceStore()
	if err != nil {
		return nil, err
	}

	// Get transactioner from provider
	transactioner, err := provider.GetDBProvider().GetUserDBTransactioner()
	if err != nil {
		return nil, err
	}

	// Create service
	service := newPreferenceService(store, transactioner)

	// Create handler
	handler := newPreferenceHandler(service)

	// Register routes with CORS middleware
	registerRoutes(mux, handler)

	return service, nil
}

// registerRoutes registers HTTP routes for preference operations.
func registerRoutes(mux *http.ServeMux, handler *preferenceHandler) {
	// GET /users/me/preferences - Get all preferences
	mux.Handle("GET /users/me/preferences", middleware.WithCORS(http.HandlerFunc(handler.handleGetPreferences)))

	// GET /users/me/preferences/{key} - Get a specific preference
	mux.Handle("GET /users/me/preferences/{key...}", middleware.WithCORS(http.HandlerFunc(handler.handleGetPreferenceByKey)))

	// PUT /users/me/preferences - Create or update preferences
	mux.Handle("PUT /users/me/preferences", middleware.WithCORS(http.HandlerFunc(handler.handleUpsertPreferences)))

	// DELETE /users/me/preferences/{key} - Delete a preference
	mux.Handle("DELETE /users/me/preferences/{key...}", middleware.WithCORS(http.HandlerFunc(handler.handleDeletePreference)))
}
