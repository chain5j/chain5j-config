## chain5j-config配置中心

## 简介
`chain5j-config` chain5j链配置中心。

## 功能
1. 提供`chainConfig`和`localConfig`配置的获取
2. 热加载`localConfig`，能够让本地配置自动更新

### 使用说明

```go
node.config, err = config.Factory{}.NewFactory(
	node.configFile,
	config.WithDB(node.database), 
	)
```

## 证书
`chain5j-config` 的源码允许用户在遵循 [Apache 2.0 开源证书](LICENSE) 规则的前提下使用。

## 版权
Copyright@2022 chain5j

![chain5j](./chain5j.png)

