package register

import (
	"container/list"
	context "context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

var (
	RegisteredDevicesMap map[string]uint64
	unusedDevicesMap     map[uint64]*DeviceMessage
	usedDevicesMap       map[uint64]*DeviceMessage
	IpPoolMap            map[uint64]string
)

func init() {
	RegisteredDevicesMap = make(map[string]uint64)
	unusedDevicesMap = make(map[uint64]*DeviceMessage)
	usedDevicesMap = make(map[uint64]*DeviceMessage)
	IpPoolMap = make(map[uint64]string)
}

type RegisterService struct {
}

func (s *RegisterService) RegisterDevice(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	//TODO implement me
	resp := new(RegisterResponse)
	devices := request.GetDevices()
	if devices == nil {
		resp.Msg = "No Device!"
		return resp, nil
	}
	result := ""
	for _, device := range devices {
		//fmt.Println("%d,%s", i, device.DeviceName)
		deviceId := uint64(time.Now().Unix())
		device.AllocedId = deviceId
		RegisteredDevicesMap[device.DeviceName] = deviceId
		IpPoolMap[deviceId] = fmt.Sprintf("%s:%d", device.DeviceAddress, device.TerminalPort)
		unusedDevicesMap[deviceId] = device
		result = fmt.Sprintf("%s\n%s has registerd as %d", result, device.DeviceName, deviceId)
	}
	resp.Msg = result
	log.Info(result)
	return resp, nil
}

func (s *RegisterService) GetAllRegisteredDevice(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	//TODO implement me
	resp := new(RegisterResponse)
	result := ""
	for deviceName, deviceId := range RegisteredDevicesMap {
		result = fmt.Sprintf("%s\n%s,%d", result, deviceName, deviceId)
	}
	resp.Msg = result
	return resp, nil
}

func (s *RegisterService) AllocDevice(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	//TODO implement me
	resp := new(RegisterResponse)
	result := ""
	devices := request.Devices
	successAllocatedDevices := list.New()
	for _, device := range devices {
		deviceName := device.GetDeviceName()
		// check the device has or not registered on server
		deviceId, ok := RegisteredDevicesMap[deviceName]
		if !ok {
			result = fmt.Sprintf("%s\nDevice:%s is not in registered device list.", result, deviceName)
			continue
		}
		// check the device has or not allocated to others
		_, ok = usedDevicesMap[deviceId]
		if ok {
			result = fmt.Sprintf("%s\nDevice:%s is not free now. You can't connect it.", result, deviceName)
			continue
		}
		// fetch the device message from unusedMap
		registerDevice := unusedDevicesMap[deviceId]
		usedDevicesMap[deviceId] = registerDevice
		delete(unusedDevicesMap, deviceId)
		result = fmt.Sprintf("%s\nDevice:%s is allocated to you.", result, deviceName)
		successAllocatedDevices.PushBack(registerDevice)
	}
	resp.Msg = result
	resp.Devices = make([]*DeviceMessage, successAllocatedDevices.Len())
	index := 0
	for i := successAllocatedDevices.Front(); i != nil; i = i.Next() {
		resp.Devices[index] = (i.Value).(*DeviceMessage)

		index++
	}
	return resp, nil
}

func (s *RegisterService) FreeDevice(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	//TODO implement me
	resp := new(RegisterResponse)
	result := ""
	devices := request.Devices
	for _, device := range devices {
		deviceName := device.GetDeviceName()
		deviceId, ok := RegisteredDevicesMap[deviceName]
		if !ok {
			result = fmt.Sprintf("%s\nDevice:%s is not in registered device list.", result, deviceName)
			continue
		}
		_, ok = unusedDevicesMap[deviceId]
		if ok {
			result = fmt.Sprintf("%s\nDevice:%s is free now. You can't free it again.", result, deviceName)
			continue
		}
		unusedDevicesMap[deviceId] = usedDevicesMap[deviceId]
		delete(usedDevicesMap, deviceId)
		result = fmt.Sprintf("%s\nDevice:%s is free, everyone can connect it now.", result, deviceName)
	}
	resp.Msg = result
	return resp, nil
}

func (s *RegisterService) mustEmbedUnimplementedReisgterServer() {
	//TODO implement me
	panic("implement me")
}

// This function is for client to fast call
func FastListAllDevices(serverIP string, serverPort int) (string, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := NewReisgterClient(conn)
	req := &RegisterRequest{}
	res, err := c.GetAllRegisteredDevice(context.Background(), req)
	if err != nil {
		return "", err
	}
	err = conn.Close()
	if err != nil {
		log.Error(err)
		return "", err
	}
	return res.Msg, nil
}

// This function is for client to fast call
func FastRegisterDevices(serverIP string, serverPort int, devices []*DeviceMessage) (string, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer conn.Close()
	c := NewReisgterClient(conn)
	req := &RegisterRequest{Devices: devices}
	res, err := c.RegisterDevice(context.Background(), req)
	if err != nil {
		log.Error(err)
		return res.Msg, err
	}
	err = conn.Close()
	if err != nil {
		log.Error(err)
		return res.Msg, err
	}
	return res.Msg, nil
}

// This function is for client to fast call
func FastAllocDevices(serverIP string, serverPort int, devicesName []string) ([]*DeviceMessage, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return []*DeviceMessage{}, err
	}
	defer conn.Close()
	c := NewReisgterClient(conn)
	// We use Slice with variable length
	devices := make([]*DeviceMessage, len(devicesName))
	for i, deviceName := range devicesName {
		devices[i] = &DeviceMessage{DeviceName: deviceName}
	}
	req := &RegisterRequest{Devices: devices}
	res, err := c.AllocDevice(context.Background(), req)
	if err != nil {
		log.Error(err)
		return []*DeviceMessage{}, err
	}
	err = conn.Close()
	if err != nil {
		log.Error(err)
		return []*DeviceMessage{}, err
	}
	return res.Devices, nil
}

// This function is for client to fast call
func FastFreeDevices(serverIP string, serverPort int, devicesName []string) (string, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverIP, serverPort), grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer conn.Close()
	c := NewReisgterClient(conn)
	// We use Slice with variable length
	devices := make([]*DeviceMessage, len(devicesName))
	for i, deviceName := range devicesName {
		devices[i] = &DeviceMessage{DeviceName: deviceName}
	}
	req := &RegisterRequest{Devices: devices}
	res, err := c.FreeDevice(context.Background(), req)
	if err != nil {
		log.Error(err)
		return res.Msg, err
	}
	err = conn.Close()
	if err != nil {
		log.Error(err)
		return res.Msg, err
	}
	return res.Msg, nil
}
