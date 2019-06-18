package hcrpcclient

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcjson"
	"github.com/HcashOrg/hcd/wire"
)


type FutureSendInstantRawTransactionResult chan *response
type FutureSendInstantTxVoteResult chan *response

func (r FutureSendInstantRawTransactionResult) Receive() (*chainhash.Hash, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal result as a string.
	var txHashStr string
	err = json.Unmarshal(res, &txHashStr)
	if err != nil {
		return nil, err
	}

	return chainhash.NewHashFromStr(txHashStr)
}

func (r FutureSendInstantTxVoteResult) Receive()  error {
	_, err := receiveFuture(r)
	if err != nil {
		return  err
	}

	return nil
}

func (c *Client) SendInstantRawTransaction(tx *wire.MsgInstantTx, allowHighFees bool) (*chainhash.Hash, error) {
	return c.SendInstantRawTransactionAsync(tx, allowHighFees).Receive()
}

func (c *Client) SendInstantRawTransactionAsync(tx *wire.MsgInstantTx, allowHighFees bool) FutureSendInstantRawTransactionResult {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return newFutureError(err)
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}

	cmd := hcjson.NewSendInstantRawTransactionCmd(txHex, &allowHighFees)
	return c.sendCmd(cmd)
}



func (c *Client)SendInstantTxVote(instantTxVote *wire.MsgInstantTxVote)error{
	return c.SendInstantTxVoteAsync(instantTxVote).Receive()
}

func (c *Client) SendInstantTxVoteAsync(msgInstantTxVote *wire.MsgInstantTxVote) FutureSendInstantTxVoteResult {
	txvoteHex := ""
	if msgInstantTxVote != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, msgInstantTxVote.SerializeSize()))
		if err := msgInstantTxVote.Serialize(buf); err != nil {
			return newFutureError(err)
		}
		txvoteHex = hex.EncodeToString(buf.Bytes())
	}

	cmd := hcjson.NewSendInstantTxVoteCmd(txvoteHex)
	return c.sendCmd(cmd)
}

