# HeterogeneousComputing

## CMD
### 安装命令
```bash
#设置国内镜像,Linux or Mac转到 https://goproxy.cn/
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
# 在项目根目录执行
go mod download
```

### 测试
```bash
go run code/main.go profile --modelPath="1.onnx" --deviceName="armcpu" --framework="tflite2.1"
```

## TODO
1. 使用GO创建一个基本的控制程序，能够识别到异构设备，USB Device、PC、mobile phone。(目前仅支持本地连接)
2. 创建远程连接。可以通过ssh连接到这些设备。
3. 模型性能分析(目前考虑支持TFLite/PaddleLite/onnxruntime)

### Models
###
1. TFLite
2. PaddleLite
3. onnxruntime