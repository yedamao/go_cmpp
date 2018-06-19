# go_cmpp
[![Build Status](https://travis-ci.org/yedamao/go_cmpp.svg?branch=master)](https://travis-ci.org/yedamao/go_cmpp)
[![Go Report Card](https://goreportcard.com/badge/github.com/yedamao/go_cmpp)](https://goreportcard.com/report/github.com/yedamao/go_cmpp)
[![codecov](https://codecov.io/gh/yedamao/go_cmpp/branch/master/graph/badge.svg)](https://codecov.io/gh/yedamao/go_cmpp)


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
  -addr string
        监听地址 (default ":8801")
```

### transmitter
提交单条短信至短信网关

```
Usage of ./bin/transmitter:
  -area-code string
        长途区号 (default "010")
  -corp-id string
        5位企业代码 (default "00000")
  -dest-number string
        接收手机号码, 86..., 多个使用，分割
  -host string
        SMSC host (default "localhost")
  -msg string
        短信内容
  -name string
        Login Name
  -passwd string
        Login Password
  -port int
        SMSC port (default 8801)
  -service-type string
        业务代码，由SP定义
  -sp-number string
        SP的接入号码
```

### receiver
接收运营商回执状态/上行短信

```
Usage of ./bin/receiver:
  -addr string
        上行监听地址 (default ":8001")
  -count int
        worker 数量 (default 5)
```


