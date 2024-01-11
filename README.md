# Envoy 配置管理

基于自己对XDS的理解，做一个envoy配置管理下发的后台程序。
主要目的是提供一个方便的管理界面管理envoy配置，在将envoy作为一个独立的网关入口（非k8s集群）时，能够方便快捷的对envoy进行配置管理。

## 目前的思路

control-plane，一个GRPC服务基于go-control-plane SDK 做配置加载、下发功能。

manage-plane, 一个API服务，采用gin,gorm 做配置管理。

control-plane与manage-plane 通过redis发布订阅进行变更通知

vueadmin ,一个纯前端项目，做界面。
ps 前端水平比较差，不做交互式的开发，仅使用yaml编辑器对yaml进行编辑

## 表设计

envoy_cluster  集群表，不是cds中的cluster，是物理环境中集群，例如 dev、test、prod。
envoy_node     节点，物理环境中的envoy节点表
cds            cluster配置表
eds            endpoint配置表
lds            listener配置表
rds            route配置表
vhds           virtualhost配置表
ecds           扩展配置表
sds            secret配置表 
rls            限流配置表 基于ratelimit服务



