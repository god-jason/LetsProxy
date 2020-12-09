# LetsProxy

使用Golang开发的HTTPS反向代理，功能特点：
1. 内嵌Let's Encrypt自动证书，
2. 支持多组代理
3. 支持多个域名
4. 支持负载均衡
5. 支持Linux和Windows系统服务，系统重启也不怕

## 编译

```shell script
go build
```

国内用户可能需要设置golang编译环境，开启代理，关闭检验
```shell script
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=*.gitlab.com,*.gitee.com,git.zgwit.com
go env -w GOSUMDB=off
```

## 运行

```shell script
LetsProxy -h
Usage of LetsProxy:
  -c string
        配置文件 (default "LetsProxy.yaml")
  -h    帮助
  -i    安装服务
  -u    卸载服务

```

运行环境：
* Windows server 2008 及以上版本
* Linux Kernel 2.6 及以上发行版

注意：部分Linux发行版安装服务异常

## 配置文件

```yaml
//证书目录
cache: certs

//letsencrypt注册邮箱（未测试）
email: ""

//域名和目标服务器均支持多个（以逗号间隔）
proxies: 
  git.zgwit.com: http://127.0.0.1:3000
  a.com,b.com: http://192.168.0.12:80,http://192.168.0.13:80
```

## 其他
1. 项目参考 [audibleblink/letsproxy](https://github.com/audibleblink/letsproxy)
2. 性能未测试，基本满足日常需要（专业用户请移步nginx）
3. 暂无界面开发计划
 

## 真格智能实验室
[<真格智能实验室>](https://zgwit.com)