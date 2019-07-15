package hcrpcclient

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcjson"
	"github.com/HcashOrg/hcd/wire"
)


type FutureSendAiRawTransactionResult chan *response
type FutureSendAiTxVoteResult chan *response

func (r FutureSendAiRawTransactionResult) Receive() (*chainhash.Hash, error) {
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

func (r FutureSendAiTxVoteResult) Receive()  error {
	_, err := receiveFuture(r)
	if err != nil {
		return  err
	}

	return nil
}

func (c *Client) SendAiRawTransaction(tx *wire.MsgAiTx, allowHighFees bool) (*chainhash.Hash, error) {
	return c.SendAiRawTransactionAsync(tx, allowHighFees).Receive()
}

func (c *Client) SendAiRawTransactionAsync(tx *wire.MsgAiTx, allowHighFees bool) FutureSendAiRawTransactionResult {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return newFutureError(err)
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}

	cmd := hcjson.NewSendAiRawTransactionCmd(txHex, &allowHighFees)
	return c.sendCmd(cmd)
}



func (c *Client)SendAiTxVote(aiTxVote *wire.MsgAiTxVote)error{
	return c.SendAiTxVoteAsync(aiTxVote).Receive()
}

func (c *Client) SendAiTxVoteAsync(msgAiTxVote *wire.MsgAiTxVote) FutureSendAiTxVoteResult {
	txvoteHex := ""
	if msgAiTxVote != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, msgAiTxVote.SerializeSize()))
		if err := msgAiTxVote.Serialize(buf); err != nil {
			return newFutureError(err)
		}
		txvoteHex = hex.EncodeToString(buf.Bytes())
	}

	cmd := hcjson.NewSendAiTxVoteCmd(txvoteHex)
	return c.sendCmd(cmd)
}

