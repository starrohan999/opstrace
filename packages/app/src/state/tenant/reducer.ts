/**
 * Copyright 2020 Opstrace, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { pluck, zipObj, pick, mergeDeepRight } from "ramda";

import { findTenantFromKey } from "state/tenant/utils";

import { createReducer, ActionType } from "typesafe-actions";
import { TenantRecords } from "./types";

import * as actions from "./actions";

type TenantActions = ActionType<typeof actions>;

type TenantState = {
  loading: boolean;
  tenants: TenantRecords;
};

const TenantInitialState: TenantState = {
  loading: true,
  tenants: {}
};

export const reducer = createReducer<TenantState, TenantActions>(
  TenantInitialState
)
  .handleAction(
    actions.setTenantList,
    (state, action): TenantState => {
      const tenantIds: string[] = pluck("id", action.payload);
      const tenants: TenantRecords = zipObj(tenantIds, action.payload);
      return mergeDeepRight(state, { loading: false, tenants: tenants });
    }
  )
  .handleAction(
    actions.alertmanagerLoaded,
    (state, action): TenantState => {
      const tenant = findTenantFromKey(state.tenants, action.payload.tenantKey);
      if (tenant) {
        return mergeDeepRight(state, {
          tenants: {
            [tenant.id]: {
              alertmanager: pick(["templates", "config", "online"])(
                action.payload
              )
            }
          }
        });
      } else {
        return state;
      }
    }
  )
  .handleAction(
    actions.updateAlertmanager,
    (state, action): TenantState => {
      const tenant = findTenantFromKey(state.tenants, action.payload.tenantKey);
      if (tenant) {
        return mergeDeepRight(state, {
          tenants: {
            [tenant.id]: {
              alertmanager: pick(["templates", "config"])(action.payload)
            }
          }
        });
      } else {
        return state;
      }
    }
  );
