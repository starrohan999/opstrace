/**
 * Copyright 2021 Opstrace, Inc.
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

import React from "react";
import { map } from "ramda";

import { useTenantByUrlSlug } from "state/tenant/hooks/useTenant";
import useAlertmanager from "state/tenant/hooks/useAlertmanager";

import { Tenant, Tenants } from "state/tenant/types";
import { PanelItem } from "client/components/Panel";

import Skeleton from "@material-ui/lab/Skeleton";

export const tenantToItem = (tenant: Tenant): PanelItem => {
  return { id: tenant.name, text: tenant.name, data: tenant };
};

export const tenantsToItems: (tenants: Tenants) => PanelItem[] = map(
  tenantToItem
);

// type WithTenantProps = {
//   Component: JSX.Element | React.ReactType;
//   urlSlug: string;
// };

// export const WithTenant = ({
//   Component,
//   urlSlug,
//   ...props
// }: WithTenantProps) => {
//   const tenant = useTenantByUrlSlug(urlSlug);

//   if (tenant) return <Component {...props} tenant={tenant} />;
//   else
//     return (
//       <Skeleton variant="rect" width="100%" height="100%" animation="wave" />
//     );
// };

export const withTenant = (
  Component: React.ReactType,
  tenantUrlSlug: string
) => {
  return (props: {}) => {
    const tenant = useTenantByUrlSlug(tenantUrlSlug);

    return tenant ? (
      <Component {...props} tenant={tenant} />
    ) : (
      <Skeleton variant="rect" width="100%" height="100%" animation="wave" />
    );
  };
};

// type WithAlertmanagerProps = {
//   Component: JSX.Element | JSX.Element[] | React.ReactType;
//   tenant?: Tenant;
// };

// export const WithAlertmanager = ({
//   Component,
//   tenant,
//   ...props
// }: WithAlertmanagerProps) => {
//   const alertmanager = useAlertmanager(tenant.name);

//   if (alertmanager)
//     return <Component {...props} tenant={tenant} alertmanager={alertmanager} />;
//   else
//     return (
//       <Skeleton variant="rect" width="100%" height="100%" animation="wave" />
//     );
// };

export const withAlertmanager = (
  Component: React.ReactType,
  tenantUrlSlug: string
) => {
  const ComponentWithAlertmanager = ({
    tenant,
    ...rest
  }: {
    tenant: Tenant;
  }) => {
    const alertmanager = useAlertmanager(tenant.name);

    return alertmanager ? (
      <Component {...rest} tenant={tenant} alertmanager={alertmanager} />
    ) : (
      <Skeleton variant="rect" width="100%" height="100%" animation="wave" />
    );
  };

  return withTenant(ComponentWithAlertmanager, tenantUrlSlug);
};
