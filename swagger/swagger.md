


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
| POST | /rest/api/1/user/apikey/generate/{name} | [create Api key](#create-api-key) | Creates an api key pair for a given user. |
| DELETE | /rest/api/1/user/apikey/delete/{keyid} | [delete Api key](#delete-api-key) | Deletes a set of api keys for a given user with a given id. |
| DELETE | /rest/api/1/command/{id} | [delete command](#delete-command) | Deletes given command. |
| GET | /rest/api/1/user/apikey/{keyid} | [get Api keys](#get-api-keys) | Returns a given api key. |
| GET | /rest/api/1/command/{id} | [get command](#get-command) | Returns a specific command. |
| POST | /rest/api/1/get-token | [get token](#get-token) | Creates a JWT token for a given api key pair. |
| POST | /rest/api/1/hooks/{rid}/{vid}/callback | [hook handler](#hook-handler) | Handle the hooks created by the platform. |
| POST | /rest/api/1/user/apikey | [list Api keys](#list-api-keys) | Lists all api keys for a given user. |
| POST | /rest/api/1/commands | [list commands](#list-commands) |  |
| GET | /rest/api/1/supported-platforms | [list supported platforms](#list-supported-platforms) | Lists all supported platforms. |
| POST | /rest/api/1/auth/refresh | [refresh token](#refresh-token) | Refresh the authentication token. |
| POST | /rest/api/1/command | [upload command](#upload-command) | Upload a command. To set up anything for the command, like schedules etc, |
| GET | /rest/api/1/auth/callback | [user callback](#user-callback) | This is the url to which Google calls back after a successful login. |
| GET | /rest/api/1/auth/login | [user login](#user-login) | User login. |
  


## Paths

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

##### <span id="create-api-key-500"></span> 500 - when failed to get user context
Status: Internal Server Error

###### <span id="create-api-key-500-schema"></span> Schema

### <span id="delete-api-key"></span> Deletes a set of api keys for a given user with a given id. (*deleteApiKey*)

```
DELETE /rest/api/1/user/apikey/delete/{keyid}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| keyid | `path` | string | `string` |  | ✓ |  | The ID of the key to delete |

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

##### <span id="delete-api-key-500"></span> 500 - when the deletion operation failed
Status: Internal Server Error

###### <span id="delete-api-key-500-schema"></span> Schema

### <span id="delete-command"></span> Deletes given command. (*deleteCommand*)

```
DELETE /rest/api/1/command/{id}
```

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | string | `string` |  | ✓ |  | The ID of the command to delete |

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

##### <span id="delete-command-500"></span> 500 - when the deletion operation failed
Status: Internal Server Error

###### <span id="delete-command-500-schema"></span> Schema

### <span id="get-api-keys"></span> Returns a given api key. (*getApiKeys*)

```
GET /rest/api/1/user/apikey/{keyid}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| keyid | `path` | string | `string` |  | ✓ |  | The ID of the key to return |

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

### <span id="get-command"></span> Returns a specific command. (*getCommand*)

```
GET /rest/api/1/command/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | string | `string` |  | ✓ |  |  |

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

##### <span id="get-command-500"></span> 500 - failed to get user context
Status: Internal Server Error

###### <span id="get-command-500-schema"></span> Schema

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

### <span id="hook-handler"></span> Handle the hooks created by the platform. (*hookHandler*)

```
POST /rest/api/1/hooks/{rid}/{vid}/callback
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| rid | `path` | string | `string` |  | ✓ |  | The ID of the repository. |
| vid | `path` | string | `string` |  | ✓ |  | The ID of the provider. |

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

##### <span id="upload-command-500"></span> 500 - failed to upload file, create plugin, create command or copy operations
Status: Internal Server Error

###### <span id="upload-command-500-schema"></span> Schema

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
| ID | int64 (formatted integer)| `int64` | ✓ | | ID of the key. This is auto-generated. |  |
| Name | string| `string` | ✓ | | Name of the key |  |
| TTL | date-time (formatted string)| `strfmt.DateTime` | ✓ | | TTL defines how long this key can live. | `time.Now().Add(10 * time.Minute)` |
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



### <span id="new-user"></span> NewUser


> NewUser is a new user in the Krok system. Specifically this exposes the token and should only be used when creating
a user for the first time.
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIKeys | [][APIKey](#api-key)| `[]*APIKey` |  | | APIKeys contains generated api access keys for this user. |  |
| DisplayName | string| `string` |  | | DisplayName is the name of the user. |  |
| Email | string| `string` | ✓ | | Email of the user. |  |
| ID | int64 (formatted integer)| `int64` | ✓ | | ID of the user. This is auto-generated. |  |
| LastLogin | date-time (formatted string)| `strfmt.DateTime` | ✓ | | LastLogin contains the timestamp of the last successful login of this user. |  |
| Token | string| `string` | ✓ | | Token is displayed once for new users. Then never again. |  |



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


