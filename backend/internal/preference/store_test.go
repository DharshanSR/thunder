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
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/asgardeo/thunder/internal/system/config"
	"github.com/asgardeo/thunder/internal/system/database/model"
	"github.com/asgardeo/thunder/internal/system/database/provider"
	"github.com/asgardeo/thunder/tests/mocks/database/providermock"
)

const (
	testUserID       = "test-user-id"
	testDeploymentID = "test-deployment"
	testKey          = "theme"
	testValue        = "dark"
)

type PreferenceStoreTestSuite struct {
	suite.Suite
	store      *preferenceStore
	dbProvider *providermock.DBProviderInterfaceMock
	dbClient   *providermock.DBClientInterfaceMock
	sqlDB      *sql.DB
	sqlMock    sqlmock.Sqlmock
}

func (suite *PreferenceStoreTestSuite) SetupTest() {
	// Create mock DB and sqlmock
	db, mock, err := sqlmock.New()
	suite.Require().NoError(err)
	suite.sqlDB = db
	suite.sqlMock = mock

	// Create mock DB provider and client
	suite.dbProvider = providermock.NewDBProviderInterfaceMock(suite.T())
	suite.dbClient = providermock.NewDBClientInterfaceMock(suite.T())

	// Setup provider to return the mock client
	provider.SetDBProvider(suite.dbProvider)

	// Setup config
	config.InitThunderRuntime("sqlite")
	config.GetThunderRuntime().SetDeploymentID(testDeploymentID)

	// Create store
	suite.store = &preferenceStore{
		dbClient: suite.dbClient,
	}
}

func (suite *PreferenceStoreTestSuite) TearDownTest() {
	suite.sqlDB.Close()
}

func (suite *PreferenceStoreTestSuite) TestGetPreferenceByKey_Success() {
	ctx := context.Background()
	expectedPreference := &Preference{
		Key:       testKey,
		Value:     testValue,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup mock expectations
	suite.dbClient.On("GetQuery", queryGetPreferenceByKey).Return("SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", nil)

	rows := sqlmock.NewRows([]string{"PREFERENCE_KEY", "PREFERENCE_VALUE", "CREATED_AT", "UPDATED_AT"}).
		AddRow(expectedPreference.Key, expectedPreference.Value, expectedPreference.CreatedAt, expectedPreference.UpdatedAt)
	suite.sqlMock.ExpectQuery("SELECT PREFERENCE_KEY").
		WithArgs(testUserID, testKey, testDeploymentID).
		WillReturnRows(rows)

	suite.dbClient.On("QueryRowContext", ctx, "SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", testUserID, testKey, testDeploymentID).
		Return(suite.sqlDB.QueryRow("SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", testUserID, testKey, testDeploymentID))

	// Execute
	preference, err := suite.store.GetPreferenceByKey(ctx, testUserID, testKey)

	// Assert
	suite.NoError(err)
	suite.NotNil(preference)
	suite.Equal(expectedPreference.Key, preference.Key)
	suite.Equal(expectedPreference.Value, preference.Value)
}

func (suite *PreferenceStoreTestSuite) TestGetPreferenceByKey_NotFound() {
	ctx := context.Background()

	// Setup mock expectations
	suite.dbClient.On("GetQuery", queryGetPreferenceByKey).Return("SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", nil)

	suite.sqlMock.ExpectQuery("SELECT PREFERENCE_KEY").
		WithArgs(testUserID, testKey, testDeploymentID).
		WillReturnError(sql.ErrNoRows)

	suite.dbClient.On("QueryRowContext", ctx, "SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", testUserID, testKey, testDeploymentID).
		Return(suite.sqlDB.QueryRow("SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", testUserID, testKey, testDeploymentID))

	// Execute
	preference, err := suite.store.GetPreferenceByKey(ctx, testUserID, testKey)

	// Assert
	suite.Error(err)
	suite.True(errors.Is(err, ErrPreferenceNotFound))
	suite.Nil(preference)
}

func (suite *PreferenceStoreTestSuite) TestGetPreferencesByUserID_Success() {
	ctx := context.Background()

	// Setup mock expectations
	suite.dbClient.On("GetQuery", queryGetPreferencesByUserID).Return("SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND DEPLOYMENT_ID = $2 ORDER BY PREFERENCE_KEY ASC", nil)

	rows := sqlmock.NewRows([]string{"PREFERENCE_KEY", "PREFERENCE_VALUE", "CREATED_AT", "UPDATED_AT"}).
		AddRow("key1", "value1", time.Now(), time.Now()).
		AddRow("key2", "value2", time.Now(), time.Now())
	suite.sqlMock.ExpectQuery("SELECT PREFERENCE_KEY").
		WithArgs(testUserID, testDeploymentID).
		WillReturnRows(rows)

	suite.dbClient.On("QueryContext", ctx, "SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND DEPLOYMENT_ID = $2 ORDER BY PREFERENCE_KEY ASC", testUserID, testDeploymentID).
		Return(suite.sqlDB.Query("SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND DEPLOYMENT_ID = $2 ORDER BY PREFERENCE_KEY ASC", testUserID, testDeploymentID))

	// Execute
	preferences, err := suite.store.GetPreferencesByUserID(ctx, testUserID)

	// Assert
	suite.NoError(err)
	suite.NotNil(preferences)
	suite.Len(preferences, 2)
}

func (suite *PreferenceStoreTestSuite) TestGetPreferencesByUserID_Empty() {
	ctx := context.Background()

	// Setup mock expectations
	suite.dbClient.On("GetQuery", queryGetPreferencesByUserID).Return("SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND DEPLOYMENT_ID = $2 ORDER BY PREFERENCE_KEY ASC", nil)

	rows := sqlmock.NewRows([]string{"PREFERENCE_KEY", "PREFERENCE_VALUE", "CREATED_AT", "UPDATED_AT"})
	suite.sqlMock.ExpectQuery("SELECT PREFERENCE_KEY").
		WithArgs(testUserID, testDeploymentID).
		WillReturnRows(rows)

	suite.dbClient.On("QueryContext", ctx, "SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND DEPLOYMENT_ID = $2 ORDER BY PREFERENCE_KEY ASC", testUserID, testDeploymentID).
		Return(suite.sqlDB.Query("SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT FROM USER_PREFERENCE WHERE USER_ID = $1 AND DEPLOYMENT_ID = $2 ORDER BY PREFERENCE_KEY ASC", testUserID, testDeploymentID))

	// Execute
	preferences, err := suite.store.GetPreferencesByUserID(ctx, testUserID)

	// Assert
	suite.NoError(err)
	suite.NotNil(preferences)
	suite.Len(preferences, 0)
}

func (suite *PreferenceStoreTestSuite) TestUpsertPreference_Success() {
	ctx := context.Background()

	// Setup mock expectations
	suite.dbClient.On("GetQuery", queryUpsertPreference).Return("INSERT INTO USER_PREFERENCE (USER_ID, PREFERENCE_KEY, PREFERENCE_VALUE, DEPLOYMENT_ID) VALUES ($1, $2, $3, $4) ON CONFLICT (USER_ID, DEPLOYMENT_ID, PREFERENCE_KEY) DO UPDATE SET PREFERENCE_VALUE = EXCLUDED.PREFERENCE_VALUE, UPDATED_AT = CURRENT_TIMESTAMP", nil)

	suite.sqlMock.ExpectExec("INSERT INTO USER_PREFERENCE").
		WithArgs(testUserID, testKey, testValue, testDeploymentID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	suite.dbClient.On("ExecContext", ctx, "INSERT INTO USER_PREFERENCE (USER_ID, PREFERENCE_KEY, PREFERENCE_VALUE, DEPLOYMENT_ID) VALUES ($1, $2, $3, $4) ON CONFLICT (USER_ID, DEPLOYMENT_ID, PREFERENCE_KEY) DO UPDATE SET PREFERENCE_VALUE = EXCLUDED.PREFERENCE_VALUE, UPDATED_AT = CURRENT_TIMESTAMP", testUserID, testKey, testValue, testDeploymentID).
		Return(suite.sqlDB.Exec("INSERT INTO USER_PREFERENCE (USER_ID, PREFERENCE_KEY, PREFERENCE_VALUE, DEPLOYMENT_ID) VALUES ($1, $2, $3, $4) ON CONFLICT (USER_ID, DEPLOYMENT_ID, PREFERENCE_KEY) DO UPDATE SET PREFERENCE_VALUE = EXCLUDED.PREFERENCE_VALUE, UPDATED_AT = CURRENT_TIMESTAMP", testUserID, testKey, testValue, testDeploymentID))

	// Execute
	err := suite.store.UpsertPreference(ctx, testUserID, testKey, testValue)

	// Assert
	suite.NoError(err)
}

func (suite *PreferenceStoreTestSuite) TestDeletePreference_Success() {
	ctx := context.Background()

	// Setup mock expectations
	suite.dbClient.On("GetQuery", queryDeletePreference).Return("DELETE FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", nil)

	suite.sqlMock.ExpectExec("DELETE FROM USER_PREFERENCE").
		WithArgs(testUserID, testKey, testDeploymentID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	suite.dbClient.On("ExecContext", ctx, "DELETE FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", testUserID, testKey, testDeploymentID).
		Return(suite.sqlDB.Exec("DELETE FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", testUserID, testKey, testDeploymentID))

	// Execute
	err := suite.store.DeletePreference(ctx, testUserID, testKey)

	// Assert
	suite.NoError(err)
}

func (suite *PreferenceStoreTestSuite) TestDeletePreference_NotFound() {
	ctx := context.Background()

	// Setup mock expectations
	suite.dbClient.On("GetQuery", queryDeletePreference).Return("DELETE FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", nil)

	suite.sqlMock.ExpectExec("DELETE FROM USER_PREFERENCE").
		WithArgs(testUserID, testKey, testDeploymentID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	suite.dbClient.On("ExecContext", ctx, "DELETE FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", testUserID, testKey, testDeploymentID).
		Return(suite.sqlDB.Exec("DELETE FROM USER_PREFERENCE WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3", testUserID, testKey, testDeploymentID))

	// Execute
	err := suite.store.DeletePreference(ctx, testUserID, testKey)

	// Assert
	suite.Error(err)
	suite.True(errors.Is(err, ErrPreferenceNotFound))
}

func (suite *PreferenceStoreTestSuite) TestGetQuery_Error() {
	ctx := context.Background()

	// Setup mock to return error when getting query
	suite.dbClient.On("GetQuery", queryGetPreferenceByKey).Return("", model.ErrQueryNotFound)

	// Execute
	preference, err := suite.store.GetPreferenceByKey(ctx, testUserID, testKey)

	// Assert
	suite.Error(err)
	suite.True(errors.Is(err, model.ErrQueryNotFound))
	suite.Nil(preference)
}

func TestPreferenceStoreTestSuite(t *testing.T) {
	suite.Run(t, new(PreferenceStoreTestSuite))
}
