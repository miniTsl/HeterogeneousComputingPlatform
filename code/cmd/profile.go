package cmd

import (
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
			serverIP, serverPort := "127.0.0.1", 9520

			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			c := exec.NewProfileClient(conn)
			req := exec.NewProfileRequest(modelPath, nnmeterPredictorName, nnmeterPredictorType)
			res, err := c.ProfileByNNMeter(context.Background(), req)
			if err != nil {
				log.Fatal(err)
			}

			//log.Info(res.Msg)
			fmt.Println(res.Msg)
		},
	}
)

func init() {
	profileCmd.PersistentFlags().StringVar(&modelPath, "modelPath", "", "file path")
	profileCmd.MarkPersistentFlagRequired("modelPath")
	profileCmd.PersistentFlags().StringVar(&profilerName, "profilerName", "nn-meter", "optional: nn-meter,paddle-lite,tensorflow-lite,onnxruntime")
	profileCmd.MarkPersistentFlagRequired("profilerName")
	profileCmd.PersistentFlags().StringVar(&nnmeterPredictorName, "nn-meter predictor", "cortexA76cpu_tflite21", "optional: cortexA76cpu_tflite21,adreno640gpu_tflite21,adreno630gpu_tflite21,myriadvpu_openvino2019r2")
	profileCmd.PersistentFlags().StringVar(&nnmeterPredictorType, "nn-meter framework", "onnx", "optional: tensorflow,onnxruntime,torch")

}
