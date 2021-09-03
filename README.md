# 基于golang 开发的轻量web框架
* 内置部分web常用包
  * session
  * sql

# 安装
```
go get -u github.com/cjie9759/goWeb
```

# demo运行
``` golang
package main

import (
	"fmt"

	"github.com/cjie9759/goWeb"
	"github.com/cjie9759/goWeb/controller"
)

func main() {
  fmt.Println(
    // 创建app
    goWeb.NewApp().
		// 注册服务
		Get(&controller.Index{}).
		// 加载中间件
		SetMiddle(goWeb.MWLog).
		// 运行
		Run(":17127"))
}
```