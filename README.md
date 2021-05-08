# tradingdb2

``tradingdb2``是新的交易系统最核心的服务，前端分析端（jupyter）会通过``tradingdb2``获取分析数据，并通过``tradingcore2``节点进行分布式运算。

# 版本更新

### v0.8 （计划中）

- 移除``SimTradingDB``模块

### v0.7 （计划中）

- 移除 ``tradingcore2`` v0.5 的支持
- 移除 ``trdb2py`` v0.2 的支持
- 依然兼容本地``SimTradingDB``数据

### v0.6

- 结构调整
- 内存峰值不会随着任务数量增加而无限制上涨
- 运行效率提升
- 支持新的``SimTradingDB2``模块，兼容``SimTradingDB``模块
- 兼容 ``tradingcore2`` v0.5、v0.6 2个版本的运算节点
- 兼容 ``trdb2py`` v0.2、v0.3

### v0.5

- 基本功能完整
- 搭配 ``tradingcore2`` v0.5版节点工作
- 可以用 ``trdb2py`` v0.2做前端分析用
