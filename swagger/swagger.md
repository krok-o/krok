


# Krok.
Documentation the Krok API.
  

## Informations

### Version

0.0.1

### License

[Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  operations

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /rest/api/1/command/add-command-rel-for-platform/{cmdid}/{repoid} | [add command rel for platform command](#add-command-rel-for-platform-command) | Adds a connection to a platform for a command. Defines what platform a command supports. These commands will only be able to run for those platforms. |
| POST | /rest/api/1/command/add-command-rel-for-repository/{cmdid}/{repoid} | [add command rel for repository command](#add-command-rel-for-repository-command) | Add a connection to a repository. This will make this command to be executed for events for that repository. |
| POST | /rest/api/1/user/apikey/generate/{name} | [create Api key](#create-api-key) | Creates an api key pair for a given user. |
| POST | /rest/api/1/repository | [create repository](#create-repository) |  |
| POST | /rest/api/1/vault/secret | [create secret](#create-secret) | Create a new secure secret. |
| POST | /rest/api/1/user | [create user](#create-user) |  |
| POST | /rest/api/1/vcs-token | [create vcs token](#create-vcs-token) | Create a new token for a platform like Github, Gitlab, Gitea... |
| DELETE | /rest/api/1/user/apikey/delete/{keyid} | [delete Api key](#delete-api-key) | Deletes a set of api keys for a given user with a given id. |
| DELETE | /rest/api/1/command/{id} | [delete command](#delete-command) | Deletes given command. |
| DELETE | /rest/api/1/command/settings/{id} | [delete command setting](#delete-command-setting) | Deletes a given command setting. |
| DELETE | /rest/api/1/repository/{id} | [delete repository](#delete-repository) | Deletes the given repository. |
| DELETE | /rest/api/1/vault/secret/{name} | [delete secret](#delete-secret) | Deletes the given secret. |
| DELETE | /rest/api/1/user/{id} | [delete user](#delete-user) | Deletes the given user. |
| GET | /rest/api/1/user/apikey/{keyid} | [get Api keys](#get-api-keys) | Returns a given api key. |
| GET | /rest/api/1/command/{id} | [get command](#get-command) | Returns a specific command. |
| GET | /rest/api/1/command/run/{id} | [get command run](#get-command-run) | Returns details about a command run. |
| GET | /rest/api/1/command/settings/{id} | [get command setting](#get-command-setting) | Get a specific setting. |
| GET | /rest/api/1/event/{id} | [get event](#get-event) | Get a specific event. |
| GET | /rest/api/1/repository/{id} | [get repository](#get-repository) | Gets the repository with the corresponding ID. |
| GET | /rest/api/1/vault/secret/{name} | [get secret](#get-secret) | Get a specific secret. |
| POST | /rest/api/1/get-token | [get token](#get-token) | Creates a JWT token for a given api key pair. |
| GET | /rest/api/1/user/{id} | [get user](#get-user) | Gets the user with the corresponding ID. |
| POST | /rest/api/1/hooks/{rid}/{vid}/callback | [hook handler](#hook-handler) | Handle the hooks created by the platform. |
| POST | /rest/api/1/user/apikey | [list Api keys](#list-api-keys) | Lists all api keys for a given user. |
| POST | /rest/api/1/command/{id}/settings | [list command settings](#list-command-settings) | List settings for a command. |
| POST | /rest/api/1/commands | [list commands](#list-commands) |  |
| POST | /rest/api/1/events/{repoid} | [list events](#list-events) | List events for a repository. |
| POST | /rest/api/1/repositories | [list repositories](#list-repositories) |  |
| POST | /rest/api/1/vault/secrets | [list secrets](#list-secrets) | List all settings without the values. |
| GET | /rest/api/1/supported-platforms | [list supported platforms](#list-supported-platforms) | Lists all supported platforms. |
| POST | /rest/api/1/users | [list users](#list-users) |  |
| POST | /rest/api/1/auth/refresh | [refresh token](#refresh-token) | Refresh the authentication token. |
| POST | /rest/api/1/command/remove-command-rel-for-platform/{cmdid}/{repoid} | [remove command rel for platform command](#remove-command-rel-for-platform-command) | Remove a relationship to a platform. This command will no longer be running for that platform events. |
| POST | /rest/api/1/command/remove-command-rel-for-repository/{cmdid}/{repoid} | [remove command rel for repository command](#remove-command-rel-for-repository-command) | Remove a relationship to a repository. This command will no longer be running for that repository events. |
| POST | /rest/api/1/command/update | [update command](#update-command) | Updates a given command. |
| POST | /rest/api/1/command/settings/update | [update command setting](#update-command-setting) | Create a new command setting. |
| POST | /rest/api/1/repository/update | [update repository](#update-repository) | Updates an existing repository. |
| POST | /rest/api/1/vault/secret/update | [update secret](#update-secret) | Updates an existing secret. |
| POST | /rest/api/1/user/update | [update user](#update-user) | Updates an existing user. |
| POST | /rest/api/1/command | [upload command](#upload-command) | Upload a command. To set up anything for the command, like schedules etc, |
| GET | /rest/api/1/auth/callback | [user callback](#user-callback) | This is the url to which Google calls back after a successful login. |
| GET | /rest/api/1/auth/login | [user login](#user-login) | User login. |
  


## Paths

### <span id="add-command-rel-for-platform-command"></span> Adds a connection to a platform for a command. Defines what platform a command supports. These commands will only be able to run for those platforms. (*addCommandRelForPlatformCommand*)

```
POST /rest/api/1/command/add-command-rel-for-platform/{cmdid}/{repoid}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| cmdid | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |
| repoid | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#add-command-rel-for-platform-command-200) | OK | successfully added relationship |  | [schema](#add-command-rel-for-platform-command-200-schema) |
| [400](#add-command-rel-for-platform-command-400) | Bad Request | invalid ids or platform not found |  | [schema](#add-command-rel-for-platform-command-400-schema) |
| [500](#add-command-rel-for-platform-command-500) | Internal Server Error | failed to add command relationship to platform |  | [schema](#add-command-rel-for-platform-command-500-schema) |

#### Responses


##### <span id="add-command-rel-for-platform-command-200"></span> 200 - successfully added relationship
Status: OK

###### <span id="add-command-rel-for-platform-command-200-schema"></span> Schema

##### <span id="add-command-rel-for-platform-command-400"></span> 400 - invalid ids or platform not found
Status: Bad Request

###### <span id="add-command-rel-for-platform-command-400-schema"></span> Schema
   
  

any

##### <span id="add-command-rel-for-platform-command-500"></span> 500 - failed to add command relationship to platform
Status: Internal Server Error

###### <span id="add-command-rel-for-platform-command-500-schema"></span> Schema
   
  

any

### <span id="add-command-rel-for-repository-command"></span> Add a connection to a repository. This will make this command to be executed for events for that repository. (*addCommandRelForRepositoryCommand*)

```
POST /rest/api/1/command/add-command-rel-for-repository/{cmdid}/{repoid}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| cmdid | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |
| repoid | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#add-command-rel-for-repository-command-200) | OK | successfully added relationship |  | [schema](#add-command-rel-for-repository-command-200-schema) |
| [400](#add-command-rel-for-repository-command-400) | Bad Request | invalid ids or repositroy not found |  | [schema](#add-command-rel-for-repository-command-400-schema) |
| [500](#add-command-rel-for-repository-command-500) | Internal Server Error | failed to add relationship |  | [schema](#add-command-rel-for-repository-command-500-schema) |

#### Responses


##### <span id="add-command-rel-for-repository-command-200"></span> 200 - successfully added relationship
Status: OK

###### <span id="add-command-rel-for-repository-command-200-schema"></span> Schema

##### <span id="add-command-rel-for-repository-command-400"></span> 400 - invalid ids or repositroy not found
Status: Bad Request

###### <span id="add-command-rel-for-repository-command-400-schema"></span> Schema
   
  

any

##### <span id="add-command-rel-for-repository-command-500"></span> 500 - failed to add relationship
Status: Internal Server Error

###### <span id="add-command-rel-for-repository-command-500-schema"></span> Schema
   
  

any

### <span id="create-api-key"></span> Creates an api key pair for a given user. (*createApiKey*)

```
POST /rest/api/1/user/apikey/generate/{name}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| name | `path` | string | `string` |  | ✓ |  | the name of the key |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#create-api-key-200) | OK | the generated api key pair |  | [schema](#create-api-key-200-schema) |
| [400](#create-api-key-400) | Bad Request | failed to generate unique key or value |  | [schema](#create-api-key-400-schema) |
| [500](#create-api-key-500) | Internal Server Error | when failed to get user context |  | [schema](#create-api-key-500-schema) |

#### Responses


##### <span id="create-api-key-200"></span> 200 - the generated api key pair
Status: OK

###### <span id="create-api-key-200-schema"></span> Schema
   
  

[APIKey](#api-key)

##### <span id="create-api-key-400"></span> 400 - failed to generate unique key or value
Status: Bad Request

###### <span id="create-api-key-400-schema"></span> Schema
   
  

any

##### <span id="create-api-key-500"></span> 500 - when failed to get user context
Status: Internal Server Error

###### <span id="create-api-key-500-schema"></span> Schema
   
  

any

### <span id="create-repository"></span> create repository (*createRepository*)

```
POST /rest/api/1/repository
```

Creates a new repository

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| repository | `body` | [Repository](#repository) | `models.Repository` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#create-repository-200) | OK | the created repository |  | [schema](#create-repository-200-schema) |
| [400](#create-repository-400) | Bad Request | failed to generate unique key or value |  | [schema](#create-repository-400-schema) |
| [500](#create-repository-500) | Internal Server Error | when failed to get user context |  | [schema](#create-repository-500-schema) |

#### Responses


##### <span id="create-repository-200"></span> 200 - the created repository
Status: OK

###### <span id="create-repository-200-schema"></span> Schema
   
  

[Repository](#repository)

##### <span id="create-repository-400"></span> 400 - failed to generate unique key or value
Status: Bad Request

###### <span id="create-repository-400-schema"></span> Schema
   
  

any

##### <span id="create-repository-500"></span> 500 - when failed to get user context
Status: Internal Server Error

###### <span id="create-repository-500-schema"></span> Schema
   
  

any

### <span id="create-secret"></span> Create a new secure secret. (*createSecret*)

```
POST /rest/api/1/vault/secret
```

#### Consumes
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| secret | `body` | [VaultSetting](#vault-setting) | `models.VaultSetting` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#create-secret-200) | OK | OK setting successfully create |  | [schema](#create-secret-200-schema) |
| [400](#create-secret-400) | Bad Request | invalid json payload |  | [schema](#create-secret-400-schema) |
| [500](#create-secret-500) | Internal Server Error | failed to create secret |  | [schema](#create-secret-500-schema) |

#### Responses


##### <span id="create-secret-200"></span> 200 - OK setting successfully create
Status: OK

###### <span id="create-secret-200-schema"></span> Schema

##### <span id="create-secret-400"></span> 400 - invalid json payload
Status: Bad Request

###### <span id="create-secret-400-schema"></span> Schema
   
  

any

##### <span id="create-secret-500"></span> 500 - failed to create secret
Status: Internal Server Error

###### <span id="create-secret-500-schema"></span> Schema
   
  

any

### <span id="create-user"></span> create user (*createUser*)

```
POST /rest/api/1/user
```

Creates a new user

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| user | `body` | [User](#user) | `models.User` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#create-user-200) | OK | the created user |  | [schema](#create-user-200-schema) |
| [400](#create-user-400) | Bad Request | invalid json payload |  | [schema](#create-user-400-schema) |
| [500](#create-user-500) | Internal Server Error | failed to create user or generating a new api key |  | [schema](#create-user-500-schema) |

#### Responses


##### <span id="create-user-200"></span> 200 - the created user
Status: OK

###### <span id="create-user-200-schema"></span> Schema
   
  

[User](#user)

##### <span id="create-user-400"></span> 400 - invalid json payload
Status: Bad Request

###### <span id="create-user-400-schema"></span> Schema
   
  

any

##### <span id="create-user-500"></span> 500 - failed to create user or generating a new api key
Status: Internal Server Error

###### <span id="create-user-500-schema"></span> Schema
   
  

any

### <span id="create-vcs-token"></span> Create a new token for a platform like Github, Gitlab, Gitea... (*createVcsToken*)

```
POST /rest/api/1/vcs-token
```

#### Consumes
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| secret | `body` | [VCSToken](#v-c-s-token) | `models.VCSToken` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#create-vcs-token-200) | OK | OK setting successfully create |  | [schema](#create-vcs-token-200-schema) |
| [400](#create-vcs-token-400) | Bad Request | invalid json payload |  | [schema](#create-vcs-token-400-schema) |
| [500](#create-vcs-token-500) | Internal Server Error | failed to create secret |  | [schema](#create-vcs-token-500-schema) |

#### Responses


##### <span id="create-vcs-token-200"></span> 200 - OK setting successfully create
Status: OK

###### <span id="create-vcs-token-200-schema"></span> Schema

##### <span id="create-vcs-token-400"></span> 400 - invalid json payload
Status: Bad Request

###### <span id="create-vcs-token-400-schema"></span> Schema
   
  

any

##### <span id="create-vcs-token-500"></span> 500 - failed to create secret
Status: Internal Server Error

###### <span id="create-vcs-token-500-schema"></span> Schema
   
  

any

### <span id="delete-api-key"></span> Deletes a set of api keys for a given user with a given id. (*deleteApiKey*)

```
DELETE /rest/api/1/user/apikey/delete/{keyid}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| keyid | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the key to delete |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-api-key-200) | OK | OK in case the deletion was successful |  | [schema](#delete-api-key-200-schema) |
| [400](#delete-api-key-400) | Bad Request | in case of missing user context or invalid ID |  | [schema](#delete-api-key-400-schema) |
| [500](#delete-api-key-500) | Internal Server Error | when the deletion operation failed |  | [schema](#delete-api-key-500-schema) |

#### Responses


##### <span id="delete-api-key-200"></span> 200 - OK in case the deletion was successful
Status: OK

###### <span id="delete-api-key-200-schema"></span> Schema

##### <span id="delete-api-key-400"></span> 400 - in case of missing user context or invalid ID
Status: Bad Request

###### <span id="delete-api-key-400-schema"></span> Schema
   
  

any

##### <span id="delete-api-key-500"></span> 500 - when the deletion operation failed
Status: Internal Server Error

###### <span id="delete-api-key-500-schema"></span> Schema
   
  

any

### <span id="delete-command"></span> Deletes given command. (*deleteCommand*)

```
DELETE /rest/api/1/command/{id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the command to delete |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-command-200) | OK | OK in case the deletion was successful |  | [schema](#delete-command-200-schema) |
| [400](#delete-command-400) | Bad Request | in case of missing user context or invalid ID |  | [schema](#delete-command-400-schema) |
| [500](#delete-command-500) | Internal Server Error | when the deletion operation failed |  | [schema](#delete-command-500-schema) |

#### Responses


##### <span id="delete-command-200"></span> 200 - OK in case the deletion was successful
Status: OK

###### <span id="delete-command-200-schema"></span> Schema

##### <span id="delete-command-400"></span> 400 - in case of missing user context or invalid ID
Status: Bad Request

###### <span id="delete-command-400-schema"></span> Schema
   
  

any

##### <span id="delete-command-500"></span> 500 - when the deletion operation failed
Status: Internal Server Error

###### <span id="delete-command-500-schema"></span> Schema
   
  

any

### <span id="delete-command-setting"></span> Deletes a given command setting. (*deleteCommandSetting*)

```
DELETE /rest/api/1/command/settings/{id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the command setting to delete |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-command-setting-200) | OK | OK in case the deletion was successful |  | [schema](#delete-command-setting-200-schema) |
| [400](#delete-command-setting-400) | Bad Request | invalid id |  | [schema](#delete-command-setting-400-schema) |
| [404](#delete-command-setting-404) | Not Found | command setting not found |  | [schema](#delete-command-setting-404-schema) |
| [500](#delete-command-setting-500) | Internal Server Error | when the deletion operation failed |  | [schema](#delete-command-setting-500-schema) |

#### Responses


##### <span id="delete-command-setting-200"></span> 200 - OK in case the deletion was successful
Status: OK

###### <span id="delete-command-setting-200-schema"></span> Schema

##### <span id="delete-command-setting-400"></span> 400 - invalid id
Status: Bad Request

###### <span id="delete-command-setting-400-schema"></span> Schema
   
  

any

##### <span id="delete-command-setting-404"></span> 404 - command setting not found
Status: Not Found

###### <span id="delete-command-setting-404-schema"></span> Schema
   
  

any

##### <span id="delete-command-setting-500"></span> 500 - when the deletion operation failed
Status: Internal Server Error

###### <span id="delete-command-setting-500-schema"></span> Schema
   
  

any

### <span id="delete-repository"></span> Deletes the given repository. (*deleteRepository*)

```
DELETE /rest/api/1/repository/{id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the repository to delete |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-repository-200) | OK | OK in case the deletion was successful |  | [schema](#delete-repository-200-schema) |
| [400](#delete-repository-400) | Bad Request | in case of missing user context or invalid ID |  | [schema](#delete-repository-400-schema) |
| [404](#delete-repository-404) | Not Found | in case of repository not found |  | [schema](#delete-repository-404-schema) |
| [500](#delete-repository-500) | Internal Server Error | when the deletion operation failed |  | [schema](#delete-repository-500-schema) |

#### Responses


##### <span id="delete-repository-200"></span> 200 - OK in case the deletion was successful
Status: OK

###### <span id="delete-repository-200-schema"></span> Schema

##### <span id="delete-repository-400"></span> 400 - in case of missing user context or invalid ID
Status: Bad Request

###### <span id="delete-repository-400-schema"></span> Schema
   
  

any

##### <span id="delete-repository-404"></span> 404 - in case of repository not found
Status: Not Found

###### <span id="delete-repository-404-schema"></span> Schema
   
  

any

##### <span id="delete-repository-500"></span> 500 - when the deletion operation failed
Status: Internal Server Error

###### <span id="delete-repository-500-schema"></span> Schema
   
  

any

### <span id="delete-secret"></span> Deletes the given secret. (*deleteSecret*)

```
DELETE /rest/api/1/vault/secret/{name}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| name | `path` | string | `string` |  | ✓ |  | The key of the secret |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-secret-200) | OK | OK in case the deletion was successful |  | [schema](#delete-secret-200-schema) |
| [400](#delete-secret-400) | Bad Request | in case of missing name |  | [schema](#delete-secret-400-schema) |
| [404](#delete-secret-404) | Not Found | in case the secret was not found |  | [schema](#delete-secret-404-schema) |
| [500](#delete-secret-500) | Internal Server Error | when the deletion operation failed |  | [schema](#delete-secret-500-schema) |

#### Responses


##### <span id="delete-secret-200"></span> 200 - OK in case the deletion was successful
Status: OK

###### <span id="delete-secret-200-schema"></span> Schema

##### <span id="delete-secret-400"></span> 400 - in case of missing name
Status: Bad Request

###### <span id="delete-secret-400-schema"></span> Schema
   
  

any

##### <span id="delete-secret-404"></span> 404 - in case the secret was not found
Status: Not Found

###### <span id="delete-secret-404-schema"></span> Schema
   
  

any

##### <span id="delete-secret-500"></span> 500 - when the deletion operation failed
Status: Internal Server Error

###### <span id="delete-secret-500-schema"></span> Schema
   
  

any

### <span id="delete-user"></span> Deletes the given user. (*deleteUser*)

```
DELETE /rest/api/1/user/{id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the user to delete |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-user-200) | OK | OK in case the deletion was successful |  | [schema](#delete-user-200-schema) |
| [400](#delete-user-400) | Bad Request | in case of missing user context or invalid ID |  | [schema](#delete-user-400-schema) |
| [404](#delete-user-404) | Not Found | in case of user not found |  | [schema](#delete-user-404-schema) |
| [500](#delete-user-500) | Internal Server Error | when the deletion operation failed |  | [schema](#delete-user-500-schema) |

#### Responses


##### <span id="delete-user-200"></span> 200 - OK in case the deletion was successful
Status: OK

###### <span id="delete-user-200-schema"></span> Schema

##### <span id="delete-user-400"></span> 400 - in case of missing user context or invalid ID
Status: Bad Request

###### <span id="delete-user-400-schema"></span> Schema
   
  

any

##### <span id="delete-user-404"></span> 404 - in case of user not found
Status: Not Found

###### <span id="delete-user-404-schema"></span> Schema
   
  

any

##### <span id="delete-user-500"></span> 500 - when the deletion operation failed
Status: Internal Server Error

###### <span id="delete-user-500-schema"></span> Schema
   
  

any

### <span id="get-api-keys"></span> Returns a given api key. (*getApiKeys*)

```
GET /rest/api/1/user/apikey/{keyid}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| keyid | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the key to return |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-api-keys-200) | OK |  |  | [schema](#get-api-keys-200-schema) |
| [500](#get-api-keys-500) | Internal Server Error | failed to get user context |  | [schema](#get-api-keys-500-schema) |

#### Responses


##### <span id="get-api-keys-200"></span> 200
Status: OK

###### <span id="get-api-keys-200-schema"></span> Schema
   
  

[APIKey](#api-key)

##### <span id="get-api-keys-500"></span> 500 - failed to get user context
Status: Internal Server Error

###### <span id="get-api-keys-500-schema"></span> Schema
   
  

any

### <span id="get-command"></span> Returns a specific command. (*getCommand*)

```
GET /rest/api/1/command/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-command-200) | OK |  |  | [schema](#get-command-200-schema) |
| [400](#get-command-400) | Bad Request | invalid command id |  | [schema](#get-command-400-schema) |
| [500](#get-command-500) | Internal Server Error | failed to get user context |  | [schema](#get-command-500-schema) |

#### Responses


##### <span id="get-command-200"></span> 200
Status: OK

###### <span id="get-command-200-schema"></span> Schema
   
  

[Command](#command)

##### <span id="get-command-400"></span> 400 - invalid command id
Status: Bad Request

###### <span id="get-command-400-schema"></span> Schema
   
  

any

##### <span id="get-command-500"></span> 500 - failed to get user context
Status: Internal Server Error

###### <span id="get-command-500-schema"></span> Schema
   
  

any

### <span id="get-command-run"></span> Returns details about a command run. (*getCommandRun*)

```
GET /rest/api/1/command/run/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-command-run-200) | OK |  |  | [schema](#get-command-run-200-schema) |
| [400](#get-command-run-400) | Bad Request | invalid command id |  | [schema](#get-command-run-400-schema) |
| [404](#get-command-run-404) | Not Found | command run not found |  | [schema](#get-command-run-404-schema) |
| [500](#get-command-run-500) | Internal Server Error | failed to get command run |  | [schema](#get-command-run-500-schema) |

#### Responses


##### <span id="get-command-run-200"></span> 200
Status: OK

###### <span id="get-command-run-200-schema"></span> Schema
   
  

[CommandRun](#command-run)

##### <span id="get-command-run-400"></span> 400 - invalid command id
Status: Bad Request

###### <span id="get-command-run-400-schema"></span> Schema
   
  

any

##### <span id="get-command-run-404"></span> 404 - command run not found
Status: Not Found

###### <span id="get-command-run-404-schema"></span> Schema

##### <span id="get-command-run-500"></span> 500 - failed to get command run
Status: Internal Server Error

###### <span id="get-command-run-500-schema"></span> Schema
   
  

any

### <span id="get-command-setting"></span> Get a specific setting. (*getCommandSetting*)

```
GET /rest/api/1/command/settings/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the command setting to retrieve |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-command-setting-200) | OK |  |  | [schema](#get-command-setting-200-schema) |
| [400](#get-command-setting-400) | Bad Request | invalid command id |  | [schema](#get-command-setting-400-schema) |
| [404](#get-command-setting-404) | Not Found | command setting not found |  | [schema](#get-command-setting-404-schema) |
| [500](#get-command-setting-500) | Internal Server Error | failed to get command setting |  | [schema](#get-command-setting-500-schema) |

#### Responses


##### <span id="get-command-setting-200"></span> 200
Status: OK

###### <span id="get-command-setting-200-schema"></span> Schema
   
  

[CommandSetting](#command-setting)

##### <span id="get-command-setting-400"></span> 400 - invalid command id
Status: Bad Request

###### <span id="get-command-setting-400-schema"></span> Schema
   
  

any

##### <span id="get-command-setting-404"></span> 404 - command setting not found
Status: Not Found

###### <span id="get-command-setting-404-schema"></span> Schema
   
  

any

##### <span id="get-command-setting-500"></span> 500 - failed to get command setting
Status: Internal Server Error

###### <span id="get-command-setting-500-schema"></span> Schema
   
  

any

### <span id="get-event"></span> Get a specific event. (*getEvent*)

```
GET /rest/api/1/event/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the event to retrieve |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-event-200) | OK |  |  | [schema](#get-event-200-schema) |
| [400](#get-event-400) | Bad Request | invalid event id |  | [schema](#get-event-400-schema) |
| [404](#get-event-404) | Not Found | event not found |  | [schema](#get-event-404-schema) |
| [500](#get-event-500) | Internal Server Error | failed to get event |  | [schema](#get-event-500-schema) |

#### Responses


##### <span id="get-event-200"></span> 200
Status: OK

###### <span id="get-event-200-schema"></span> Schema
   
  

[Event](#event)

##### <span id="get-event-400"></span> 400 - invalid event id
Status: Bad Request

###### <span id="get-event-400-schema"></span> Schema
   
  

any

##### <span id="get-event-404"></span> 404 - event not found
Status: Not Found

###### <span id="get-event-404-schema"></span> Schema
   
  

any

##### <span id="get-event-500"></span> 500 - failed to get event
Status: Internal Server Error

###### <span id="get-event-500-schema"></span> Schema
   
  

any

### <span id="get-repository"></span> Gets the repository with the corresponding ID. (*getRepository*)

```
GET /rest/api/1/repository/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-repository-200) | OK |  |  | [schema](#get-repository-200-schema) |
| [400](#get-repository-400) | Bad Request | invalid repository id |  | [schema](#get-repository-400-schema) |
| [404](#get-repository-404) | Not Found | repository not found |  | [schema](#get-repository-404-schema) |
| [500](#get-repository-500) | Internal Server Error | failed to get repository |  | [schema](#get-repository-500-schema) |

#### Responses


##### <span id="get-repository-200"></span> 200
Status: OK

###### <span id="get-repository-200-schema"></span> Schema
   
  

[Repository](#repository)

##### <span id="get-repository-400"></span> 400 - invalid repository id
Status: Bad Request

###### <span id="get-repository-400-schema"></span> Schema
   
  

any

##### <span id="get-repository-404"></span> 404 - repository not found
Status: Not Found

###### <span id="get-repository-404-schema"></span> Schema
   
  

any

##### <span id="get-repository-500"></span> 500 - failed to get repository
Status: Internal Server Error

###### <span id="get-repository-500-schema"></span> Schema
   
  

any

### <span id="get-secret"></span> Get a specific secret. (*getSecret*)

```
GET /rest/api/1/vault/secret/{name}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| name | `path` | string | `string` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-secret-200) | OK |  |  | [schema](#get-secret-200-schema) |
| [400](#get-secret-400) | Bad Request | invalid name |  | [schema](#get-secret-400-schema) |
| [404](#get-secret-404) | Not Found | secret not found |  | [schema](#get-secret-404-schema) |
| [500](#get-secret-500) | Internal Server Error | failed to load secrets |  | [schema](#get-secret-500-schema) |

#### Responses


##### <span id="get-secret-200"></span> 200
Status: OK

###### <span id="get-secret-200-schema"></span> Schema
   
  

[VaultSetting](#vault-setting)

##### <span id="get-secret-400"></span> 400 - invalid name
Status: Bad Request

###### <span id="get-secret-400-schema"></span> Schema
   
  

any

##### <span id="get-secret-404"></span> 404 - secret not found
Status: Not Found

###### <span id="get-secret-404-schema"></span> Schema
   
  

any

##### <span id="get-secret-500"></span> 500 - failed to load secrets
Status: Internal Server Error

###### <span id="get-secret-500-schema"></span> Schema
   
  

any

### <span id="get-token"></span> Creates a JWT token for a given api key pair. (*getToken*)

```
POST /rest/api/1/get-token
```

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-token-200) | OK | the generated JWT token |  | [schema](#get-token-200-schema) |
| [500](#get-token-500) | Internal Server Error | when there was a problem with matching the email, or the api key or generating the token |  | [schema](#get-token-500-schema) |

#### Responses


##### <span id="get-token-200"></span> 200 - the generated JWT token
Status: OK

###### <span id="get-token-200-schema"></span> Schema
   
  

any

##### <span id="get-token-500"></span> 500 - when there was a problem with matching the email, or the api key or generating the token
Status: Internal Server Error

###### <span id="get-token-500-schema"></span> Schema

### <span id="get-user"></span> Gets the user with the corresponding ID. (*getUser*)

```
GET /rest/api/1/user/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-user-200) | OK |  |  | [schema](#get-user-200-schema) |
| [400](#get-user-400) | Bad Request | invalid user id |  | [schema](#get-user-400-schema) |
| [404](#get-user-404) | Not Found | user not found |  | [schema](#get-user-404-schema) |
| [500](#get-user-500) | Internal Server Error | failed to get user |  | [schema](#get-user-500-schema) |

#### Responses


##### <span id="get-user-200"></span> 200
Status: OK

###### <span id="get-user-200-schema"></span> Schema
   
  

[User](#user)

##### <span id="get-user-400"></span> 400 - invalid user id
Status: Bad Request

###### <span id="get-user-400-schema"></span> Schema
   
  

any

##### <span id="get-user-404"></span> 404 - user not found
Status: Not Found

###### <span id="get-user-404-schema"></span> Schema
   
  

any

##### <span id="get-user-500"></span> 500 - failed to get user
Status: Internal Server Error

###### <span id="get-user-500-schema"></span> Schema
   
  

any

### <span id="hook-handler"></span> Handle the hooks created by the platform. (*hookHandler*)

```
POST /rest/api/1/hooks/{rid}/{vid}/callback
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| rid | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the repository. |
| vid | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the provider. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#hook-handler-200) | OK | success in case the hook event was processed without problems |  | [schema](#hook-handler-200-schema) |
| [400](#hook-handler-400) | Bad Request | for invalid parameters |  | [schema](#hook-handler-400-schema) |
| [404](#hook-handler-404) | Not Found | if the repository or the provider does not exist |  | [schema](#hook-handler-404-schema) |

#### Responses


##### <span id="hook-handler-200"></span> 200 - success in case the hook event was processed without problems
Status: OK

###### <span id="hook-handler-200-schema"></span> Schema

##### <span id="hook-handler-400"></span> 400 - for invalid parameters
Status: Bad Request

###### <span id="hook-handler-400-schema"></span> Schema

##### <span id="hook-handler-404"></span> 404 - if the repository or the provider does not exist
Status: Not Found

###### <span id="hook-handler-404-schema"></span> Schema

### <span id="list-api-keys"></span> Lists all api keys for a given user. (*listApiKeys*)

```
POST /rest/api/1/user/apikey
```

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-api-keys-200) | OK |  |  | [schema](#list-api-keys-200-schema) |
| [500](#list-api-keys-500) | Internal Server Error | failed to get user context |  | [schema](#list-api-keys-500-schema) |

#### Responses


##### <span id="list-api-keys-200"></span> 200
Status: OK

###### <span id="list-api-keys-200-schema"></span> Schema
   
  

[][APIKey](#api-key)

##### <span id="list-api-keys-500"></span> 500 - failed to get user context
Status: Internal Server Error

###### <span id="list-api-keys-500-schema"></span> Schema
   
  

any

### <span id="list-command-settings"></span> List settings for a command. (*listCommandSettings*)

```
POST /rest/api/1/command/{id}/settings
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the command to list settings for |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-command-settings-200) | OK |  |  | [schema](#list-command-settings-200-schema) |
| [400](#list-command-settings-400) | Bad Request | invalid id |  | [schema](#list-command-settings-400-schema) |
| [500](#list-command-settings-500) | Internal Server Error | failed to list settings |  | [schema](#list-command-settings-500-schema) |

#### Responses


##### <span id="list-command-settings-200"></span> 200
Status: OK

###### <span id="list-command-settings-200-schema"></span> Schema
   
  

[][CommandSetting](#command-setting)

##### <span id="list-command-settings-400"></span> 400 - invalid id
Status: Bad Request

###### <span id="list-command-settings-400-schema"></span> Schema
   
  

any

##### <span id="list-command-settings-500"></span> 500 - failed to list settings
Status: Internal Server Error

###### <span id="list-command-settings-500-schema"></span> Schema
   
  

any

### <span id="list-commands"></span> list commands (*listCommands*)

```
POST /rest/api/1/commands
```

List commands

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| listOptions | `body` | [ListOptions](#list-options) | `models.ListOptions` | |  | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-commands-200) | OK |  |  | [schema](#list-commands-200-schema) |
| [500](#list-commands-500) | Internal Server Error | failed to get user context |  | [schema](#list-commands-500-schema) |

#### Responses


##### <span id="list-commands-200"></span> 200
Status: OK

###### <span id="list-commands-200-schema"></span> Schema
   
  

[][Command](#command)

##### <span id="list-commands-500"></span> 500 - failed to get user context
Status: Internal Server Error

###### <span id="list-commands-500-schema"></span> Schema
   
  

any

### <span id="list-events"></span> List events for a repository. (*listEvents*)

```
POST /rest/api/1/events/{repoid}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| repoid | `path` | int (formatted integer) | `int64` |  | ✓ |  | The ID of the repository to list events for. |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-events-200) | OK |  |  | [schema](#list-events-200-schema) |
| [400](#list-events-400) | Bad Request | invalid repository id |  | [schema](#list-events-400-schema) |
| [500](#list-events-500) | Internal Server Error | failed to list events |  | [schema](#list-events-500-schema) |

#### Responses


##### <span id="list-events-200"></span> 200
Status: OK

###### <span id="list-events-200-schema"></span> Schema
   
  

[][Event](#event)

##### <span id="list-events-400"></span> 400 - invalid repository id
Status: Bad Request

###### <span id="list-events-400-schema"></span> Schema
   
  

any

##### <span id="list-events-500"></span> 500 - failed to list events
Status: Internal Server Error

###### <span id="list-events-500-schema"></span> Schema
   
  

any

### <span id="list-repositories"></span> list repositories (*listRepositories*)

```
POST /rest/api/1/repositories
```

List repositories

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| listOptions | `body` | [ListOptions](#list-options) | `models.ListOptions` | |  | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-repositories-200) | OK |  |  | [schema](#list-repositories-200-schema) |
| [500](#list-repositories-500) | Internal Server Error | failed to list repositories |  | [schema](#list-repositories-500-schema) |

#### Responses


##### <span id="list-repositories-200"></span> 200
Status: OK

###### <span id="list-repositories-200-schema"></span> Schema
   
  

[][Repository](#repository)

##### <span id="list-repositories-500"></span> 500 - failed to list repositories
Status: Internal Server Error

###### <span id="list-repositories-500-schema"></span> Schema
   
  

any

### <span id="list-secrets"></span> List all settings without the values. (*listSecrets*)

```
POST /rest/api/1/vault/secrets
```

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-secrets-200) | OK |  |  | [schema](#list-secrets-200-schema) |
| [500](#list-secrets-500) | Internal Server Error | failed to load secrets |  | [schema](#list-secrets-500-schema) |

#### Responses


##### <span id="list-secrets-200"></span> 200
Status: OK

###### <span id="list-secrets-200-schema"></span> Schema
   
  

[][VaultSetting](#vault-setting)

##### <span id="list-secrets-500"></span> 500 - failed to load secrets
Status: Internal Server Error

###### <span id="list-secrets-500-schema"></span> Schema
   
  

any

### <span id="list-supported-platforms"></span> Lists all supported platforms. (*listSupportedPlatforms*)

```
GET /rest/api/1/supported-platforms
```

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-supported-platforms-200) | OK | the list of supported platform ids |  | [schema](#list-supported-platforms-200-schema) |

#### Responses


##### <span id="list-supported-platforms-200"></span> 200 - the list of supported platform ids
Status: OK

###### <span id="list-supported-platforms-200-schema"></span> Schema
   
  

[][Platform](#platform)

### <span id="list-users"></span> list users (*listUsers*)

```
POST /rest/api/1/users
```

List users

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-users-200) | OK |  |  | [schema](#list-users-200-schema) |
| [500](#list-users-500) | Internal Server Error | failed to list user |  | [schema](#list-users-500-schema) |

#### Responses


##### <span id="list-users-200"></span> 200
Status: OK

###### <span id="list-users-200-schema"></span> Schema
   
  

[][User](#user)

##### <span id="list-users-500"></span> 500 - failed to list user
Status: Internal Server Error

###### <span id="list-users-500-schema"></span> Schema
   
  

any

### <span id="refresh-token"></span> Refresh the authentication token. (*refreshToken*)

```
POST /rest/api/1/auth/refresh
```

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#refresh-token-200) | OK | Status OK |  | [schema](#refresh-token-200-schema) |
| [401](#refresh-token-401) | Unauthorized | refresh token cookie not found|error refreshing token |  | [schema](#refresh-token-401-schema) |

#### Responses


##### <span id="refresh-token-200"></span> 200 - Status OK
Status: OK

###### <span id="refresh-token-200-schema"></span> Schema

##### <span id="refresh-token-401"></span> 401 - refresh token cookie not found|error refreshing token
Status: Unauthorized

###### <span id="refresh-token-401-schema"></span> Schema

### <span id="remove-command-rel-for-platform-command"></span> Remove a relationship to a platform. This command will no longer be running for that platform events. (*removeCommandRelForPlatformCommand*)

```
POST /rest/api/1/command/remove-command-rel-for-platform/{cmdid}/{repoid}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| cmdid | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |
| repoid | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#remove-command-rel-for-platform-command-200) | OK | successfully removed relationship |  | [schema](#remove-command-rel-for-platform-command-200-schema) |
| [400](#remove-command-rel-for-platform-command-400) | Bad Request | invalid ids or platform not found |  | [schema](#remove-command-rel-for-platform-command-400-schema) |
| [500](#remove-command-rel-for-platform-command-500) | Internal Server Error | failed to add relationship |  | [schema](#remove-command-rel-for-platform-command-500-schema) |

#### Responses


##### <span id="remove-command-rel-for-platform-command-200"></span> 200 - successfully removed relationship
Status: OK

###### <span id="remove-command-rel-for-platform-command-200-schema"></span> Schema

##### <span id="remove-command-rel-for-platform-command-400"></span> 400 - invalid ids or platform not found
Status: Bad Request

###### <span id="remove-command-rel-for-platform-command-400-schema"></span> Schema
   
  

any

##### <span id="remove-command-rel-for-platform-command-500"></span> 500 - failed to add relationship
Status: Internal Server Error

###### <span id="remove-command-rel-for-platform-command-500-schema"></span> Schema
   
  

any

### <span id="remove-command-rel-for-repository-command"></span> Remove a relationship to a repository. This command will no longer be running for that repository events. (*removeCommandRelForRepositoryCommand*)

```
POST /rest/api/1/command/remove-command-rel-for-repository/{cmdid}/{repoid}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| cmdid | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |
| repoid | `path` | int (formatted integer) | `int64` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#remove-command-rel-for-repository-command-200) | OK | successfully removed relationship |  | [schema](#remove-command-rel-for-repository-command-200-schema) |
| [400](#remove-command-rel-for-repository-command-400) | Bad Request | invalid ids or repositroy not found |  | [schema](#remove-command-rel-for-repository-command-400-schema) |
| [500](#remove-command-rel-for-repository-command-500) | Internal Server Error | failed to add relationship |  | [schema](#remove-command-rel-for-repository-command-500-schema) |

#### Responses


##### <span id="remove-command-rel-for-repository-command-200"></span> 200 - successfully removed relationship
Status: OK

###### <span id="remove-command-rel-for-repository-command-200-schema"></span> Schema

##### <span id="remove-command-rel-for-repository-command-400"></span> 400 - invalid ids or repositroy not found
Status: Bad Request

###### <span id="remove-command-rel-for-repository-command-400-schema"></span> Schema
   
  

any

##### <span id="remove-command-rel-for-repository-command-500"></span> 500 - failed to add relationship
Status: Internal Server Error

###### <span id="remove-command-rel-for-repository-command-500-schema"></span> Schema
   
  

any

### <span id="update-command"></span> Updates a given command. (*updateCommand*)

```
POST /rest/api/1/command/update
```

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| command | `body` | [Command](#command) | `models.Command` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#update-command-200) | OK | successfully updated command |  | [schema](#update-command-200-schema) |
| [400](#update-command-400) | Bad Request | binding error |  | [schema](#update-command-400-schema) |
| [500](#update-command-500) | Internal Server Error | failed to update the command |  | [schema](#update-command-500-schema) |

#### Responses


##### <span id="update-command-200"></span> 200 - successfully updated command
Status: OK

###### <span id="update-command-200-schema"></span> Schema
   
  

[Command](#command)

##### <span id="update-command-400"></span> 400 - binding error
Status: Bad Request

###### <span id="update-command-400-schema"></span> Schema
   
  

any

##### <span id="update-command-500"></span> 500 - failed to update the command
Status: Internal Server Error

###### <span id="update-command-500-schema"></span> Schema
   
  

any

### <span id="update-command-setting"></span> Create a new command setting. (*updateCommandSetting*)

```
POST /rest/api/1/command/settings/update
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| setting | `body` | [CommandSetting](#command-setting) | `models.CommandSetting` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#update-command-setting-200) | OK | successfully created command setting |  | [schema](#update-command-setting-200-schema) |
| [400](#update-command-setting-400) | Bad Request | binding error |  | [schema](#update-command-setting-400-schema) |
| [500](#update-command-setting-500) | Internal Server Error | failed to create the command setting |  | [schema](#update-command-setting-500-schema) |

#### Responses


##### <span id="update-command-setting-200"></span> 200 - successfully created command setting
Status: OK

###### <span id="update-command-setting-200-schema"></span> Schema

##### <span id="update-command-setting-400"></span> 400 - binding error
Status: Bad Request

###### <span id="update-command-setting-400-schema"></span> Schema
   
  

any

##### <span id="update-command-setting-500"></span> 500 - failed to create the command setting
Status: Internal Server Error

###### <span id="update-command-setting-500-schema"></span> Schema
   
  

any

### <span id="update-repository"></span> Updates an existing repository. (*updateRepository*)

```
POST /rest/api/1/repository/update
```

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| repository | `body` | [Repository](#repository) | `models.Repository` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#update-repository-200) | OK | the updated repository |  | [schema](#update-repository-200-schema) |
| [400](#update-repository-400) | Bad Request | failed to bind repository |  | [schema](#update-repository-400-schema) |
| [404](#update-repository-404) | Not Found | repository not found |  | [schema](#update-repository-404-schema) |
| [500](#update-repository-500) | Internal Server Error | failed to update repository |  | [schema](#update-repository-500-schema) |

#### Responses


##### <span id="update-repository-200"></span> 200 - the updated repository
Status: OK

###### <span id="update-repository-200-schema"></span> Schema
   
  

[Repository](#repository)

##### <span id="update-repository-400"></span> 400 - failed to bind repository
Status: Bad Request

###### <span id="update-repository-400-schema"></span> Schema
   
  

any

##### <span id="update-repository-404"></span> 404 - repository not found
Status: Not Found

###### <span id="update-repository-404-schema"></span> Schema
   
  

any

##### <span id="update-repository-500"></span> 500 - failed to update repository
Status: Internal Server Error

###### <span id="update-repository-500-schema"></span> Schema
   
  

any

### <span id="update-secret"></span> Updates an existing secret. (*updateSecret*)

```
POST /rest/api/1/vault/secret/update
```

#### Consumes
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| secret | `body` | [VaultSetting](#vault-setting) | `models.VaultSetting` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#update-secret-200) | OK | OK setting successfully updated |  | [schema](#update-secret-200-schema) |
| [400](#update-secret-400) | Bad Request | invalid json payload |  | [schema](#update-secret-400-schema) |
| [404](#update-secret-404) | Not Found | setting not found |  | [schema](#update-secret-404-schema) |
| [500](#update-secret-500) | Internal Server Error | failed to update secret |  | [schema](#update-secret-500-schema) |

#### Responses


##### <span id="update-secret-200"></span> 200 - OK setting successfully updated
Status: OK

###### <span id="update-secret-200-schema"></span> Schema

##### <span id="update-secret-400"></span> 400 - invalid json payload
Status: Bad Request

###### <span id="update-secret-400-schema"></span> Schema
   
  

any

##### <span id="update-secret-404"></span> 404 - setting not found
Status: Not Found

###### <span id="update-secret-404-schema"></span> Schema
   
  

any

##### <span id="update-secret-500"></span> 500 - failed to update secret
Status: Internal Server Error

###### <span id="update-secret-500-schema"></span> Schema
   
  

any

### <span id="update-user"></span> Updates an existing user. (*updateUser*)

```
POST /rest/api/1/user/update
```

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| user | `body` | [User](#user) | `models.User` | | ✓ | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#update-user-200) | OK | user successfully updated |  | [schema](#update-user-200-schema) |
| [400](#update-user-400) | Bad Request | invalid json payload |  | [schema](#update-user-400-schema) |
| [404](#update-user-404) | Not Found | user not found |  | [schema](#update-user-404-schema) |
| [500](#update-user-500) | Internal Server Error | failed to update user |  | [schema](#update-user-500-schema) |

#### Responses


##### <span id="update-user-200"></span> 200 - user successfully updated
Status: OK

###### <span id="update-user-200-schema"></span> Schema
   
  

[User](#user)

##### <span id="update-user-400"></span> 400 - invalid json payload
Status: Bad Request

###### <span id="update-user-400-schema"></span> Schema
   
  

any

##### <span id="update-user-404"></span> 404 - user not found
Status: Not Found

###### <span id="update-user-404-schema"></span> Schema
   
  

any

##### <span id="update-user-500"></span> 500 - failed to update user
Status: Internal Server Error

###### <span id="update-user-500-schema"></span> Schema
   
  

any

### <span id="upload-command"></span> Upload a command. To set up anything for the command, like schedules etc, (*uploadCommand*)

```
POST /rest/api/1/command
```

the command has to be edited. We don't support uploading the same thing twice.
If the command binary needs to be updated, delete the command and upload the
new binary.

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#upload-command-201) | Created | in case of successful file upload |  | [schema](#upload-command-201-schema) |
| [400](#upload-command-400) | Bad Request | invalid file format or command already exists |  | [schema](#upload-command-400-schema) |
| [500](#upload-command-500) | Internal Server Error | failed to upload file, create plugin, create command or copy operations |  | [schema](#upload-command-500-schema) |

#### Responses


##### <span id="upload-command-201"></span> 201 - in case of successful file upload
Status: Created

###### <span id="upload-command-201-schema"></span> Schema
   
  

[Command](#command)

##### <span id="upload-command-400"></span> 400 - invalid file format or command already exists
Status: Bad Request

###### <span id="upload-command-400-schema"></span> Schema
   
  

any

##### <span id="upload-command-500"></span> 500 - failed to upload file, create plugin, create command or copy operations
Status: Internal Server Error

###### <span id="upload-command-500-schema"></span> Schema
   
  

any

### <span id="user-callback"></span> This is the url to which Google calls back after a successful login. (*userCallback*)

```
GET /rest/api/1/auth/callback
```

Creates a cookie which will hold the authenticated user.

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| code | `query` | string | `string` |  | ✓ |  | the authorization code provided by Google |
| state | `query` | string | `string` |  | ✓ |  | the state variable defined by Google |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [308](#user-callback-308) | Permanent Redirect | the permanent redirect url |  | [schema](#user-callback-308-schema) |
| [401](#user-callback-401) | Unauthorized | error verifying state | error during token exchange |  | [schema](#user-callback-401-schema) |
| [404](#user-callback-404) | Not Found | error invalid state|code |  | [schema](#user-callback-404-schema) |

#### Responses


##### <span id="user-callback-308"></span> 308 - the permanent redirect url
Status: Permanent Redirect

###### <span id="user-callback-308-schema"></span> Schema

##### <span id="user-callback-401"></span> 401 - error verifying state | error during token exchange
Status: Unauthorized

###### <span id="user-callback-401-schema"></span> Schema

##### <span id="user-callback-404"></span> 404 - error invalid state|code
Status: Not Found

###### <span id="user-callback-404-schema"></span> Schema

### <span id="user-login"></span> User login. (*userLogin*)

```
GET /rest/api/1/auth/login
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| redirect_url | `query` | string | `string` |  | ✓ |  | the redirect URL coming from Google to redirect login to |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [307](#user-login-307) | Temporary Redirect | the redirect url to the login |  | [schema](#user-login-307-schema) |
| [401](#user-login-401) | Unauthorized | error generating state |  | [schema](#user-login-401-schema) |
| [404](#user-login-404) | Not Found | error invalid redirect_url |  | [schema](#user-login-404-schema) |

#### Responses


##### <span id="user-login-307"></span> 307 - the redirect url to the login
Status: Temporary Redirect

###### <span id="user-login-307-schema"></span> Schema

##### <span id="user-login-401"></span> 401 - error generating state
Status: Unauthorized

###### <span id="user-login-401-schema"></span> Schema

##### <span id="user-login-404"></span> 404 - error invalid redirect_url
Status: Not Found

###### <span id="user-login-404-schema"></span> Schema

## Models

### <span id="api-key"></span> APIKey


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIKeyID | string| `string` | ✓ | | APIKeyID is a generated id of the key. |  |
| APIKeySecret | string| `string` | ✓ | | APIKeySecret is a generated secret, aka, the key. |  |
| CreateAt | date-time (formatted string)| `strfmt.DateTime` | ✓ | | CreateAt defines when this key was created. | `time.Now()` |
| ID | int64 (formatted integer)| `int64` | ✓ | | ID of the key. This is auto-generated. |  |
| Name | string| `string` | ✓ | | Name of the key |  |
| TTL | string| `string` | ✓ | | TTL defines how long this key can live in duration. | `1h10m10s` |
| UserID | int64 (formatted integer)| `int64` | ✓ | | UserID is the ID of the user to which this key belongs. |  |



### <span id="auth"></span> Auth


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Password | string| `string` |  | | Password is the password required to access the platform for this repositroy. |  |
| SSH | string| `string` |  | | SSH private key. |  |
| Secret | string| `string` | ✓ | | Hook secret to create a hook with on the respective platform. |  |
| Username | string| `string` |  | | Username is the username required to access the platform for this repositroy. |  |



### <span id="command"></span> Command


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Enabled | boolean| `bool` |  | | Enabled defines if this command can be executed or not. | `false` |
| Filename | string| `string` | ✓ | | Filename is the name of the file which holds this command. | `my_awesome_command` |
| Hash | string| `string` | ✓ | | Hash is the hash of the command file. |  |
| ID | int64 (formatted integer)| `int64` | ✓ | | ID of the command. Generated. |  |
| Location | string| `string` | ✓ | | Location is where this command is located at. This is the full path of the containing folder. | `/tmp/krok-commands` |
| Name | string| `string` | ✓ | | Name of the command. |  |
| Repositories | [][Repository](#repository)| `[]*Repository` |  | | Repositories that this command can execute on. |  |
| Schedule | string| `string` |  | | Schedule of the command. | `0 * * * * // follows cron job syntax.` |



### <span id="command-run"></span> CommandRun


> CommandRun is a single run of a command belonging to an event
including things like, state, event, and created at.
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| CommandName | string| `string` | ✓ | | CommandName is the name of the command that is being executed. |  |
| CreateAt | date-time (formatted string)| `strfmt.DateTime` | ✓ | | CreatedAt is the time when this command run was created. |  |
| EventID | int64 (formatted integer)| `int64` | ✓ | | EventID is the ID of the event that this run belongs to. |  |
| ID | int64 (formatted integer)| `int64` | ✓ | | ID is a generatd identifier. |  |
| Outcome | string| `string` |  | | Outcome is any output of the command. Stdout and stderr combined. |  |
| Status | string| `string` | ✓ | | Status is the current state of the command run. | `running, failed, success` |



### <span id="command-setting"></span> CommandSetting


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| CommandID | int64 (formatted integer)| `int64` | ✓ | | CommandID is the ID of the command to which these settings belong to. |  |
| ID | int64 (formatted integer)| `int64` | ✓ | | ID is a generated ID. |  |
| InVault | boolean| `bool` |  | | InVault defines if this is sensitive information and should be stored securely. |  |
| Key | string| `string` | ✓ | | Key is the name of the setting. |  |
| Value | string| `string` | ✓ | | Value is the value of the setting. |  |



### <span id="event"></span> Event


> Event contains details about a platform event, such as
the repository it belongs to and the event that created it...
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| CommandRuns | [][CommandRun](#command-run)| `[]*CommandRun` |  | | CommandRuns contains a list of CommandRuns which executed for this event. |  |
| CreateAt | date-time (formatted string)| `strfmt.DateTime` | ✓ | | CreatedAt contains the timestamp when this event occurred. |  |
| EventID | string| `string` | ✓ | | EvenID is the ID of the corresponding event on the given platform. If that cannot be determined
an ID is generated. |  |
| ID | int64 (formatted integer)| `int64` | ✓ | | ID is a generated ID. |  |
| Payload | string| `string` | ✓ | | Payload defines the information received from the platform for this event. |  |
| RepositoryID | int64 (formatted integer)| `int64` | ✓ | | RepositoryID contains the ID of the repository for which this event occurred. |  |
| VCS | int64 (formatted integer)| `int64` | ✓ | | VCS is the ID of the platform on which this even occurred. |  |



### <span id="git-lab"></span> GitLab


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ProjectID | int64 (formatted integer)| `int64` |  | | ProjectID is an optional ID which defines a project in Gitlab. |  |



### <span id="list-options"></span> ListOptions


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| EndDate | date-time (formatted string)| `strfmt.DateTime` |  | | EndDate defines a date of end to look for events. Not Inclusive. | `2021-02-03` |
| Name | string| `string` |  | | Name of the context for which this option is used. | `\"partialNameOfACommand\` |
| Page | int64 (formatted integer)| `int64` |  | | Current Page | `0` |
| PageSize | int64 (formatted integer)| `int64` |  | | Items per Page

required false | `10` |
| StartingDate | date-time (formatted string)| `strfmt.DateTime` |  | | StartingDate defines a date of start to look for events. Inclusive. | `2021-02-02` |
| VCS | int64 (formatted integer)| `int64` |  | | Only list all entries for a given platform ID. |  |



### <span id="platform"></span> Platform


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ID | int64 (formatted integer)| `int64` | ✓ | | ID of the platform. This is choosen. |  |
| Name | string| `string` | ✓ | | Name of the platform. | `github, gitlab, gitea` |



### <span id="repository"></span> Repository


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Commands | [][Command](#command)| `[]*Command` |  | | Commands contains all the commands which this repository is attached to. |  |
| Events | []string| `[]string` |  | | TODO: Think about storing this |  |
| ID | int64 (formatted integer)| `int64` | ✓ | | ID of the repository. Auto-generated. |  |
| Name | string| `string` |  | | Name of the repository. |  |
| URL | string| `string` | ✓ | | URL of the repository. |  |
| UniqueURL | string| `string` | ✓ | | This field is not saved in the DB but generated every time the repository
details needs to be displayed. |  |
| VCS | int64 (formatted integer)| `int64` | ✓ | | VCS Defines which handler will be used. For values, see platforms.go. |  |
| auth | [Auth](#auth)| `Auth` | ✓ | |  |  |
| git_lab | [GitLab](#git-lab)| `GitLab` |  | |  |  |



### <span id="user"></span> User


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIKeys | [][APIKey](#api-key)| `[]*APIKey` |  | | APIKeys contains generated api access keys for this user. |  |
| DisplayName | string| `string` |  | | DisplayName is the name of the user. |  |
| Email | string| `string` | ✓ | | Email of the user. |  |
| ID | int64 (formatted integer)| `int64` | ✓ | | ID of the user. This is auto-generated. |  |
| LastLogin | date-time (formatted string)| `strfmt.DateTime` | ✓ | | LastLogin contains the timestamp of the last successful login of this user. |  |



### <span id="user-auth-details"></span> UserAuthDetails


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Email | string| `string` | ✓ | | Email is the email of the registered user. |  |
| FirstName | string| `string` | ✓ | | FirstName is the first name of the user. |  |
| LastName | string| `string` | ✓ | | LastName is the last name of the user. |  |



### <span id="v-c-s-token"></span> VCSToken


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Token | string| `string` | ✓ | | Token is the actual token. |  |
| VCS | int64 (formatted integer)| `int64` | ✓ | | VCS is the ID of the platform to which this token belongs to. |  |



### <span id="vault-setting"></span> VaultSetting


> VaultSetting defines a setting that comes from the vault
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Key | string| `string` | ✓ | | Key is the name of the setting. |  |
| Value | string| `string` | ✓ | | Value is the value of the setting. |  |


