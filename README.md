# 新版testvm
## 功能
1. 提供接口给客户端返回可用的服务器节点
2. 提供接口支持写入某个服务器进程负载信息
3. 提供接口获取当前IDC的webrtc的负载信息

## 第三方依赖包
- github.com/beego/bee
- github.com/astaxie/beego
- github.com/Sirupsen/logrus
- github.com/garyburd/redigo/redis

### 下载依赖包
> go get *package name*

## 运行
> cd 到testvm目录
> $ bee run