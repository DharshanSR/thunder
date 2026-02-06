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
	"strings"

	"github.com/asgardeo/thunder/internal/system/database/transaction"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/log"
)

const (
	maxPreferenceKeyLength   = 255
	maxPreferenceValueLength = 10000 // Allow reasonably large preference values
)

// PreferenceServiceInterface defines the contract for preference business logic operations.
type PreferenceServiceInterface interface {
	GetPreferenceByKey(ctx context.Context, userID, key string) (*Preference, *serviceerror.ServiceError)
	GetPreferencesByUserID(ctx context.Context, userID string) ([]Preference, *serviceerror.ServiceError)
	UpsertPreferences(ctx context.Context, userID string, preferences map[string]string) ([]string, *serviceerror.ServiceError)
	DeletePreference(ctx context.Context, userID, key string) *serviceerror.ServiceError
}

// preferenceService implements the PreferenceServiceInterface.
type preferenceService struct {
	store        preferenceStoreInterface
	transactioner transaction.Transactioner
}

// newPreferenceService creates a new instance of preferenceService.
func newPreferenceService(store preferenceStoreInterface, transactioner transaction.Transactioner) PreferenceServiceInterface {
	return &preferenceService{
		store:        store,
		transactioner: transactioner,
	}
}

// GetPreferenceByKey retrieves a single preference by user ID and key.
func (ps *preferenceService) GetPreferenceByKey(ctx context.Context, userID, key string) (*Preference, *serviceerror.ServiceError) {
	logger := log.GetLogger().With(log.String("userID", userID), log.String("key", key))

	// Validate input
	if err := validatePreferenceKey(key); err != nil {
		return nil, err
	}

	preference, err := ps.store.GetPreferenceByKey(ctx, userID, key)
	if err != nil {
		if errors.Is(err, ErrPreferenceNotFound) {
			return nil, &ErrorPreferenceNotFound
		}
		logger.Error("Failed to get preference", log.Error(err))
		return nil, &ErrorInternalServerError
	}

	return preference, nil
}

// GetPreferencesByUserID retrieves all preferences for a user.
func (ps *preferenceService) GetPreferencesByUserID(ctx context.Context, userID string) ([]Preference, *serviceerror.ServiceError) {
	logger := log.GetLogger().With(log.String("userID", userID))

	preferences, err := ps.store.GetPreferencesByUserID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get preferences", log.Error(err))
		return nil, &ErrorInternalServerError
	}

	return preferences, nil
}

// UpsertPreferences creates or updates multiple preferences.
func (ps *preferenceService) UpsertPreferences(ctx context.Context, userID string, preferences map[string]string) ([]string, *serviceerror.ServiceError) {
	logger := log.GetLogger().With(log.String("userID", userID))

	// Validate all preferences before upserting
	for key, value := range preferences {
		if err := validatePreferenceKey(key); err != nil {
			return nil, err
		}
		if err := validatePreferenceValue(value); err != nil {
			return nil, err
		}
	}

	updatedKeys := make([]string, 0, len(preferences))

	// Use transaction to ensure atomicity
	err := ps.transactioner.Transact(ctx, func(txCtx context.Context) error {
		for key, value := range preferences {
			if err := ps.store.UpsertPreference(txCtx, userID, key, value); err != nil {
				return err
			}
			updatedKeys = append(updatedKeys, key)
		}
		return nil
	})

	if err != nil {
		logger.Error("Failed to upsert preferences", log.Error(err))
		return nil, &ErrorInternalServerError
	}

	return updatedKeys, nil
}

// DeletePreference deletes a preference by user ID and key.
func (ps *preferenceService) DeletePreference(ctx context.Context, userID, key string) *serviceerror.ServiceError {
	logger := log.GetLogger().With(log.String("userID", userID), log.String("key", key))

	// Validate input
	if err := validatePreferenceKey(key); err != nil {
		return err
	}

	err := ps.store.DeletePreference(ctx, userID, key)
	if err != nil {
		if errors.Is(err, ErrPreferenceNotFound) {
			return &ErrorPreferenceNotFound
		}
		logger.Error("Failed to delete preference", log.Error(err))
		return &ErrorInternalServerError
	}

	return nil
}

// validatePreferenceKey validates that a preference key is valid.
func validatePreferenceKey(key string) *serviceerror.ServiceError {
	if strings.TrimSpace(key) == "" {
		return &ErrorInvalidPreferenceKey
	}
	if len(key) > maxPreferenceKeyLength {
		return &ErrorInvalidPreferenceKey
	}
	return nil
}

// validatePreferenceValue validates that a preference value is valid.
func validatePreferenceValue(value string) *serviceerror.ServiceError {
	if len(value) > maxPreferenceValueLength {
		return &ErrorInvalidPreferenceValue
	}
	return nil
}
