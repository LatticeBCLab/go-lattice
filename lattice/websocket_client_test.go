package lattice

import (
	"context"
	"testing"
	"time"

	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/lattice/client"
	"github.com/stretchr/testify/assert"
)

func setupWebsocketClient() client.WebSocketApi {
	initWebsocketClientArgs := &client.WebSocketApiInitParam{
		WebSocketUrl: "ws://192.168.0.191:50000",
	}
	wsApi := client.NewWebSocketApi(initWebsocketClientArgs)

	return wsApi
}

func TestWebsocketClientSubscribe(t *testing.T) {
	wsApi := setupWebsocketClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	t.Run("subscribe monitorData", func(t *testing.T) {
		subMonitorData, err := wsApi.Subscribe(ctx, "latc_subscribe", "monitorData")
		assert.NoError(t, err)
		assert.NotNil(t, subMonitorData)
		for i := 0; i < 5; i++ { // 读5条数据然后关闭
			monitorData, err := subMonitorData.Read()
			assert.NoError(t, err)
			assert.NotNil(t, monitorData)
			t.Log(i, monitorData)
		}
		err = subMonitorData.Close()
		assert.NoError(t, err)
	})
}

func TestWorkflow(t *testing.T) {
	wsApi := setupWebsocketClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	t.Run("subscribe workflow", func(t *testing.T) {
		subWorkflow, err := wsApi.Workflow(ctx, nil)
		assert.Nil(t, err)
		assert.NotNil(t, subWorkflow)
		for i := range 5 { // 读5条数据然后关闭
			workflow, err := subWorkflow.Read()
			assert.NoError(t, err)
			assert.NotNil(t, workflow)
			t.Logf("%d: %#v", i, workflow)
		}
		err = subWorkflow.Close()
		assert.NoError(t, err)
	})
	t.Run("subscribe workflow daemon block", func(t *testing.T) {
		subWorkflow, err := wsApi.Workflow(ctx, &types.WorkflowSubscribeCondition{Type: types.WorkflowType_DAEMON_BLOCK})
		assert.Nil(t, err)
		assert.NotNil(t, subWorkflow)
		for i := range 5 { // 读5条数据然后关闭
			workflow, err := subWorkflow.Read()
			assert.NoError(t, err)
			assert.NotNil(t, workflow)
			t.Logf("%d: %#v", i, workflow)
		}
		err = subWorkflow.Close()
		assert.NoError(t, err)
	})
}
