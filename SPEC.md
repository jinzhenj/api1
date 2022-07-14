# api1 specification

api1(api first) 是一个代码和文档生成工具。这个工具会从 `*.api` 文件中读取api定义，并支持生成：
- golang代码
- js/ts代码
- swagger文档

我们说的api，主要指两部分内容：
- 接口函数声明
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

## 内建基本类型

|api|golang|javascript|typescript|openapi|
|--|--|--|--|--|
|int|int|number|number|integer|
|float|float64|number|number|number|
|string|string|string|string|string|
|boolean|bool|boolean|boolean|boolean|
|object|map|object|object|object|
|any|interface{}|-|any|-|

## 自定义基本类型

例子

```
scalar datetime
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

## 与package相关的注释

## 与web api相关的注释

## 问题：
1、合法性检查
2、高亮
3、代码之间跳转


