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

import React from "react";
import { useParams } from "react-router-dom";
import Skeleton from "@material-ui/lab/Skeleton";
import Avatar from "@material-ui/core/Avatar";
import { useDispatch } from "react-redux";

import { Box } from "client/components/Box";
import Attribute from "client/components/Attribute";

import { deleteUser } from "state/user/actions";

import useUserList from "state/user/hooks/useUserList";
import Layout from "client/layout/MainContent";
import SideBar from "./Sidebar";

import { Card, CardContent, CardHeader } from "client/components/Card";
import { Button } from "client/components/Button";
import { usePickerService } from "client/services/Picker";
import { useCommandService } from "client/services/Command";
import useUser from "state/user/hooks/useUser";
import useCurrentUser from "state/user/hooks/useCurrentUser";

const UserDetail = () => {
  const params = useParams<{ userId: string }>();
  const users = useUserList();
  const currentUser = useCurrentUser();
  const user = useUser(params.userId);
  const dispatch = useDispatch();

  const { activatePickerWithText } = usePickerService(
    {
      title: `Delete ${user?.email}?`,
      activationPrefix: "delete user directly?:",
      disableFilter: true,
      disableInput: true,
      options: [
        {
          id: "yes",
          text: `yes`
        },
        {
          id: "no",
          text: "no"
        }
      ],
      onSelected: option => {
        if (option.id === "yes" && user?.id) dispatch(deleteUser(user?.id));
      }
    },
    [user?.id]
  );

  const cmdService = useCommandService();

  if (!user)
    return (
      <Layout sidebar={SideBar}>
        <Skeleton variant="rect" width="100%" height="100%" animation="wave" />
      </Layout>
    );
  else
    return (
      <Layout sidebar={SideBar}>
        <Box
          width="100%"
          height="100%"
          display="flex"
          justifyContent="center"
          alignItems="center"
          flexWrap="wrap"
          p={1}
        >
          <Box maxWidth={700}>
            <Card p={3}>
              <CardHeader
                titleTypographyProps={{ variant: "h5" }}
                avatar={
                  user.avatar ? (
                    <Avatar
                      alt={user.username}
                      style={{ width: 100, height: 100 }}
                      src={user.avatar}
                    />
                  ) : (
                    <Avatar
                      alt={user.username}
                      style={{ width: 100, height: 100 }}
                    >
                      {user.username.slice(0, 1).toUpperCase()}
                    </Avatar>
                  )
                }
                action={
                  <Box ml={3} display="flex" flexWrap="wrap">
                    {user.email === currentUser?.email ? (
                      <Box p={1}>
                        <Button
                          variant="outlined"
                          size="medium"
                          onClick={() => cmdService.executeCommand("logout")}
                        >
                          Logout
                        </Button>
                      </Box>
                    ) : null}
                    <Box p={1}>
                      <Button
                        variant="outlined"
                        size="medium"
                        disabled={users.length < 2}
                        onClick={() =>
                          activatePickerWithText("delete user directly?: ")
                        }
                      >
                        Delete
                      </Button>
                    </Box>
                  </Box>
                }
                title={user.username}
              />
              <CardContent>
                <Box display="flex">
                  <Box display="flex" flexDirection="column">
                    <Attribute.Key>Role:</Attribute.Key>
                    <Attribute.Key>Username:</Attribute.Key>
                    <Attribute.Key>Email:</Attribute.Key>
                    <Attribute.Key>Last Login:</Attribute.Key>
                    <Attribute.Key>Created:</Attribute.Key>
                  </Box>
                  <Box display="flex" flexDirection="column" flexGrow={1}>
                    <Attribute.Value>{user.role}</Attribute.Value>
                    <Attribute.Value>{user.username}</Attribute.Value>
                    <Attribute.Value>{user.email}</Attribute.Value>
                    <Attribute.Value>
                      {user.session_last_updated
                        ? user.session_last_updated
                        : "-"}
                    </Attribute.Value>
                    <Attribute.Value>{user.created_at}</Attribute.Value>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Box>
        </Box>
      </Layout>
    );
};

export default UserDetail;
