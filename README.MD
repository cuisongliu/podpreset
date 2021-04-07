
## 命令参数
webhook-port webhook端口 默认9443
metrics-addr 监控地址 默认:8080
health-addr  健康检查地址 默认:9090
enable-leader-election 是否开启选举，开启后只会有一个工作. 默认 false
log-level 日志级别(debug,info,warn,error,dpanic,panic,fatal) 默认info

## 环境变量
NAMESPACE_NAME 存储秘钥和svc的命名空间 默认值 kube-system
SVC_NAME 当前部署的svc名称 默认值 webhook-service
SECRET_NAME 存储证书的秘钥名称 默认值 webhook-secret
CSR_NAME 提供证书的资源名称 默认值 webhook-csr
MUTATING_NAME mutating webhook名称 默认值 mutating-cfg
VALIDATING_NAME  validating webhook名称 默认值 validating-cfg
