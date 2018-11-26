package sakuracloud

import (
	"context"
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/pkg/errors"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/stretchr/testify/assert"
)

type testStepServerInfoAPIResults struct {
	readResult   *sacloud.Server
	readError    error
	vncInfo      *sacloud.VNCProxyResponse
	vncInfoError error
}

func (t *testStepServerInfoAPIResults) init() {
	t.readResult = nil
	t.readError = nil
	t.vncInfo = nil
	t.vncInfoError = nil
}

var testStepServerInfoResult = &testStepServerInfoAPIResults{}

func TestStepServerInfo(t *testing.T) {
	ctx := context.Background()

	t.Run("with reading server info error", func(t *testing.T) {
		testStepServerInfoResult.init()

		step := &stepServerInfo{}
		state := initStepServerInfoState()

		testStepServerInfoResult.readError = errors.New("error")

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("with reading server info error", func(t *testing.T) {
		testStepServerInfoResult.init()

		step := &stepServerInfo{}
		state := initStepServerInfoState()

		testStepServerInfoResult.readError = errors.New("error")

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("with empty server info", func(t *testing.T) {
		testStepServerInfoResult.init()

		step := &stepServerInfo{}
		state := initStepServerInfoState()

		testStepServerInfoResult.readResult = dummyServerEmptyInfo()
		testStepServerInfoResult.vncInfo = dummyVNCProxyInfo()

		action := step.Run(ctx, state)
		_, hasError := state.GetOk("error")

		assert.Equal(t, multistep.ActionContinue, action)
		assert.False(t, hasError)

		keys := []string{
			"server_ip",
			"default_route",
			"network_mask_len",
			"dns1",
			"dns2",
			"vnc",
		}
		for _, key := range keys {
			_, ok := state.GetOk(key)
			assert.Truef(t, ok, "%s is not exists in StateBag", key)
		}
	})

	t.Run("with reading VPCInfo error", func(t *testing.T) {
		testStepServerInfoResult.init()

		step := &stepServerInfo{}
		state := initStepServerInfoState()

		testStepServerInfoResult.readResult = dummyServerEmptyInfo()
		testStepServerInfoResult.vncInfoError = errors.New("error")

		action := step.Run(ctx, state)
		err := state.Get("error").(error)

		assert.Equal(t, multistep.ActionHalt, action)
		assert.Error(t, err)
	})

	t.Run("normal case", func(t *testing.T) {
		testStepServerInfoResult.init()

		step := &stepServerInfo{}
		state := initStepServerInfoState()

		expectVNCInfo := dummyVNCProxyInfo()
		testStepServerInfoResult.readResult = dummyServerInfo()
		testStepServerInfoResult.vncInfo = expectVNCInfo

		action := step.Run(ctx, state)
		_, hasError := state.GetOk("error")

		assert.Equal(t, multistep.ActionContinue, action)
		assert.False(t, hasError)

		assert.Equal(t, dummyServerIP, state.Get("server_ip").(string))
		assert.Equal(t, dummyServerDefaultRoute, state.Get("default_route").(string))
		assert.Equal(t, dummyServerNwMaskLen, state.Get("network_mask_len").(int))
		assert.Equal(t, dummyDNSServers[0], state.Get("dns1").(string))
		assert.Equal(t, dummyDNSServers[1], state.Get("dns2").(string))

		actualVNCInfo, ok := state.GetOk("vnc")
		assert.True(t, ok)
		assert.EqualValues(t, expectVNCInfo, actualVNCInfo.(*sacloud.VNCProxyResponse))
	})
}

func initStepServerInfoState() multistep.StateBag {
	state := dummyMinimumStateBag(nil)
	state.Put("server_id", dummyServerID)

	state.Put("serverClient", &dummyServerClient{
		readFunc: func(int64) (*sacloud.Server, error) {
			return testStepServerInfoResult.readResult, testStepServerInfoResult.readError
		},
		getVNCProxyFunc: func(int64) (*sacloud.VNCProxyResponse, error) {
			return testStepServerInfoResult.vncInfo, testStepServerInfoResult.vncInfoError
		},
	})

	return state
}

func dummyServerInfo() *sacloud.Server {
	server := &sacloud.Server{
		Resource: sacloud.NewResource(dummyServerID),
	}
	nic := sacloud.Interface{}
	nic.IPAddress = dummyServerIP
	nic.Switch = &sacloud.Switch{}
	nic.Switch.Subnet = &sacloud.Subnet{
		DefaultRoute:   dummyServerDefaultRoute,
		NetworkMaskLen: dummyServerNwMaskLen,
	}
	server.Interfaces = []sacloud.Interface{nic}
	server.Zone = &sacloud.Zone{}
	server.Zone.Region = &sacloud.Region{
		NameServers: dummyDNSServers,
	}
	return server
}

func dummyServerEmptyInfo() *sacloud.Server {
	server := &sacloud.Server{
		Resource: sacloud.NewResource(dummyServerID),
	}
	server.Zone = &sacloud.Zone{}
	server.Zone.Region = &sacloud.Region{}

	return server
}

func dummyVNCProxyInfo() *sacloud.VNCProxyResponse {
	return &sacloud.VNCProxyResponse{
		Status:       "status",
		Host:         "ftps.example.com",
		IOServerHost: "io.example.com",
		Port:         "31313",
		Password:     "password",
		VNCFile:      "vncfile",
	}
}
