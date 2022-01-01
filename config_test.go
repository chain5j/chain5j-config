// Package config
//
// @author: xwc1125
package config

import (
	"github.com/chain5j/chain5j-protocol/mock"
	"github.com/chain5j/chain5j-protocol/models"
	"github.com/chain5j/logger"
	"github.com/chain5j/logger/zap"
	"github.com/golang/mock/gomock"
	"testing"
)

func init() {
	zap.InitWithConfig(&logger.LogConfig{
		Console: logger.ConsoleLogConfig{
			Level:    4,
			Modules:  "*",
			ShowPath: false,
			Format:   "",
			UseColor: true,
			Console:  true,
		},
		File: logger.FileLogConfig{},
	})
}

func TestConfig(t *testing.T) {
	mockCtl := gomock.NewController(nil)
	mockDbReader := mock.NewMockDatabaseReader(mockCtl)
	mockDbReader.EXPECT().ChainConfig().Return(&models.ChainConfig{
		GenesisHeight: 0,
		ChainID:       1,
		ChainName:     "chain5j",
		VersionName:   "v0.0.1",
		VersionCode:   1,
		TxSizeLimit:   1024,
		Packer: &models.PackerConfig{
			WorkerType:           0,
			BlockMaxTxsCapacity:  10000,
			BlockMaxSize:         2048,
			BlockMaxIntervalTime: 1000,
			BlockGasLimit:        50000000,
			Period:               3000,
			EmptyBlocks:          0,
			Timeout:              100,
		},
		StateApp: &models.StateAppConfig{
			UseEthereum: false,
		},
	}, nil).AnyTimes()

	config, err := NewConfig(
		"../../conf/config.yaml",
		WithDB(mockDbReader),
	)
	if err != nil {
		t.Fatal(err)
	}

	txPoolConfig := config.TxPoolConfig()
	t.Logf("txPoolConfig=%v", txPoolConfig)
	txSize := config.TxSizeLimit()
	t.Logf("txSize=%v", txSize)
	nodeKeyConfig := config.NodeKeyConfig()
	t.Logf("nodeKeyConfig=%v", nodeKeyConfig)
	packer := config.EnablePacker()
	t.Logf("packer=%v", packer)
	chainConfig := config.ChainConfig()
	t.Logf("%v", chainConfig)
}
