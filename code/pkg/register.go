package pkg

import (
	"HCPlatform/code/protos/register"
	context "context"
)

var (
	registeredDevicesMap map[string]string
	allocedDevicesMap    map[string]string
)

type RegisterService struct {
}

func (s *RegisterService) ResgisterDevice(ctx context.Context, request *protos.RegisterRequest) (*protos.RegisterResponse, error) {
	//TODO implement me
	resp := new(protos.RegisterResponse)
	resp.Msg = "OK"
	return resp, nil
}

func (s *RegisterService) mustEmbedUnimplementedReisgterServer() {
	//TODO implement me
	panic("implement me")
}
