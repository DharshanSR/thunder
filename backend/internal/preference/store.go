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
	"database/sql"
	"errors"

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
	dbClient provider.DBClientInterface
}

// newPreferenceStore creates a new instance of preferenceStore.
func newPreferenceStore() (*preferenceStore, error) {
	dbClient, err := provider.GetDBProvider().GetUserDBClient()
	if err != nil {
		return nil, err
	}

	return &preferenceStore{
		dbClient: dbClient,
	}, nil
}

// GetPreferenceByKey retrieves a single preference by user ID and key.
func (ps *preferenceStore) GetPreferenceByKey(ctx context.Context, userID, key string) (*Preference, error) {
	query, err := ps.dbClient.GetQuery(queryGetPreferenceByKey)
	if err != nil {
		return nil, err
	}

	deploymentID := config.GetThunderRuntime().GetDeploymentID()

	var preference Preference
	err = ps.dbClient.QueryRowContext(ctx, query, userID, key, deploymentID).Scan(
		&preference.Key,
		&preference.Value,
		&preference.CreatedAt,
		&preference.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPreferenceNotFound
		}
		return nil, err
	}

	return &preference, nil
}

// GetPreferencesByUserID retrieves all preferences for a user.
func (ps *preferenceStore) GetPreferencesByUserID(ctx context.Context, userID string) ([]Preference, error) {
	query, err := ps.dbClient.GetQuery(queryGetPreferencesByUserID)
	if err != nil {
		return nil, err
	}

	deploymentID := config.GetThunderRuntime().GetDeploymentID()

	rows, err := ps.dbClient.QueryContext(ctx, query, userID, deploymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var preferences []Preference
	for rows.Next() {
		var preference Preference
		err := rows.Scan(
			&preference.Key,
			&preference.Value,
			&preference.CreatedAt,
			&preference.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		preferences = append(preferences, preference)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for consistency
	if preferences == nil {
		preferences = []Preference{}
	}

	return preferences, nil
}

// UpsertPreference creates or updates a preference.
func (ps *preferenceStore) UpsertPreference(ctx context.Context, userID, key, value string) error {
	query, err := ps.dbClient.GetQuery(queryUpsertPreference)
	if err != nil {
		return err
	}

	deploymentID := config.GetThunderRuntime().GetDeploymentID()

	_, err = ps.dbClient.ExecContext(ctx, query, userID, key, value, deploymentID)
	return err
}

// DeletePreference deletes a preference by user ID and key.
func (ps *preferenceStore) DeletePreference(ctx context.Context, userID, key string) error {
	query, err := ps.dbClient.GetQuery(queryDeletePreference)
	if err != nil {
		return err
	}

	deploymentID := config.GetThunderRuntime().GetDeploymentID()

	result, err := ps.dbClient.ExecContext(ctx, query, userID, key, deploymentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPreferenceNotFound
	}

	return nil
}
