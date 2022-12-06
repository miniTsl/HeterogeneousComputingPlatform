package protos

import (
	context "context"
)

type RegisterService struct {
}

func (s *RegisterService) ResgisterDevice(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	//TODO implement me
	resp := new(RegisterResponse)
	resp.Msg = "OK"
	return resp, nil
}

func (s *RegisterService) mustEmbedUnimplementedReisgterServer() {
	//TODO implement me
	panic("implement me")
}
