package exec

import (
	"HCPlatform/code/internal"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

type ProfileService struct {
}

func (s *ProfileService) ProfileByNNMeter(ctx context.Context, request *ProfileRequest) (*ProfileResponse, error) {
	resp := new(ProfileResponse)

	path, err := convertBytesToFile(request.ModelFile.Filename, request.ModelFile.Data)
	if err != nil {
		resp.Msg = err.Error()
	}
	args := request.GetNnmeterArgs()
	res := profileByNNMeter(path, args.Predictor, args.Version, args.Framework)
	resp.Msg = fmt.Sprintf("nn-meter predictor:%s predictor-version:%s framework:%s\n%s", args.Predictor, args.Version, args.Framework, res)
	return resp, nil
}

func profileByNNMeter(path string, predictor string, version string, framework string) string {
	//打开shell查看执行状态...
	shell, _ := internal.NewPowerShell()
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

func (s *ProfileService) mustEmbedUnimplementedProfileServer() {
	//TODO implement me
	panic("implement me")
}

func NewProfileRequest(path string, nnmeterPredictorName string, nnmeterPredictorFramework string) *ProfileRequest {
	file := new(File)
	data, size := convertFileToBytes(path)
	file.Data = data
	file.Size = size
	if size == 0 {
		return nil
	}
	//fmt.Println(path)
	file.Filename = splitFilenameFromFilePath(path)
	//fmt.Println(file.Filename)
	if file.Filename == "" {
		return nil
	}
	rq := new(ProfileRequest)
	rq.ModelFile = file
	rq.Args = &ProfileRequest_NnmeterArgs{NnmeterArgs: &NNMeterArgs{
		Predictor: nnmeterPredictorName,
		Version:   "1.0",
		Framework: nnmeterPredictorFramework,
	}}
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
