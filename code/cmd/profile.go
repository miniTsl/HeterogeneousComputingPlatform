package cmd

import (
	"HCPlatform/code/protos/exec"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	modelPath           string
	profilerName        string
	dstDeviceName       string
	warmsup_rounds      int     //模型热身次数
	num_rounds          int     //模型正式运行次数
	delayBetweenRound   float32 //每轮之间的运行延迟
	enable_op_profiling bool    //模型是否启用OP分析

	//
	profileByMobileCPU bool
	num_threads        int
	profileByMobileGPU bool

	profileByNNAPI   bool
	profileByCoreML  bool
	profileByXNNPack bool

	//nn-meter
	nnmeterPredictorName string
	nnmeterPredictorType string
	profileCmd           = &cobra.Command{
		Use:   "profile",
		Short: "profile a model on specific device",
		Long:  "profile a model on specific device",
		Run: func(cmd *cobra.Command, args []string) {
			serverIP, serverPort := "192.168.1.106", 9520

			if profilerName == "nn-meter" {
				res := exec.FastNNMeterProfile(serverIP, serverPort, dstDeviceName, modelPath, nnmeterPredictorName, nnmeterPredictorType)
				fmt.Println(res)
			} else if profilerName == "TFLite" {
				res := exec.FastTFLiteProfile(serverIP, serverPort, dstDeviceName, modelPath, warmsup_rounds, num_rounds, delayBetweenRound, enable_op_profiling, num_threads, profileByMobileCPU, profileByMobileGPU)
				fmt.Println(res)
			}

			//log.Info(res.Msg)

		},
	}
)

// go run code/main.go profile --profilerName=nn-meter --nnmeterPredictorName=cortexA76cpu_tflite21 --nnmeterPredictorType=onnx --modelPath=D:\code\HeterogeneousComputingPlatform\model\resnet18-12.onnx
// go run code/main.go profile --profilerName=nn-meter --nnmeterPredictorName=adreno640gpu_tflite21 --nnmeterPredictorType=onnx --modelPath=D:\code\HeterogeneousComputingPlatform\model\resnet18-12.onnx
func init() {
	//connectCmd.PersistentFlags().BoolVar(&NNMeter, "nn-meter", true, "profile by nn-meter")
	profileCmd.PersistentFlags().StringVar(&modelPath, "modelPath", "", "The path to the model file.")
	profileCmd.MarkPersistentFlagRequired("modelPath")
	profileCmd.PersistentFlags().StringVar(&profilerName, "profilerName", "nn-meter", "optional: nn-meter,PaddleLite,TFLite,onnxruntime")
	profileCmd.MarkPersistentFlagRequired("profilerName")
	profileCmd.PersistentFlags().StringVar(&dstDeviceName, "deviceName", "any", "the device that you selected")
	profileCmd.PersistentFlags().StringVar(&nnmeterPredictorName, "nnmeterPredictorName", "cortexA76cpu_tflite21", "optional: cortexA76cpu_tflite21,adreno640gpu_tflite21,adreno630gpu_tflite21,myriadvpu_openvino2019r2")
	profileCmd.PersistentFlags().StringVar(&nnmeterPredictorType, "nnmeterPredictorType", "onnx", "optional: tensorflow,onnx,torch")

	profileCmd.PersistentFlags().IntVar(&warmsup_rounds, "warmsup_runs", 1, "The number of warmup runs to do before starting the benchmark.")
	profileCmd.PersistentFlags().IntVar(&num_rounds, "num_runs", 10, "The number of runs. Increase this to reduce variance.")
	profileCmd.PersistentFlags().Float32Var(&delayBetweenRound, "run_delay", -1.0, "The delay in seconds between subsequent benchmark runs. Non-positive values mean use no delay.")
	profileCmd.PersistentFlags().BoolVar(&enable_op_profiling, "enable_op_profiling", false, "Whether to enable per-operator profiling measurement.")

	profileCmd.PersistentFlags().BoolVar(&profileByMobileCPU, "use_cpu", true, "Use Mobile CPU profile the model.")
	profileCmd.PersistentFlags().IntVar(&num_threads, "num_threads", -1, "The number of threads to use for running TFLite interpreter. By default, this is set to the platform default value -1.")
	profileCmd.PersistentFlags().BoolVar(&profileByMobileGPU, "use_gpu", false, "Use Mobile GPU profile the model.")
	//profileCmd.PersistentFlags().BoolVar(&profileByNNAPI, "use_nnapi", false, "Note some Android P devices will fail to use NNAPI for models in /data/local/tmp/ and this benchmark tool will not correctly use NNAPI.")
	//profileCmd.PersistentFlags().BoolVar(&profileByCoreML,"use_coreml",false,"")

}
