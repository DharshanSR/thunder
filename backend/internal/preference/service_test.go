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
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/asgardeo/thunder/internal/system/database/transaction"
)

type PreferenceServiceTestSuite struct {
	suite.Suite
	service       PreferenceServiceInterface
	mockStore     *preferenceStoreInterfaceMock
	transactioner *transactionerMock
}

func (suite *PreferenceServiceTestSuite) SetupTest() {
	suite.mockStore = newPreferenceStoreInterfaceMock(suite.T())
	suite.transactioner = newTransactionerMock(suite.T())
	suite.service = newPreferenceService(suite.mockStore, suite.transactioner)
}

func (suite *PreferenceServiceTestSuite) TestGetPreferenceByKey_Success() {
	ctx := context.Background()
	expectedPreference := &Preference{
		Key:   testKey,
		Value: testValue,
	}

	suite.mockStore.On("GetPreferenceByKey", ctx, testUserID, testKey).Return(expectedPreference, nil)

	// Execute
	preference, err := suite.service.GetPreferenceByKey(ctx, testUserID, testKey)

	// Assert
	suite.NoError(err)
	suite.NotNil(preference)
	suite.Equal(expectedPreference.Key, preference.Key)
	suite.Equal(expectedPreference.Value, preference.Value)
	suite.mockStore.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestGetPreferenceByKey_NotFound() {
	ctx := context.Background()

	suite.mockStore.On("GetPreferenceByKey", ctx, testUserID, testKey).Return(nil, ErrPreferenceNotFound)

	// Execute
	preference, err := suite.service.GetPreferenceByKey(ctx, testUserID, testKey)

	// Assert
	suite.Error(err)
	suite.Equal(ErrorPreferenceNotFound.Code, err.Code)
	suite.Nil(preference)
	suite.mockStore.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestGetPreferenceByKey_InvalidKey() {
	ctx := context.Background()

	// Test with empty key
	preference, err := suite.service.GetPreferenceByKey(ctx, testUserID, "")

	// Assert
	suite.Error(err)
	suite.Equal(ErrorInvalidPreferenceKey.Code, err.Code)
	suite.Nil(preference)

	// Test with key that's too long
	longKey := strings.Repeat("a", maxPreferenceKeyLength+1)
	preference, err = suite.service.GetPreferenceByKey(ctx, testUserID, longKey)

	// Assert
	suite.Error(err)
	suite.Equal(ErrorInvalidPreferenceKey.Code, err.Code)
	suite.Nil(preference)
}

func (suite *PreferenceServiceTestSuite) TestGetPreferenceByKey_StoreError() {
	ctx := context.Background()

	suite.mockStore.On("GetPreferenceByKey", ctx, testUserID, testKey).Return(nil, errors.New("database error"))

	// Execute
	preference, err := suite.service.GetPreferenceByKey(ctx, testUserID, testKey)

	// Assert
	suite.Error(err)
	suite.Equal(ErrorInternalServerError.Code, err.Code)
	suite.Nil(preference)
	suite.mockStore.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestGetPreferencesByUserID_Success() {
	ctx := context.Background()
	expectedPreferences := []Preference{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}

	suite.mockStore.On("GetPreferencesByUserID", ctx, testUserID).Return(expectedPreferences, nil)

	// Execute
	preferences, err := suite.service.GetPreferencesByUserID(ctx, testUserID)

	// Assert
	suite.NoError(err)
	suite.NotNil(preferences)
	suite.Len(preferences, 2)
	suite.mockStore.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestGetPreferencesByUserID_Empty() {
	ctx := context.Background()

	suite.mockStore.On("GetPreferencesByUserID", ctx, testUserID).Return([]Preference{}, nil)

	// Execute
	preferences, err := suite.service.GetPreferencesByUserID(ctx, testUserID)

	// Assert
	suite.NoError(err)
	suite.NotNil(preferences)
	suite.Len(preferences, 0)
	suite.mockStore.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestGetPreferencesByUserID_StoreError() {
	ctx := context.Background()

	suite.mockStore.On("GetPreferencesByUserID", ctx, testUserID).Return(nil, errors.New("database error"))

	// Execute
	preferences, err := suite.service.GetPreferencesByUserID(ctx, testUserID)

	// Assert
	suite.Error(err)
	suite.Equal(ErrorInternalServerError.Code, err.Code)
	suite.Nil(preferences)
	suite.mockStore.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestUpsertPreferences_Success() {
	ctx := context.Background()
	preferences := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	// Setup transactioner to execute the transaction function
	suite.transactioner.On("Transact", ctx, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			txFunc := args.Get(1).(func(context.Context) error)
			_ = txFunc(ctx)
		}).
		Return(nil)

	suite.mockStore.On("UpsertPreference", ctx, testUserID, "key1", "value1").Return(nil)
	suite.mockStore.On("UpsertPreference", ctx, testUserID, "key2", "value2").Return(nil)

	// Execute
	updatedKeys, err := suite.service.UpsertPreferences(ctx, testUserID, preferences)

	// Assert
	suite.NoError(err)
	suite.NotNil(updatedKeys)
	suite.Len(updatedKeys, 2)
	suite.mockStore.AssertExpectations(suite.T())
	suite.transactioner.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestUpsertPreferences_InvalidKey() {
	ctx := context.Background()
	preferences := map[string]string{
		"": "value1",
	}

	// Execute
	updatedKeys, err := suite.service.UpsertPreferences(ctx, testUserID, preferences)

	// Assert
	suite.Error(err)
	suite.Equal(ErrorInvalidPreferenceKey.Code, err.Code)
	suite.Nil(updatedKeys)
}

func (suite *PreferenceServiceTestSuite) TestUpsertPreferences_InvalidValue() {
	ctx := context.Background()
	longValue := strings.Repeat("a", maxPreferenceValueLength+1)
	preferences := map[string]string{
		"key1": longValue,
	}

	// Execute
	updatedKeys, err := suite.service.UpsertPreferences(ctx, testUserID, preferences)

	// Assert
	suite.Error(err)
	suite.Equal(ErrorInvalidPreferenceValue.Code, err.Code)
	suite.Nil(updatedKeys)
}

func (suite *PreferenceServiceTestSuite) TestUpsertPreferences_TransactionError() {
	ctx := context.Background()
	preferences := map[string]string{
		"key1": "value1",
	}

	// Setup transactioner to return error
	suite.transactioner.On("Transact", ctx, mock.AnythingOfType("func(context.Context) error")).
		Return(errors.New("transaction error"))

	// Execute
	updatedKeys, err := suite.service.UpsertPreferences(ctx, testUserID, preferences)

	// Assert
	suite.Error(err)
	suite.Equal(ErrorInternalServerError.Code, err.Code)
	suite.Nil(updatedKeys)
	suite.transactioner.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestDeletePreference_Success() {
	ctx := context.Background()

	suite.mockStore.On("DeletePreference", ctx, testUserID, testKey).Return(nil)

	// Execute
	err := suite.service.DeletePreference(ctx, testUserID, testKey)

	// Assert
	suite.NoError(err)
	suite.mockStore.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestDeletePreference_NotFound() {
	ctx := context.Background()

	suite.mockStore.On("DeletePreference", ctx, testUserID, testKey).Return(ErrPreferenceNotFound)

	// Execute
	err := suite.service.DeletePreference(ctx, testUserID, testKey)

	// Assert
	suite.Error(err)
	suite.Equal(ErrorPreferenceNotFound.Code, err.Code)
	suite.mockStore.AssertExpectations(suite.T())
}

func (suite *PreferenceServiceTestSuite) TestDeletePreference_InvalidKey() {
	ctx := context.Background()

	// Execute
	err := suite.service.DeletePreference(ctx, testUserID, "")

	// Assert
	suite.Error(err)
	suite.Equal(ErrorInvalidPreferenceKey.Code, err.Code)
}

func (suite *PreferenceServiceTestSuite) TestDeletePreference_StoreError() {
	ctx := context.Background()

	suite.mockStore.On("DeletePreference", ctx, testUserID, testKey).Return(errors.New("database error"))

	// Execute
	err := suite.service.DeletePreference(ctx, testUserID, testKey)

	// Assert
	suite.Error(err)
	suite.Equal(ErrorInternalServerError.Code, err.Code)
	suite.mockStore.AssertExpectations(suite.T())
}

func TestPreferenceServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PreferenceServiceTestSuite))
}

// Mock implementations for testing

type preferenceStoreInterfaceMock struct {
	mock.Mock
}

func newPreferenceStoreInterfaceMock(t *testing.T) *preferenceStoreInterfaceMock {
	mockStore := &preferenceStoreInterfaceMock{}
	mockStore.Test(t)
	return mockStore
}

func (m *preferenceStoreInterfaceMock) GetPreferenceByKey(ctx context.Context, userID, key string) (*Preference, error) {
	args := m.Called(ctx, userID, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Preference), args.Error(1)
}

func (m *preferenceStoreInterfaceMock) GetPreferencesByUserID(ctx context.Context, userID string) ([]Preference, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Preference), args.Error(1)
}

func (m *preferenceStoreInterfaceMock) UpsertPreference(ctx context.Context, userID, key, value string) error {
	args := m.Called(ctx, userID, key, value)
	return args.Error(0)
}

func (m *preferenceStoreInterfaceMock) DeletePreference(ctx context.Context, userID, key string) error {
	args := m.Called(ctx, userID, key)
	return args.Error(0)
}

type transactionerMock struct {
	mock.Mock
}

func newTransactionerMock(t *testing.T) *transactionerMock {
	mockTx := &transactionerMock{}
	mockTx.Test(t)
	return mockTx
}

func (m *transactionerMock) Transact(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *transactionerMock) GetTransactionerInterface() transaction.TransactionerInterface {
	args := m.Called()
	return args.Get(0).(transaction.TransactionerInterface)
}
