// Package config
//
// @author: xwc1125
package config

import (
	"sync"
	"sync/atomic"

	"github.com/chain5j/chain5j-pkg/cli"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-protocol/models"
	"github.com/chain5j/chain5j-protocol/protocol"
	"github.com/chain5j/logger"
	"github.com/chain5j/logger/zap"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	_ protocol.Config = new(config)
)

var (
	// defaultTxSize 默认的最大交易size
	defaultTxSize = types.StorageSize(128 * 1024)
)

type config struct {
	log        logger.Logger           // log
	configFile string                  // 本地配置文件
	db         protocol.DatabaseReader // 数据库

	chainConfig     *models.ChainConfig // 链配置
	localConfig     *sync.Pool          // 节点本地启动配置,*models.LocalConfig
	chainConfigTemp atomic.Value        // chainConfig缓存
	txSize          atomic.Value        // txSize的缓存
	genesisBlock    atomic.Value        // 创世块的缓存
}

// NewConfig 加载本地配置文件
func NewConfig(configFile string, opts ...option) (protocol.Config, error) {
	c := &config{
		log:        logger.New("config"),
		configFile: configFile,
	}
	if err := apply(c, opts...); err != nil {
		c.log.Error("apply is error", "err", err)
		return nil, err
	}
	var localConfig = new(models.LocalConfig)
	// 从config file中加载配置
	err := cli.LoadConfig(configFile, "", localConfig, func(e fsnotify.Event) {
		c.log.Info("Config file changed", "eventName", e.Name)
		if err := cli.LoadConfig(configFile, "", localConfig, nil); err == nil {
			c.localConfig.Put(localConfig)
		}
	})
	if err != nil {
		c.log.Error("load config err", "config", configFile, "err", err)
		return nil, err
	}

	c.localConfig = &sync.Pool{
		New: func() interface{} {
			return localConfig
		},
	}
	// 初始化日志
	zap.InitWithConfig(&localConfig.Log)

	return c, err
}

// SetDatabase 设置数据库
func (c *config) SetDatabase(db protocol.DatabaseReader) error {
	c.db = db
	chainConfig, err := c.db.ChainConfig()
	if err != nil {
		return err
	}
	c.chainConfig = chainConfig
	return nil
}

func (c *config) getChainConfig() *models.ChainConfig {
	if c.db != nil {
		var err error
		c.chainConfig, err = c.db.ChainConfig()
		if err != nil {
			c.log.Error("db getChainConfig err", "err", err)
			return nil
		}
	}
	if c.chainConfig.StateApp == nil {
		c.chainConfig.StateApp = &models.StateAppConfig{
			UseEthereum: false,
		}
	}
	return c.chainConfig
}

// ChainConfig 获取链配置
func (c *config) ChainConfig() models.ChainConfig {
	if c.chainConfig == nil {
		if c.getChainConfig() == nil {
			return models.ChainConfig{}
		}
	}
	// 需要深度copy
	chainConfig := c.copyChainConfig()
	return chainConfig
}

// copyChainConfig chainConfig 复制
func (c *config) copyChainConfig() models.ChainConfig {
	if c.chainConfigTemp.Load() != nil {
		return c.chainConfigTemp.Load().(models.ChainConfig)
	}
	packer := *c.chainConfig.Packer
	chainConfig := models.ChainConfig{
		GenesisHeight: c.chainConfig.GenesisHeight,
		ChainID:       c.chainConfig.ChainID,
		ChainName:     c.chainConfig.ChainName,
		VersionName:   c.chainConfig.VersionName,
		VersionCode:   c.chainConfig.VersionCode,
		TxSizeLimit:   c.chainConfig.TxSizeLimit,
		Packer:        &packer,
		Consensus:     c.chainConfig.Consensus,
		StateApp:      c.chainConfig.StateApp,
	}
	c.chainConfigTemp.Store(chainConfig)
	return chainConfig
}

// GenesisBlock 获取创世块
func (c *config) GenesisBlock() *models.Block {
	if c.genesisBlock.Load() != nil {
		return c.genesisBlock.Load().(*models.Block)
	}
	genesisBlock, err := c.db.GetBlockByHeight(c.chainConfig.GenesisHeight)
	if err != nil {
		c.log.Crit("get genesis block err", "err", err)
	}
	c.genesisBlock.Store(genesisBlock)
	return genesisBlock
}

// TxSizeLimit 获取交易池本地配置
func (c *config) TxSizeLimit() types.StorageSize {
	if c.txSize.Load() != nil {
		return c.txSize.Load().(types.StorageSize)
	}
	if c.chainConfig == nil {
		c.txSize.Store(defaultTxSize)
		return defaultTxSize
	}
	size := types.StorageSize(c.chainConfig.TxSizeLimit * 1024)
	c.txSize.Store(size)
	return size
}

func (c *config) DatabaseConfig() models.DatabaseConfig {
	return c.localConfig.Get().(*models.LocalConfig).Database
}
func (c *config) BlockchainConfig() models.BlockchainLocalConfig {
	return c.LocalConfig().Blockchain
}

func (c *config) LocalConfig() *models.LocalConfig {
	return c.localConfig.Get().(*models.LocalConfig)
}

// TxPoolConfig 交易池本地配置
func (c *config) TxPoolConfig() models.TxPoolLocalConfig {
	return c.localConfig.Get().(*models.LocalConfig).TxPool
}

// NodeKeyConfig nodeKey本地配置
func (c *config) NodeKeyConfig() models.NodeKeyLocalConfig {
	return c.localConfig.Get().(*models.LocalConfig).NodeKey
}

// PackerConfig packer本地配置
func (c *config) PackerConfig() models.PackerLocalConfig {
	return c.LocalConfig().Packer
}

// EnablePacker 是否启动打包器
func (c *config) EnablePacker() bool {
	return viper.GetBool("mine")
}

func (c *config) BroadcasterConfig() models.BroadcasterLocalConfig {
	return c.LocalConfig().Broadcaster
}

func (c *config) P2PConfig() models.P2PConfig {
	return c.LocalConfig().P2P
}

func (c *config) ConsensusConfig() models.ConsensusLocalConfig {
	return c.LocalConfig().Consensus
}
