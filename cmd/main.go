package main

import (
	"fmt"
	"os"

	"github.com/trmigor/distr-model/internal/messages"
	"github.com/trmigor/distr-model/internal/process"
	"github.com/trmigor/distr-model/internal/world"
	"github.com/trmigor/distr-model/user/context"
)

func workFunctionSETX(dp *process.Process, m *messages.Message) bool {
	m.Ptr = 0
	s := m.GetString()
	nl := dp.Network
	if !dp.IsMyMessage([]byte("SETX"), s) {
		return false
	}
	if string(s) == "SETX_INIT" {
		arg := m.GetInt32()
		nl.SendMessage(dp.Node, dp.Node, messages.NewMessageByArgs(messages.NewMessageArg([]byte("SETX_SET")), messages.NewMessageArg(arg)))
	} else if string(s) == "SETX_SET" {
		arg := m.GetInt32()
		fmt.Printf("[%v]: SETX_SET received, arg=%v\n", dp.Node, arg)
		if dp.Context["SetX"].(context.SetX).X != int(arg) {
			for v := range *dp.Neibs() {
				nl.SendMessage(dp.Node, v.(int32), messages.NewMessageByArgs(messages.NewMessageArg([]byte("SETX_SET")), messages.NewMessageArg(arg)))
			}
		}
		ctx := dp.Context["SetX"].(context.SetX)
		ctx.X = int(arg)
		dp.Context["SetX"] = ctx
	}
	return true
}

func main() {
	args := os.Args[1:]
	config := "configs/config.data"
	if len(args) > 0 {
		config = args[0]
	}

	w := world.New()
	defer w.Stop()
	w.RegisterWorkFunction([]byte("SETX"), workFunctionSETX)
	if !w.ParseConfig([]byte(config)) {
		fmt.Printf("can't open file '%v'\n", config)
		os.Exit(1)
	}
}
