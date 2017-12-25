# 用来实时监控支付宝，进行交易

## 只需个人用户，不需要企业或商家认证

## 采用golang标准库编写，无第三方库依赖

## 采用web界面输入cookie进行登陆方式，简单

## 部署步骤
### 1.安装go，配置好GOPATH和GOBIN环境变量(如果没有配置好，自行google)
### 2.下载项目
```shell
$ go get github.com/changvvb/alinotify 
 ```
### 3.编译
```shell
$ go install github.com/changvvb/alinotify 
 ```
### 4.运行
```shell
$ alinotify
```

### 5.打开浏览其 https://my.alipay.com 登陆后，获得请求头中的cookies(许多cookie!)复制
### 6.打开浏览器 http://127.0.0.1:2048/setcookie ,将复制到的cookies粘贴提交
### 7.交易后到 http://127.0.0.1:2048/exam?tel=<电话>&email=< Email> 查看

