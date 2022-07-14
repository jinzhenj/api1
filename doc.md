# sem comments

## `@omitempty`

used for: `StructField`

Field won't return if it is empty.

Example:

```
struct User {
  name: string

  # @omitempty
  address: string
}
```

## `@ignore`

used for: `StructField`

Field won't return.

Example:

```
struct User {
  name: string

  # @ignore
  password: string
}
```

## `@route`

used for: `Fun`

Function is a rest route handler.

Example:

```
interface user {

  # @route get /users/:id
  getUser(id: int): User
}
```

## `@go.type`

used for: `Scalar`, `StructField`, `Param`, `Fun`

Specify golang type for scalar or typeRef.

Example:

```
# @go.type uint64
scalar Timestamp

struct User {
  
  # @go.type map[string]string
  properties: object
}

interface user {

  # @route get /users/:id/properties
  # @go.type map[string]string
  getProperties(id: int): object
 
  # @route put /users/:id/properties
  setProperties(
    id: int,

    # @go.type map[string]string
    properties: object
  )
}
```

## `@go.middleware`

used for: `Fun`

Add middleware for golang route handler

Example:

```
interface TypeController {

  # @route get /types
  listTypes(): [Type]

  # @route put /types/:id
  # @go.middleware adminRequired
  modifyType(id: int, type: Type)
}
```

## `@go.validator`

used for: `StructField`, `Param`

Add validator for field or param.

Example:

```
struct UserLoginRequest {
  # @go.validator: email
  email: string
  password: string
}

interface LoginController {

  # @route post /actions/checkUserExists
  checkUserExists(
    # @go.validator email
    email: string
  ): boolean
}
```

## `@go.tag`

used for: `StructField`

Add struct tag for golang struct field.

Example:

```
struct User {
  # @go.tag gorm:"uniqueIndex"
  email: string
}
```

## `@go.package`

## `@go.import`

## `@ts.modifier`

## `@deprecated`

## `@default`

## `@minimum` & `@maximum`

## `@minLength` & `@maxLength`

