


# mwallet.
Documentation of mwallet API.
  

## Informations

### Version

1.0.0

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  accounts

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /opening/accounts | [create account](#create-account) | Create new account. |
| DELETE | /opening/accounts/{id} | [delete account](#delete-account) | Delete an account. |
| GET | /opening/accounts/{id} | [find account](#find-account) | Find account. |
| GET | /opening/accounts | [list accounts](#list-accounts) | List all accounts. |
  


###  payments

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /transferring/payments/{id} | [find payment](#find-payment) | Find payment. |
| GET | /transferring/payments | [list payments](#list-payments) | List all payments. |
| POST | /transferring/payments | [send payment](#send-payment) | Transfer payment. |
  


## Paths

### <span id="create-account"></span> Create new account. (*createAccount*)

```
POST /opening/accounts
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| create account | `body` | [AddAccount](#add-account) | `models.AddAccount` | | ✓ | | account payload |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#create-account-200) | OK | Created |  | [schema](#create-account-200-schema) |
| [500](#create-account-500) | Internal Server Error | Internal server error |  | [schema](#create-account-500-schema) |

#### Responses


##### <span id="create-account-200"></span> 200 - Created
Status: OK

###### <span id="create-account-200-schema"></span> Schema
   
  

[AddAccount](#add-account)

##### <span id="create-account-500"></span> 500 - Internal server error
Status: Internal Server Error

###### <span id="create-account-500-schema"></span> Schema

### <span id="delete-account"></span> Delete an account. (*deleteAccount*)

```
DELETE /opening/accounts/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | string | `string` |  | ✓ |  | id of the account |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-account-200) | OK |  |  | [schema](#delete-account-200-schema) |
| [500](#delete-account-500) | Internal Server Error | Internal server error |  | [schema](#delete-account-500-schema) |

#### Responses


##### <span id="delete-account-200"></span> 200
Status: OK

###### <span id="delete-account-200-schema"></span> Schema
   
  

[DeleteAccountResponse](#delete-account-response)

##### <span id="delete-account-500"></span> 500 - Internal server error
Status: Internal Server Error

###### <span id="delete-account-500-schema"></span> Schema

### <span id="find-account"></span> Find account. (*findAccount*)

```
GET /opening/accounts/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | string | `string` |  | ✓ |  | id of the account |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#find-account-200) | OK |  |  | [schema](#find-account-200-schema) |
| [500](#find-account-500) | Internal Server Error | Internal server error |  | [schema](#find-account-500-schema) |

#### Responses


##### <span id="find-account-200"></span> 200
Status: OK

###### <span id="find-account-200-schema"></span> Schema
   
  

[AddAccount](#add-account)

##### <span id="find-account-500"></span> 500 - Internal server error
Status: Internal Server Error

###### <span id="find-account-500-schema"></span> Schema

### <span id="find-payment"></span> Find payment. (*findPayment*)

```
GET /transferring/payments/{id}
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | string | `string` |  | ✓ |  | id of the account |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#find-payment-200) | OK |  |  | [schema](#find-payment-200-schema) |
| [500](#find-payment-500) | Internal Server Error | Internal server error |  | [schema](#find-payment-500-schema) |

#### Responses


##### <span id="find-payment-200"></span> 200
Status: OK

###### <span id="find-payment-200-schema"></span> Schema
   
  

[FindPaymentResponse](#find-payment-response)

##### <span id="find-payment-500"></span> 500 - Internal server error
Status: Internal Server Error

###### <span id="find-payment-500-schema"></span> Schema

### <span id="list-accounts"></span> List all accounts. (*listAccounts*)

```
GET /opening/accounts
```

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-accounts-200) | OK |  |  | [schema](#list-accounts-200-schema) |
| [500](#list-accounts-500) | Internal Server Error | Internal server error |  | [schema](#list-accounts-500-schema) |

#### Responses


##### <span id="list-accounts-200"></span> 200
Status: OK

###### <span id="list-accounts-200-schema"></span> Schema
   
  

[ListAccountsResponse](#list-accounts-response)

##### <span id="list-accounts-500"></span> 500 - Internal server error
Status: Internal Server Error

###### <span id="list-accounts-500-schema"></span> Schema

### <span id="list-payments"></span> List all payments. (*listPayments*)

```
GET /transferring/payments
```

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-payments-200) | OK |  |  | [schema](#list-payments-200-schema) |
| [500](#list-payments-500) | Internal Server Error | Internal server error |  | [schema](#list-payments-500-schema) |

#### Responses


##### <span id="list-payments-200"></span> 200
Status: OK

###### <span id="list-payments-200-schema"></span> Schema
   
  

[ListPaymentResponse](#list-payment-response)

##### <span id="list-payments-500"></span> 500 - Internal server error
Status: Internal Server Error

###### <span id="list-payments-500-schema"></span> Schema

### <span id="send-payment"></span> Transfer payment. (*sendPayment*)

```
POST /transferring/payments
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| transfer payment | `body` | [SendPayment](#send-payment) | `models.SendPayment` | | ✓ | | payment payload |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#send-payment-200) | OK |  |  | [schema](#send-payment-200-schema) |
| [500](#send-payment-500) | Internal Server Error | Internal server error |  | [schema](#send-payment-500-schema) |

#### Responses


##### <span id="send-payment-200"></span> 200
Status: OK

###### <span id="send-payment-200-schema"></span> Schema
   
  

[SendPayment](#send-payment)

##### <span id="send-payment-500"></span> 500 - Internal server error
Status: Internal Server Error

###### <span id="send-payment-500-schema"></span> Schema

## Models

### <span id="account"></span> Account


> Account represents a mobile wallet account
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Balance | double (formatted number)| `float64` |  | |  |  |
| Currency | string| `string` |  | |  |  |
| ID | string| `string` |  | |  |  |



### <span id="find-payment-response"></span> FindPaymentResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Err | string| `string` |  | |  |  |
| Payments | [][Payment](#payment)| `[]*Payment` |  | |  |  |



### <span id="payment"></span> Payment


> Payment represents a transfer payment
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Account | string| `string` |  | |  |  |
| Amount | double (formatted number)| `float64` |  | |  |  |
| Direction | string| `string` |  | |  |  |
| FromAccount | string| `string` |  | |  |  |
| ID | string| `string` |  | |  |  |
| ToAccount | string| `string` |  | |  |  |



### <span id="add-account"></span> addAccount


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Balance | double (formatted number)| `float64` | ✓ | |  | `100` |
| Currency | string| `string` | ✓ | |  | `USD` |
| ID | string| `string` | ✓ | |  | `bob123` |



### <span id="delete-account-request"></span> deleteAccountRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ID | string| `string` | ✓ | |  | `bob123` |



### <span id="delete-account-response"></span> deleteAccountResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Err | string| `string` |  | |  |  |



### <span id="find-payment-request"></span> findPaymentRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| AccountID | string| `string` | ✓ | |  | `bob123` |



### <span id="get-account-request"></span> getAccountRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ID | string| `string` | ✓ | |  | `bob123` |



### <span id="list-accounts-response"></span> listAccountsResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Accounts | [][Account](#account)| `[]*Account` |  | |  |  |
| Err | string| `string` |  | |  |  |



### <span id="list-payment-response"></span> listPaymentResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Err | string| `string` |  | |  |  |
| Payments | [][Payment](#payment)| `[]*Payment` |  | |  |  |



### <span id="send-payment"></span> sendPayment


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Amount | double (formatted number)| `float64` | ✓ | |  | `50` |
| FromAccount | string| `string` | ✓ | |  | `bob123` |
| ToAccount | string| `string` | ✓ | |  | `alice456` |


