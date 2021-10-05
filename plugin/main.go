package main

import (
	"encoding/json"
	"github.com/hashicorp/go-plugin"
	"github.com/wulie/go-plugin-bidirectional/shared"
	"io/ioutil"
)

type Counter struct {
}

type data struct {
	Value int64
}

func (c Counter) Put(key string, value int64, a shared.Add) error {
	get, err := c.Get(key)
	if err != nil {
		return err
	}
	sum, err := a.Sum(get, value)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(&data{Value: sum})
	if err != nil {
		return err
	}
	//err = ioutil.WriteFile("aaa"+key, bytes, 0644)
	return ioutil.WriteFile("kv_"+key, bytes, 0644)
}

//func (c *Counter) Put(key string, value int64, a shared.Add) error {
//	v, _ := c.Get(key)
//
//	r, err := a.Sum(v, value)
//	if err != nil {
//		return err
//	}
//
//	buf, err := json.Marshal(&data{r})
//	if err != nil {
//		return err
//	}
//
//	return ioutil.WriteFile("kv_"+key, buf, 0644)
//}

func (c *Counter) Get(key string) (int64, error) {
	bytes, err := ioutil.ReadFile("kv_" + key)
	if err != nil {
		return 0, err
	}
	d := &data{}
	err = json.Unmarshal(bytes, d)
	if err != nil {
		return 0, err
	}
	return d.Value, nil
}

//func (c *Counter) Get(key string) (int64, error) {
//	dataRaw, err := ioutil.ReadFile("kv_" + key)
//	if err != nil {
//		return 0, err
//	}
//
//	data := &data{}
//	err = json.Unmarshal(dataRaw, data)
//	if err != nil {
//		return 0, err
//	}
//
//	return data.Value, nil
//}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"counter": &shared.CounterPlugin{
				NetRPCUnsupportedPlugin: plugin.NetRPCUnsupportedPlugin{},
				Impl:                    &Counter{},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
