package cmd

import (
	"HCPlatform/code/pkg"
	"HCPlatform/code/protos/exec"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	modelPath    string
	profilerName string
	deviceName   string

	nnmeterPredictorName string
	nnmeterPredictorType string
	profileCmd           = &cobra.Command{
		Use:   "profile",
		Short: "profile a model on specific device",
		Long:  "profile a model on specific device",
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: 根据设备查询连接方式
			//向服务器发送一个模型文件
			//fmt.Print(modelPath, profilerName)
			serverIP, serverPort := "192.168.1.106", 9520

			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			c := exec.NewProfileClient(conn)
			fmt.Println("connect to ", serverIP)
			req := pkg.NewProfileRequest(modelPath, nnmeterPredictorName, nnmeterPredictorType)
			fmt.Println("upload model", modelPath, " to ", serverIP)
			res, err := c.ProfileWithArgs(context.Background(), req)
			if err != nil {
				log.Fatal(err)
			}

			//log.Info(res.Msg)
			fmt.Println(res.Msg)
		},
	}
)

// go run code/main.go profile --profilerName=nn-meter --nnmeterPredictorName=cortexA76cpu_tflite21 --nnmeterPredictorType=onnx --modelPath=D:\code\HeterogeneousComputingPlatform\model\resnet18-12.onnx
// go run code/main.go profile --profilerName=nn-meter --nnmeterPredictorName=adreno640gpu_tflite21 --nnmeterPredictorType=onnx --modelPath=D:\code\HeterogeneousComputingPlatform\model\resnet18-12.onnx
func init() {
	//connectCmd.PersistentFlags().BoolVar(&NNMeter, "nn-meter", true, "profile by nn-meter")
	profileCmd.PersistentFlags().StringVar(&modelPath, "modelPath", "", "devcie configuration")
	profileCmd.MarkPersistentFlagRequired("modelPath")
	profileCmd.PersistentFlags().StringVar(&profilerName, "profilerName", "nn-meter", "optional: nn-meter,paddle-lite,tensorflow-lite,onnxruntime")
	profileCmd.MarkPersistentFlagRequired("profilerName")
	profileCmd.PersistentFlags().StringVar(&nnmeterPredictorName, "nnmeterPredictorName", "cortexA76cpu_tflite21", "optional: cortexA76cpu_tflite21,adreno640gpu_tflite21,adreno630gpu_tflite21,myriadvpu_openvino2019r2")
	profileCmd.PersistentFlags().StringVar(&nnmeterPredictorType, "nnmeterPredictorType", "onnx", "optional: tensorflow,onnx,torch")

}
