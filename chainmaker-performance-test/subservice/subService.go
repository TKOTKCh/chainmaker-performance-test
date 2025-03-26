/*
Copyright (C) 2023 Beijing Advanced Innovation Center for Future Blockchain and Privacy Computing (未来区块链与隐私计算高精尖创新中心). All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package subservice

import (
	logger "chain-performance-test/log"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"context"
	"errors"
	"time"
)

var BlockInfos []BlockInfo

// BlockInfo 区块信息定义
type BlockInfo struct {
	Timestamp   time.Time
	BlockHeight uint64
	TxCount     uint32
}

// SubscribeBlock 订阅链定义
type SubscribeBlock struct {
	client  *sdk.ChainClient
	closeCh chan struct{}
}

// BlockSubscriber 订阅链接口
type BlockSubscriber interface {
	Run(ctx context.Context, startBlock, endBlock int64) error
	Close() error
	SetClient(c *sdk.ChainClient)
}

func (s *SubscribeBlock) SetClient(c *sdk.ChainClient) {
	if c != nil {
		s.client = c
	}
	s.closeCh = make(chan struct{})
}

// Run 运行消息订阅服务
func (s *SubscribeBlock) Run(ctx context.Context, startBlock, endBlock int64) error {

	logger.Logger.Println("start sub")
	var Info BlockInfo
	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)

	ch, _ := s.client.SubscribeBlock(ctx, startBlock, endBlock, false, false)
	for {
		select {
		case <-s.closeCh:
			// close关闭服务
			return nil
		case <-ctx.Done():
			// 超时关闭消息订阅服务
			return errors.New("subscribe time out")
		case <-timer.C:
			// 5秒内ch中没有信息，触发取消操作
			logger.Logger.Println("5秒内无消息，执行取消订阅操作")
			ctx.Done()
			return nil
		case blockInfoTemp, ok := <-ch:
			// 获得到新的落块信息
			if !ok {
				logger.Logger.Println("chainConfig Changed, re-subscribe ...... ")
				return nil
			}
			if blockInfoTemp == nil {
				logger.Logger.Println("require not nil")
				return nil
			}
			blockInfo, ok := blockInfoTemp.(*common.BlockInfo)
			if !ok {
				logger.Logger.Println("require true")
			}
			logger.Logger.Printf("time is %v, block %v with %v txs, \n", blockInfo.Block.GetTimestamp(), blockInfo.Block.Header.BlockHeight, blockInfo.Block.Header.TxCount)
			// 存储区块信息（落块时间、区块高度、区块交易数）到 BlockInfos
			Info.Timestamp = blockInfo.Block.GetTimestamp()
			Info.BlockHeight = blockInfo.Block.Header.BlockHeight
			Info.TxCount = blockInfo.Block.Header.TxCount
			BlockInfos = append(BlockInfos, Info)

			// 重置定时器
			timer.Reset(timeout)
		}
	}

}

// Close 关闭消息订阅服务
func (s *SubscribeBlock) Close() error {
	close(s.closeCh)
	return nil
}
