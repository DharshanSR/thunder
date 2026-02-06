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

// Query definitions.
var (
	queryGetPreferenceByKey = model.DBQuery{
		ID: "PREF-01",
		Query: `SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT 
			FROM USER_PREFERENCE 
			WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3`,
	}

	queryGetPreferencesByUserID = model.DBQuery{
		ID: "PREF-02",
		Query: `SELECT PREFERENCE_KEY, PREFERENCE_VALUE, CREATED_AT, UPDATED_AT 
			FROM USER_PREFERENCE 
			WHERE USER_ID = $1 AND DEPLOYMENT_ID = $2
			ORDER BY PREFERENCE_KEY ASC`,
	}

	queryUpsertPreference = model.DBQuery{
		ID: "PREF-03",
		Query: `INSERT INTO USER_PREFERENCE (USER_ID, PREFERENCE_KEY, PREFERENCE_VALUE, DEPLOYMENT_ID) 
			VALUES ($1, $2, $3, $4) 
			ON CONFLICT (USER_ID, DEPLOYMENT_ID, PREFERENCE_KEY) 
			DO UPDATE SET PREFERENCE_VALUE = EXCLUDED.PREFERENCE_VALUE, UPDATED_AT = CURRENT_TIMESTAMP`,
	}

	queryDeletePreference = model.DBQuery{
		ID: "PREF-04",
		Query: `DELETE FROM USER_PREFERENCE 
			WHERE USER_ID = $1 AND PREFERENCE_KEY = $2 AND DEPLOYMENT_ID = $3`,
	}
)
