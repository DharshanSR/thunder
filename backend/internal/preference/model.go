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

import "time"

// Preference represents a user preference.
type Preference struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetPreferenceResponse represents the response for getting a single preference.
type GetPreferenceResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetPreferencesResponse represents the response for getting all preferences.
type GetPreferencesResponse struct {
	Preferences []Preference `json:"preferences"`
}

// UpsertPreferencesRequest represents the request to create or update preferences.
type UpsertPreferencesRequest struct {
	Preferences map[string]string `json:"preferences"`
}

// UpsertPreferencesResponse represents the response for upserting preferences.
type UpsertPreferencesResponse struct {
	UpdatedKeys []string `json:"updated_keys"`
}

// DeletePreferenceResponse represents the response for deleting a preference.
type DeletePreferenceResponse struct {
	Message string `json:"message"`
}
