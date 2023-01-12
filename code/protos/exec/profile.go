package exec

import (
	"HCPlatform/code/pkg"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"strings"
)

type ProfileService struct {
}

func (s *ProfileService) GetModelStaticAttr(ctx context.Context, request *ProfileRequest) (*ProfileResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ProfileService) GetDeviceStaticAttr(ctx context.Context, request *ProfileRequest) (*ProfileResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ProfileService) GetProfileAbility(ctx context.Context, request *ProfileRequest) (*ProfileResponse, error) {
	//TODO implement me
	/*
		returns the platform profiling ability
	*/
	panic("implement me")
}

func (s *ProfileService) ProfileWithArgs(ctx context.Context, request *ProfileRequest) (*ProfileResponse, error) {
	resp := new(ProfileResponse)

	path, err := convertBytesToFile(request.ModelFile.Filename, request.ModelFile.Data)
	if err != nil {
		resp.Msg = err.Error()
	}

	//TODO 申请该设备
	//register.FastAllocDevices()
	//deviceId:=register.RegisteredDevicesMap[request.DeviceName]
	//serverAddr := register.IpPoolMap[deviceId]
	//shellId,err:=term.FastNewTerminalAddr(serverAddr,deviceId,0)

	//打开device的远程shell
	profileType := request.Type
	if profileType == ProfileRequest_nnMeter {
		args := request.GetNnmeterArgs()
		res := profileByNNMeter(path, args.Predictor, args.Version, args.Framework)
		resp.Msg = fmt.Sprintf("nn-meter predictor:%s predictor-version:%s framework:%s\n%s", args.Predictor, args.Version, args.Framework, res)
	} else if profileType == ProfileRequest_tflite {
		args := request.GetTfliteArgs()
		deviceType := args.DeviceType
		params := args.Params
		delegateParams := args.DelegateParams
		res := profileByTFLite(path, request.DeviceName, deviceType, params, delegateParams)
		resp.Msg = fmt.Sprintf("TFLite\n%s", res)
	}

	return resp, nil
}

func (s *ProfileService) mustEmbedUnimplementedProfileServer() {
	//TODO implement me
	panic("implement me")
}

func profileByNNMeter(path string, predictor string, version string, framework string) string {
	//打开shell查看执行状态...
	shell, _ := pkg.NewPowerShell()
	//打开nn-meter执行环境
	sout, serr, err := shell.Execute("conda activate nn-meter")
	sout, serr, err = shell.Execute(fmt.Sprintf("nn-meter predict --predictor %s --predictor-version %s --%s %s", predictor, version, framework, path))
	if err != nil {
		fmt.Println(sout, "\n", serr)
	} else {
		fmt.Println(sout)
	}
	return sout
}

func profileByTFLite(path string, deviceSerial string, deviceType TFLiteArgs_DeviceType, params *TFLiteParameters, delegateParams *TFLiteDelegateParameters) string {
	// TODO 根据deviceIP来找shell服务.
	shell, _ := pkg.NewPowerShell()

	//打开nn-meter执行环境
	sout, serr, err := shell.Execute(fmt.Sprintf("adb -s %s  shell \"mkdir -p /data/local/tmp/tflite_models\"", deviceSerial))
	if err != nil {
		fmt.Println(sout, "\n", serr)
	} else {
		fmt.Println(sout)
	}
	sout, serr, err = shell.Execute(fmt.Sprintf("adb -s %s push  D:\\code\\HeterogeneousComputingPlatform\\tmp\\android_aarch64_benchmark_model_performance_options /data/local/tmp", deviceSerial))
	if err != nil {
		fmt.Println(sout, "\n", serr)
	} else {
		fmt.Println(sout)
	}
	sout, serr, err = shell.Execute(fmt.Sprintf("adb -s %s  shell \"chmod +x /data/local/tmp/android_aarch64_benchmark_model_performance_options\"", deviceSerial))
	if err != nil {
		fmt.Println(sout, "\n", serr)
	} else {
		fmt.Println(sout)
	}
	sout, serr, err = shell.Execute(fmt.Sprintf("adb -s %s push %s /data/local/tmp/tflite_models", deviceSerial, path))
	if err != nil {
		fmt.Println(sout, "\n", serr)
	} else {
		fmt.Println(sout)
	}
	modelName := splitFilenameFromFilePath(path)
	if deviceType == TFLiteArgs_cpu {
		num_threads := params.NumThreads
		warmup_runs := params.WarmupRuns
		enable_op_profiling := params.EnableOpProfiling
		num_runs := params.NumRuns
		command := fmt.Sprintf("adb  -s %s shell \"/data/local/tmp/android_aarch64_benchmark_model_performance_options "+
			"--num_threads=%d --graph=/data/local/tmp/tflite_models/%s --warmup_runs=%d --enable_op_profiling=%t --num_runs=%d\"", deviceSerial, num_threads, modelName, warmup_runs, enable_op_profiling, num_runs)
		sout, serr, err = shell.Execute(command)
		if err != nil {
			fmt.Println(sout, "\n", serr)
		} else {
			fmt.Println(sout)
		}
	} else if deviceType == TFLiteArgs_gpu {
		warmup_runs := params.WarmupRuns
		enable_op_profiling := params.EnableOpProfiling
		num_runs := params.NumRuns
		command := fmt.Sprintf("adb -s %s  shell \"/data/local/tmp/android_aarch64_benchmark_model_performance_options "+
			"--use_gpu=%t --graph=/data/local/tmp/tflite_models/%s --warmup_runs=%d --enable_op_profiling=%t --num_runs=%d\"", deviceSerial, true, modelName, warmup_runs, enable_op_profiling, num_runs)
		fmt.Println(command)
		sout, serr, err = shell.Execute(command)
		if err != nil {
			fmt.Println(sout, "\n", serr)
		} else {
			fmt.Println(sout)
		}

	}

	return sout

}

func profileByPaddleLite(path string, version string, deviceName string) string {
	// First,call RPC Register Service to alloc deviceId.

	// Second,call RPC Terminal Service to open a shell by allocated deviceId.

	// Third, run specified commands to profile the submitted model.

	//

	return ""
}

//func NewProfilePaddleLiteRequest(path string, deviceName string, paddleVersion string) *ProfileRequest {
//	file := new(File)
//	data, size := convertFileToBytes(path)
//	file.Data = data
//	file.Size = size
//	if size == 0 {
//		return nil
//	}
//	file.Filename = splitFilenameFromFilePath(path)
//	if file.Filename == "" {
//		return nil
//	}
//	rq := new(ProfileRequest)
//	rq.ModelFile = file
//	rq.Args = &ProfileRequest_PaddleLiteArgs{PaddleLiteArgs: &PaddleLiteArgs{
//		Version: "latest",
//	}}
//	return rq
//}

func NewProfileNNMeterRequest(path string, deviceName string, nnmeterPredictorName string, nnmeterPredictorFramework string) *ProfileRequest {
	file := new(File)
	data, size := convertFileToBytes(path)
	file.Data = data
	file.Size = size
	if size == 0 {
		return nil
	}
	file.Filename = splitFilenameFromFilePath(path)
	if file.Filename == "" {
		return nil
	}
	rq := new(ProfileRequest)
	rq.ModelFile = file
	rq.DeviceName = deviceName
	rq.Type = ProfileRequest_nnMeter
	rq.Args = &ProfileRequest_NnmeterArgs{NnmeterArgs: &NNMeterArgs{
		Predictor: nnmeterPredictorName,
		Version:   "1.0",
		Framework: nnmeterPredictorFramework,
	}}
	return rq
}

func NewProfileTFLiteCPURequest(path string, deviceName string, warmsup_rounds int, num_rounds int, delayBetweenRound float32, enable_op_profiling bool,
	num_threads int) *ProfileRequest {
	file := new(File)
	data, size := convertFileToBytes(path)
	file.Data = data
	file.Size = size
	if size == 0 {
		return nil
	}
	file.Filename = splitFilenameFromFilePath(path)
	if file.Filename == "" {
		return nil
	}
	rq := new(ProfileRequest)
	rq.ModelFile = file
	rq.DeviceName = deviceName
	rq.Type = ProfileRequest_tflite
	rq.Args = &ProfileRequest_TfliteArgs{
		TfliteArgs: &TFLiteArgs{
			DeviceType: TFLiteArgs_cpu,
			Params: &TFLiteParameters{
				WarmupRuns:        int32(warmsup_rounds),
				NumRuns:           int32(num_rounds),
				RunDelay:          delayBetweenRound,
				EnableOpProfiling: enable_op_profiling,
				NumThreads:        int32(num_threads),
			},
			DelegateParams: &TFLiteDelegateParameters{
				UseGpu:     false,
				UseNnapi:   false,
				UseCoreml:  false,
				UseHexagon: false,
				UseXnnpack: false,
			},
		},
	}
	return rq
}

func NewProfileTFLiteGPURequest(path string, deviceName string, warmsup_rounds int, num_rounds int, delayBetweenRound float32, enable_op_profiling bool) *ProfileRequest {
	file := new(File)
	data, size := convertFileToBytes(path)
	file.Data = data
	file.Size = size
	if size == 0 {
		return nil
	}
	file.Filename = splitFilenameFromFilePath(path)
	if file.Filename == "" {
		return nil
	}
	rq := new(ProfileRequest)
	rq.ModelFile = file
	rq.DeviceName = deviceName
	rq.Type = ProfileRequest_tflite
	rq.Args = &ProfileRequest_TfliteArgs{
		TfliteArgs: &TFLiteArgs{
			DeviceType: TFLiteArgs_gpu,
			Params: &TFLiteParameters{
				WarmupRuns:        int32(warmsup_rounds),
				NumRuns:           int32(num_rounds),
				RunDelay:          delayBetweenRound,
				EnableOpProfiling: enable_op_profiling,
			},
			DelegateParams: &TFLiteDelegateParameters{
				UseGpu:     true,
				UseNnapi:   false,
				UseCoreml:  false,
				UseHexagon: false,
				UseXnnpack: false,
			},
		},
	}
	return rq
}

func convertFileToBytes(path string) ([]byte, uint32) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("read file fail", err)
		return nil, 0
	}
	defer f.Close()

	fd, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return nil, 0
	}
	return fd, uint32(len(fd))
}

func convertBytesToFile(filename string, data []byte) (string, error) {
	// we will save the stream to TmpDir.
	tmp, err := os.CreateTemp("", filename)
	if err != nil {
		return "", err
	}
	tmp.Write(data)
	return tmp.Name(), err
}

func splitFilenameFromFilePath(path string) string {
	var separator string
	var name = ""
	if strings.ContainsAny(path, "\\") {
		separator = "\\"

	} else if strings.ContainsAny(path, "/") {
		separator = "/"
	}
	arr := strings.Split(path, separator)
	name = arr[len(arr)-1]
	return name
}

func FastNNMeterProfile(serverIP string, serverPort int, deviceName string, path string, nnmeterPredictorName string, nnmeterPredictorType string) string {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := NewProfileClient(conn)
	req := NewProfileNNMeterRequest(path, deviceName, nnmeterPredictorName, nnmeterPredictorType)
	res, err := c.ProfileWithArgs(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	return res.Msg
}

func FastTFLiteProfile(serverIP string, serverPort int, deviceName string, path string,
	warmsup_rounds int, num_rounds int, delayBetweenRound float32, enable_op_profiling bool,
	num_threads int, profileByMobileCPU bool, profileByMobileGPU bool) string {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := NewProfileClient(conn)
	var req *ProfileRequest
	if profileByMobileCPU {
		req = NewProfileTFLiteCPURequest(path, deviceName, warmsup_rounds, num_rounds, delayBetweenRound, enable_op_profiling, num_threads)
	} else if profileByMobileGPU {
		req = NewProfileTFLiteGPURequest(path, deviceName, warmsup_rounds, num_rounds, delayBetweenRound, enable_op_profiling)
	}
	res, err := c.ProfileWithArgs(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	return res.Msg
}
