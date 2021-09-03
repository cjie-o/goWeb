# 基于golang 开发的轻量web框架
* 内置部分web常用包
  * session
  * sql

# 安装

# demo运行
``` golang
fun main(){
  fmt.Println(NewApp().Get(&controller.Index{}).SetMiddle(MWLog).Run(":17127"))
}
```