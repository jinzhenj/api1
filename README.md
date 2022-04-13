

# go-swagger

## usage
* cd e2e && go test 后观察生成的 swagger.json 和其他的文件

## definitions
* api 
```
  service user {

  	@doc(
  		summary: 注册
  	)
  	@handler register 
  	post /api/user/register (pkg.types.RegisterReq)  # pkg.types.RegisterReq 代表 pkg/types 目录下某个文件定义的结构体 RegisterReq

  	@doc(
  		summary: 登录
  	)
  	@handler login
  	post /api/user/login (pkg.config.LoginReq)

  	@doc(
  		summary: 获取用户信息
  	)
  	@handler getUserInfo
  	get /api/user/:id (pkg.config.UserInfoReq) return (pkg.config.UserInfoReply)

  	@doc(
  		summary: 用户搜索
  	)
  	@handler searchUser
  	get /api/user/search (pkg.config.UserSearchReq{list=[]string, records=[]int}) return (pkg.config.UserInfoReply)

  }

```

* struct
```
type Field struct {
	Name     string `json:"name,omitempty" position:"query"`  // query 中的参数
	Tag      *Tag   `json:"tag,omitempty" position:"path"`    // 路径中的参数
	Kind     TypeD  `json:"kind,omitempty" position:"path"`   
	Comments string `json:"comments,omitempty" position:"body"` // json body 中的参数
}

-- 目前要求的所有的结构体都必须定义 pkg/{any folder}/xxx.go 文件，即必须定义在 pkg 子目录下文件。不能是子目录的子目录下的文件
```


## Features 
* 解析 golang 结构体
   * 解析无 Field 注释的结构体 [done]
   * 解析有 Field 注释的，并将注释作为 swagger descrition 信息 [done]
   * 解析自定义类型的 golang 结构体 [done]
   * 支持 type TypeA TypeB 写法
   * 支持结构体定义不在 pkg/ 目录下
   * 支持结构体定义在深层次的目录下
* 生成 swagger
   * 生成 swager json 文件[done]
   * 只渲染使用到的 go struct 结构体而非项目中 pkg/ 目录下的结构体
* 生成的 go 路由接口文件
  * 生成 go 路由接口文件 [doing]
  * 覆盖更新路由条目时，如果路由上的注册函数有装饰器则输出提示信息
  * 生成 go 路由接口文件改用 template 生成  
* Misc
   * 日志修改成可以在运行时修改日志等级，根据 flag 设定
   * 补充说明使用文档 

   

