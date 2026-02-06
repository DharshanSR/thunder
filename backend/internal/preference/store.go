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
	"context"
	"errors"
	"fmt"

	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/database/provider"
)

var (
	// ErrPreferenceNotFound is returned when a preference is not found.
	ErrPreferenceNotFound = errors.New("preference not found")
)

// preferenceStoreInterface defines the contract for preference data access operations.
type preferenceStoreInterface interface {
	GetPreferenceByKey(ctx context.Context, userID, key string) (*Preference, error)
	GetPreferencesByUserID(ctx context.Context, userID string) ([]Preference, error)
	UpsertPreference(ctx context.Context, userID, key, value string) error
	DeletePreference(ctx context.Context, userID, key string) error
}

// preferenceStore implements the preferenceStoreInterface.
type preferenceStore struct {
	dbClient     provider.DBClientInterface
	deploymentID string
}

// newPreferenceStore creates a new instance of preferenceStore.
func newPreferenceStore() (*preferenceStore, error) {
	runtime := config.GetThunderRuntime()
	dbClient, err := provider.GetDBProvider().GetUserDBClient()
	if err != nil {
		return nil, err
	}

	return &preferenceStore{
		dbClient:     dbClient,
		deploymentID: runtime.Config.Server.Identifier,
	}, nil
}

// GetPreferenceByKey retrieves a single preference by user ID and key.
func (ps *preferenceStore) GetPreferenceByKey(ctx context.Context, userID, key string) (*Preference, error) {
	results, err := ps.dbClient.QueryContext(ctx, queryGetPreferenceByKey, userID, key, ps.deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if len(results) == 0 {
		return nil, ErrPreferenceNotFound
	}

	row := results[0]
	preference := &Preference{
		Key:   toString(row["preference_key"]),
		Value: toString(row["preference_value"]),
	}

	// Handle timestamps which might be returned as different types
	if createdAt := row["created_at"]; createdAt != nil {
		preference.CreatedAt = toString(createdAt)
	}
	if updatedAt := row["updated_at"]; updatedAt != nil {
		preference.UpdatedAt = toString(updatedAt)
	}

	return preference, nil
}

// GetPreferencesByUserID retrieves all preferences for a user.
func (ps *preferenceStore) GetPreferencesByUserID(ctx context.Context, userID string) ([]Preference, error) {
	results, err := ps.dbClient.QueryContext(ctx, queryGetPreferencesByUserID, userID, ps.deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	preferences := make([]Preference, 0, len(results))
	for _, row := range results {
		preference := Preference{
			Key:   toString(row["preference_key"]),
			Value: toString(row["preference_value"]),
		}
		
		// Handle timestamps which might be returned as different types
		if createdAt := row["created_at"]; createdAt != nil {
			preference.CreatedAt = toString(createdAt)
		}
		if updatedAt := row["updated_at"]; updatedAt != nil {
			preference.UpdatedAt = toString(updatedAt)
		}
		
		preferences = append(preferences, preference)
	}

	return preferences, nil
}

// UpsertPreference creates or updates a preference.
func (ps *preferenceStore) UpsertPreference(ctx context.Context, userID, key, value string) error {
	_, err := ps.dbClient.ExecuteContext(ctx, queryUpsertPreference, userID, key, value, ps.deploymentID)
	if err != nil {
		return fmt.Errorf("failed to upsert preference: %w", err)
	}
	return nil
}

// DeletePreference deletes a preference by user ID and key.
func (ps *preferenceStore) DeletePreference(ctx context.Context, userID, key string) error {
	rowsAffected, err := ps.dbClient.ExecuteContext(ctx, queryDeletePreference, userID, key, ps.deploymentID)
	if err != nil {
		return fmt.Errorf("failed to delete preference: %w", err)
	}

	if rowsAffected == 0 {
		return ErrPreferenceNotFound
	}

	return nil
}

// toString safely converts an interface{} to a string.
func toString(val interface{}) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}
