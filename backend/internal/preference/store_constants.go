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
	"github.com/asgardeo/thunder/internal/system/database/model"
)

// Database query constants.
const (
	queryIDGetPreferenceByKey     = "get_preference_by_key"
	queryIDGetPreferencesByUserID = "get_preferences_by_user_id"
	queryIDUpsertPreference       = "upsert_preference"
	queryIDDeletePreference       = "delete_preference"
)

// Query definitions.
var (
	queryGetPreferenceByKey = model.DBQuery{
		ID: queryIDGetPreferenceByKey,
		Queries: map[string]string{
			model.DefaultDBType: `SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT 
				FROM USER_PREFERENCE 
				WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3`,
		},
	}

	queryGetPreferencesByUserID = model.DBQuery{
		ID: queryIDGetPreferencesByUserID,
		Queries: map[string]string{
			model.DefaultDBType: `SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT 
				FROM USER_PREFERENCE 
				WHERE USER_ID = $1 AND DEPLOYMENT_ID = $2
				ORDER BY PREFERENCE_KEY ASC`,
		},
	}

	queryUpsertPreference = model.DBQuery{
		ID: queryIDUpsertPreference,
		Queries: map[string]string{
			model.DefaultDBType: `INSERT INTO USER_PREFERENCE (USER_ID, PREFERENCE_KEY, PREFERENCE_VALUE, DEPLOYMENT_ID) 
				VALUES ($1, $2, $3, $4) 
				ON CONFLICT (USER_ID, DEPLOYMENT_ID, PREFERENCE_KEY) 
				DO UPDATE SET PREFERENCE_VALUE = EXCLUDED.PREFERENCE_VALUE, UPDATED_AT = CURRENT_TIMESTAMP`,
		},
	}

	queryDeletePreference = model.DBQuery{
		ID: queryIDDeletePreference,
		Queries: map[string]string{
			model.DefaultDBType: `DELETE FROM USER_PREFERENCE 
				WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3`,
		},
	}
)
