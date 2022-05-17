# tmdei-project
Master's Thesis Project

## Context

1. Create Sonarqube instance
2. Download all plugins that were requested
3. Move them to the plugins folder of the sonarqube instance
4. Delete them from the local machine
5. Login as admin
6. Change the admin password
7. Persist that information so that it can be used when doing work on the instance
8. Create a new user with enough permissions
9. Generate user token
10. Return the instance URL and the token to access it

### Create Sonarqube instance

In order to create the instance we can make use of [docker](https://hub.docker.com/_/sonarqube/)

We create a template Docker compose that can be then replaced by the details wanted by the user such as the specific version and the database to use

### Download all plugins that were requested

In order to download the plugins we can periodically download the list of versions of each plugin from [here](https://update.sonarsource.org/) where we populate a list of plugins with their versions and their download links which we use to download them.

### Login as admin

In order to login as admin we need to make use of the Web API

The endpoint to login is the following: `POST {url}/api/authentication/login` with the body: `{ login: "username", password: "password" }`

### Change the admin password

The endpoint to change the password is: `POST {url}/api/users/change_password` with the body: `{ login: "username", previousPassword: "prevPassword", password: "password" }`

### Create a new user with enough permissions

The endpoint to create a user is: `POST {url}/api/users/create` with the body: `{ email: "email", local: true, login: "username", name: "name", password: "password" }`

Permissions should be set for the user, the endpoint to do that is: `POST {url}/api/permissions/add_user` with the body: `{ login: "username", permission: "permission" }`

The permissions that should be set are: `scan` and `provisioning`.

### Generate user token

The endpoint to generate a token is: `POST {url}/api/user_tokens/generate` with the body: `{ login: "username", name: "token name" }`

### Return

Once everything is done and the Sonarqube instance is correctly configured we return the URL of the instance and the relevant information such as the Authentication Token that was generated for the created user.