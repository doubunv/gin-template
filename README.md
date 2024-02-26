## 全局说明：
1. 所有函数错误需要显示返回，即需要定义一个返回error
2. 在自定义协程里面，必须手动加上 recover 
3. 日志集成调用了链路号，在配置文件控制打印级别，打印统一采用 `tool.Debug() tool.Info() tool.Error()` 输出打印

## API接口层：
1. 该层的设计是基于gin框架进行的二次封装，具有gin框架的所有能力。
2. router中定义路由是二级路由[eg: /controller/function]，但是暴露给client端为三级路由 [eg: /api-server/controller/function]。
3. 接口请求参数和返回值，采用proto文件定义，采用json参数验证逻辑。
4. 接口文档采用proto-swagger插件，自动生成swagger接口文档。 
5. 通过make proto 命令自动生成，shell proto.pb.go 文件, 会覆盖旧文件。
6. 在该层可以通过`headInfo.Getxxx()`获取当前登陆的一些信息。