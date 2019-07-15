package hcrpcclient

import (
	"encoding/json"
	"github.com/HcashOrg/hcd/hcjson"
)

type FutureGetTxlockPoolInfoResult chan *response

func (c *Client) GetItTxInLockPool() (map[string]*hcjson.TxLockInfo, error) {
	return c.GetItTxInLockPoolAsync().Receive()
}

func (c *Client) GetItTxInLockPoolAsync() FutureGetTxlockPoolInfoResult {

	cmd := hcjson.NewGetTxlockPoolInfoCmd()
	return c.sendCmd(cmd)
}

// Receive waits for the response promised by the future and returns a
// transaction given its hash.
func (r FutureGetTxlockPoolInfoResult) Receive() (map[string]*hcjson.TxLockInfo, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	var lockInfo hcjson.GetTxLockpoolInfoResult
	err = json.Unmarshal(res, &lockInfo)
	if err != nil {
		return nil, err
	}

	return lockInfo.Info, nil
}
