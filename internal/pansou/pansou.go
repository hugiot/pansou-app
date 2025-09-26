package pansou

import (
	"fmt"
	"log"
	"pansou/config"
	"pansou/model"
	"pansou/plugin"
	"pansou/service"
	"pansou/util"
	"pansou/util/cache"
	"pansou/util/json"
	"time"

	// 以下是插件的空导入，用于触发各插件的init函数，实现自动注册
	// 添加新插件时，只需在此处添加对应的导入语句即可
	// _ "pansou/plugin/hdr4k"
	// _ "pansou/plugin/pan666"
	_ "pansou/plugin/bixin"
	_ "pansou/plugin/cldi"
	_ "pansou/plugin/clmao"
	_ "pansou/plugin/clxiong"
	_ "pansou/plugin/cyg"
	_ "pansou/plugin/ddys"
	_ "pansou/plugin/duoduo"
	_ "pansou/plugin/erxiao"
	_ "pansou/plugin/fox4k"
	_ "pansou/plugin/haisou"
	_ "pansou/plugin/hdmoli"
	_ "pansou/plugin/huban"
	_ "pansou/plugin/hunhepan"
	_ "pansou/plugin/javdb"
	_ "pansou/plugin/jikepan"
	_ "pansou/plugin/jutoushe"
	_ "pansou/plugin/labi"
	_ "pansou/plugin/leijing"
	_ "pansou/plugin/libvio"
	_ "pansou/plugin/miaoso"
	_ "pansou/plugin/muou"
	_ "pansou/plugin/ouge"
	_ "pansou/plugin/pansearch"
	_ "pansou/plugin/panta"
	_ "pansou/plugin/panwiki"
	_ "pansou/plugin/panyq"
	_ "pansou/plugin/pianku"
	_ "pansou/plugin/qupansou"
	_ "pansou/plugin/sdso"
	_ "pansou/plugin/shandian"
	_ "pansou/plugin/susu"
	_ "pansou/plugin/thepiratebay"
	_ "pansou/plugin/u3c3"
	_ "pansou/plugin/wanou"
	_ "pansou/plugin/wuji"
	_ "pansou/plugin/xb6v"
	_ "pansou/plugin/xdyh"
	_ "pansou/plugin/xiaoji"
	_ "pansou/plugin/xiaozhang"
	_ "pansou/plugin/xuexizhinan"
	_ "pansou/plugin/xys"
	_ "pansou/plugin/yuhuage"
	_ "pansou/plugin/zhizhen"
)

type Pansou interface {
	// Health 健康检查
	Health(req HealthRequest) HealthResponse
	// Search 搜索
	Search(req model.SearchRequest) (SearchResponse, error)
	// Release 释放资源
	Release() error
}

type pansouImpl struct {
	searchService     *service.SearchService
	cacheWriteManager *cache.DelayedBatchWriteManager
}

func (p *pansouImpl) Health(_ HealthRequest) HealthResponse {
	// 根据配置决定是否返回插件信息
	pluginCount := 0
	var pluginNames []string
	pluginsEnabled := config.AppConfig.AsyncPluginEnabled

	if pluginsEnabled && p.searchService != nil && p.searchService.GetPluginManager() != nil {
		plugins := p.searchService.GetPluginManager().GetPlugins()
		pluginCount = len(plugins)
		for _, p := range plugins {
			pluginNames = append(pluginNames, p.Name())
		}
	}

	// 获取频道信息
	channels := config.AppConfig.DefaultChannels
	channelsCount := len(channels)

	response := HealthResponse{
		Status:         "ok",
		PluginsEnabled: pluginsEnabled,
		Channels:       channels,
		ChannelsCount:  channelsCount,
		PluginCount:    0,
		Plugins:        nil,
	}

	// 只有当插件启用时才返回插件相关信息
	if pluginsEnabled {
		response.PluginCount = pluginCount
		response.Plugins = pluginNames
	}

	return response
}

func (p *pansouImpl) Search(req model.SearchRequest) (SearchResponse, error) {
	response, err := p.searchService.Search(req.Keyword, req.Channels, req.Concurrency, req.ForceRefresh, req.ResultType, req.SourceType, req.Plugins, req.CloudTypes, req.Ext)
	if err != nil {
		return SearchResponse{}, err
	}
	return p.convResponse(response), nil
}

func (p *pansouImpl) Release() error {
	// 优先保存缓存数据到磁盘（数据安全第一）
	// 增加关闭超时时间，确保数据有足够时间保存
	shutdownTimeout := 10 * time.Second

	if p.cacheWriteManager != nil {
		if err := p.cacheWriteManager.Shutdown(shutdownTimeout); err != nil {
			log.Printf("缓存数据保存失败: %v", err)
			return err
		}
	}

	// 额外确保内存缓存也被保存（双重保障）
	if mainCache := service.GetEnhancedTwoLevelCache(); mainCache != nil {
		if err := mainCache.FlushMemoryToDisk(); err != nil {
			log.Printf("内存缓存同步失败: %v", err)
			return err
		}
	}

	fmt.Println("服务器已安全关闭")
	return nil
}

func (p *pansouImpl) convResponse(response model.SearchResponse) SearchResponse {
	bs, _ := json.Marshal(response)
	var result SearchResponse
	_ = json.Unmarshal(bs, &result)
	return result
}

func New() Pansou {
	instance := &pansouImpl{}

	// 初始化配置
	config.Init()

	// 初始化HTTP客户端
	util.InitHTTPClient()

	// 初始化缓存写入管理器
	var err error
	instance.cacheWriteManager, err = cache.NewDelayedBatchWriteManager()
	if err != nil {
		log.Fatalf("缓存写入管理器创建失败: %v", err)
	}
	if err := instance.cacheWriteManager.Initialize(); err != nil {
		log.Fatalf("缓存写入管理器初始化失败: %v", err)
	}
	// 将缓存写入管理器注入到service包
	service.SetGlobalCacheWriteManager(instance.cacheWriteManager)

	// 延迟设置主缓存更新函数，确保service初始化完成
	go func() {
		// 等待一小段时间确保service包完全初始化
		time.Sleep(100 * time.Millisecond)
		if mainCache := service.GetEnhancedTwoLevelCache(); mainCache != nil {
			instance.cacheWriteManager.SetMainCacheUpdater(func(key string, data []byte, ttl time.Duration) error {
				return mainCache.SetBothLevels(key, data, ttl)
			})
		}
	}()

	// 确保异步插件系统初始化
	plugin.InitAsyncPluginSystem()

	// 初始化插件管理器
	pluginManager := plugin.NewPluginManager()

	// 获取所有异步插件名称
	plugins := plugin.GetRegisteredPlugins()
	var pluginNames []string
	for _, p := range plugins {
		pluginNames = append(pluginNames, p.Name())
	}

	// 注册全局插件（根据配置过滤）
	if config.AppConfig.AsyncPluginEnabled {
		pluginManager.RegisterGlobalPluginsWithFilter(pluginNames)
	}

	// 更新默认并发数（如果插件被禁用则使用0）
	pluginCount := 0
	if config.AppConfig.AsyncPluginEnabled {
		pluginCount = len(pluginManager.GetPlugins())
	}
	config.UpdateDefaultConcurrency(pluginCount)

	// 初始化搜索服务
	instance.searchService = service.NewSearchService(pluginManager)

	return instance
}
