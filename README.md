# go_cmpp

go_cmpp是为SP设计实现的CMPP2.0协议开发工具包。包括cmpp协议包和命令行工具。

## 安装
```
go get github.com/yedamao/go_cmpp/...
cd $GOPATH/src/github.com/yedamao/go_cmpp && make
```

## Cmpp协议包

### 支持操作
- [x] CMPP_CONNECT
- [x] CMPP_TERMINATE
- [x] CMPP_TERMINATE_RESP
- [x] CMPP_SUBMIT
- [x] CMPP_QUERY
- [x] CMPP_DELIVER_RESP
- [x] CMPP_CANCEL
- [x] CMPP_ACTIVE_TEST
- [x] CMPP_ACTIVE_TEST_RESP

### Example
参照cmd/transmitter/main.go, cmd/receiver/main.go

## 命令行工具

### mockserver
ISMG短信网关模拟器

```
Usage of ./bin/mockserver:
  -activeTest
        是否activeTest
  -addr string
        addr(本地监听地址) (default ":7890")
  -mo
        是否模拟上行短信
  -rpt
        是否模拟上行状态报告
```

### transmitter
提交单条短信至短信网关

```
Usage of ./bin/transmitter:
  -addr string
        smgw addr(运营商地址) (default ":7890")
  -dest-number string
        接收手机号码, 86..., 多个使用，分割
  -msg string
        短信内容
  -secret string
        登陆密码
  -serviceId string
        业务类型，是数字、字母和符号的组合
  -sourceAddr string
        源地址，即SP的企业代码
  -srcId string
        SP的接入号码
```

### receiver
接收运营商回执状态/上行短信

```
Usage of ./bin/receiver:
  -addr string
        smgw addr(运营商地址) (default ":7890")
  -secret string
        登陆密码
  -sourceAddr string
        源地址，即SP的企业代码
```


