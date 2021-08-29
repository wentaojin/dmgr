### DMGR 数据迁移管理平台

MySQL -> TiDB DM 数据迁移任务管理，已集成 [DM](https://docs.pingcap.com/zh/tidb-data-migration/stable/overview) 集群级别管理功能

1. 集群部署
2. 集群启停
3. 集群状态查看
4. 集群扩缩容
5. 集群滚更
6. 集群补丁
7. 集群升级
8. 集群销毁

### DMGR 运行示例

```
# dmgr.toml 配置文件示例参见 conf 目录
$ ./dmgr --config dmgr.toml

# 程序默认运行 debug 模式，展示 api 接口，若无需展示，请切换至 release 模式运行
$ export GIN_MODE=release
$ ./dmgr --config dmgr.toml
```

数据同步任务功能待设计阶段...

### DMGR 集群管理级别功能设计

#### dmgr 管理表结构设计
- 部署机器列表新增(部署、升级前必备)  -> machine

  * 机器新增用户链接信息，需要 root 用户或者具备 sudo 权限的用户 
  * 离线包上传(部署、升级前必备) -> warehouse
  * 集群运维管理
  - 保存集群元信息 -> cluster_meta
  - 保存集群拓扑  -> cluster_topology
  - 用户登录       -> user

#### dmgr 集群管理目录层级设计

```
集群管理元目录层级 {cluster_path}/cluster/{cluster_name}
v2.0.1               -> 离线安装包解压后的存放目录，版本号区分
cache                -> 模板文件生成文件
ssh                  -> 集群 ssh 认证存放路径
```

#### dmgr 集群部署目录层级设计

```
dm 集群常见部署目录层级
{deploy_dir}/bin 
{deploy_dir}/scripts
{deploy_dir}/conf
{data_dir}/data
{log_dir}/log
```

#### dm 安装包结构设计

```
dm-v2.0.1.tar.gz 压缩包内容【格式必须】
  bin                     -> 二进制文件
    * prometheus
    * alertmanager
    * dm-worker
    * dm-master

  conf
    * dm_worker.rules.yml -> 告警规则模板文件
    * alertmanager.yml    -> 配置文件
    * dm-master.toml      -> 默认配置文件
    * dm-worker.toml      -> 默认配置文件
    
  template               -> 非 template 生成的文件支持手工修改
    * grafana.ini.tmpl    -> grafana 模板文件
      - dashboard.yml.tmpl -> 监控面板模板
      - datasource.yml.tmpl -> 监控数据源模板
    * run_grafana.sh.tmpl -> 运行脚本模板
    
    * systemd.service.tmpl -> 二进制文件运行 systemd 服务模板文件

    
    * prometheus.yml.tmpl -> 配置模板文件
    * run_prometheus.sh.tmpl -> 运行脚本模板
    
    * run_alertmanager.sh.tmpl -> 运行脚本模板
    
    * run_dm-master.sh.tmpl -> 运行脚本模板   
    * run_dm-worker.sh.tmpl -> 运行脚本模板
    * run_dm-master-scale.sh.tmpl -> 扩容 dm-master 脚本模板  [扩容阶段]  


  
  grafana.tar.gz
    * bin
      & grafana-server     -> 二进制文件
      & grafana-cli        -> 二进制文件
    
    * dashboards   -> 存放 grafana 监控面板原始 json 文件模板   
      & dm_instances.json
      & dm.json
    
    * plugins/            -> 存放插件目录
    
    * provisioning/      -> 包含 grafana 将在启动和运行时应用的配置文件的文件夹
      & dashboards/      -> 存放 dashboard 层面目录  dashboard.yml.tmpl
      & datasources/     -> 存放监控数据源 datasource.yml.tmpl
      & conf/            -> 存放 grafana 默认 defaults.ini
    
    * homepath/           -> 二进制文件运行家目录
        * public/         -> 存放所有
        * scripts/        -> 存放脚本
        * notifiers/      -> 存放其他信息

  others
    * dmctl             -> 二进制文件
    * task_advanced.yml  -> 任务同步示例文件
    * task_basic.yml     -> 任务同步示例文件
```

