# api1 (api first)

api1 is:
1. An api definition language
2. An api doc generating tool (openapi for now)
3. An api code generating tool (golang server code for now)

## api1 definition language specification

Api definition contains 2 parts:
1. function definitions (which can be invoked)
2. type/struct definitions (as input/output types of these functions)

Most programming languages provides these 2 abilities.

api1(api first) 是一个代码和文档生成工具。这个工具会从 `*.api` 文件中读取api定义，并支持生成：

- 函数声明中用到的结构体类型定义（含输入参数类型或返回类型）

本文档将详细介绍 `*.api` 文件的语法规则。作为api1的使用者，请务必仔细阅读本文档。

## `*.api` 文件语法的设计原则

1. 尽量使用各种语言中常见的数据类型和语法特征。
2. 会在简单易用和强大的定义能力中寻求平衡点，仅以满足常见api定义需求为目标。

## 参考语言

- golang
- javascript/typescript
- kotlin
- graphql
- protobuf

## Built-in Scalar Types

|api|golang|javascript|typescript|openapi|python|
|--|--|--|--|--|--|
|int|int64|number|number|integer|int|
|float|float64|number|number|number|float|
|string|string|string|string|string|str|
|boolean|bool|boolean|boolean|boolean|bool|
|object|map|object|object|object|dict|
|any|interface{}|-|any|oneOf[]|typing.Any|

## Custom Scalar Types

例子

```
scalar datetime
```

## Real world scalar types

```
# @go.type string
# @openapi.type string
# @openapi.format password
scalar Password

# example: 2022-11-29T03:09:18.031Z
# @go.type time/Time
# @openapi.type string
# @openapi.format date-time
scalar Time

# example: 1669691224000
# @go.type int64
# @go.typeDef
# @openapi.type integer
scalar Timestamp

# example: 2006/01/02 15:04
# @go.type string
# @go.typeDef
# @openapi.type string
scalar TimeTillMinute

# @go.type *multipart.FileHeader
# @go.typePkg mime/multipart
# @openapi.type string
# @openapi.format binary
scalar MultipartFile
```

## 枚举类型

枚举类型在api定义中十分常见，枚举类型的定义语法如下：

例子：

```
enum ProtocolType {
  TCP
  UDP
}
```

枚举项使用自定义Value的例子：

```
enum Colour {
  RED = "red"
  GREEN = "green"
  BLUE = "blue"
}
```

枚举项使用数值型Value的例子：

```
enum Colour {
  RED = 1
  GREEN = 2
  BLUE = 3
}
```

## 结构体类型

结构体类型中包含明确的字段定义，每个字段的类型也都是明确的。
结构体类型的定义语法如下：

例子：

```
struct User {
  id: int
  name: string
  email: string
  admin: boolean
}
```

## 数组类型

基本类型、枚举类型、以及结构体类型，均可定义为数组类型。
定义数组类型的语法如下：

例子：

```
[string]
[ProtocolType]
[User]
[[int]]  # 数组的数据，即二维数组
```

## 类型可为空

当类型后面带有 `?` 时，表示当前类型可为空。

例子：

```
string?    # 可为null的字符串
[string]?  # 可为null的字符串数组
[string?]  # 数组中的元素可为null
```

更多例子：

```
struct User {
  name: string
  address: string?
  phoneNumbers: [string]?
}
```

## 接口

接口中定义了可交互的方法（函数），该方法接受0～N个入参，返回0～1个出参。
参数类型可以是基本类型、枚举类型、结构体类型或数组类型。

接口类型的定义语法如下：

```
interface UserController {

  register(req: RegisterReq)

  login(req: LoginReq)

  getUserInfo(req: UserInfoReq): UserInfoResp

  listUsers(pageSize: int, pageNo: int): [User]
}
```

## 注释

以 `#` 符号来定义注释

## 语义注释

语义注释是注释的扩展，它类似SQL语言中的Hit，不以语法形式存在。
语义注释定例子如下：

```
interface UserController {

  # @route post /user/login
  login(req: LoginReq)
}

# @route get /another
# @go.package some/package
# @go.someAttr:json|
#   {
#     "a": 1,
#     "b": 2,
#   }
interface AnotherController {

}
```

## Known Semantic Comments

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

## `@form` & `@accept`

```
# @go.type *multipart.FileHeader
# @go.typePkg mime/multipart
# @openapi.type string
# @openapi.format binary
scalar MultipartFile

# @form
struct UploadInfo {
  file: MultipartFile
  path: string
}

interface UploadApi {
  # @route post /actions/upload
  # @accept multipart/form-data
  doUpload(req: UploadInfo): object
}
```

## 与package相关的注释

## 与web api相关的注释

## 问题：
1、合法性检查
2、高亮
3、代码之间跳转


