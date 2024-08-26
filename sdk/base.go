package fastsearch

import (
	"runtime"

	"gitee.com/rachel_os/fastsearch/searcher/model"
	"gitee.com/rachel_os/fastsearch/searcher/system"
)

// Query 查询
func (c *Client) Query(req *model.SearchRequest) (any, error) {
	r, err := c.container.GetDataBase(req.Database).MultiSearch(req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (*Client) GC() {
	runtime.GC()
}
func (c *Client) Status() (map[string]interface{}, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// TODO 其他系统信息
	r := map[string]interface{}{
		"memory": system.GetMemStat(),
		"cpu":    system.GetCPUStatus(),
		"disk":   system.GetDiskStat(),
	}
	return r, nil
}
