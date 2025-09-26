package main

import (
	"context"
	"pansou/config"
	"pansou/model"

	"github.com/hugiot/pansou-app/internal/pansou"
)

// App struct
type App struct {
	ctx     context.Context
	service pansou.Pansou
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.service = pansou.New()
}

func (a *App) shutdown(_ context.Context) {
	_ = a.service.Release()
}

func (a *App) Health(req pansou.HealthRequest) pansou.HealthResponse {
	return a.service.Health(req)
}

func (a *App) Search(req model.SearchRequest) (pansou.SearchResponse, error) {
	// 检查并设置默认值
	if len(req.Channels) == 0 {
		req.Channels = config.AppConfig.DefaultChannels
	}

	// 如果未指定结果类型，默认返回merge并转换为merged_by_type
	if req.ResultType == "" {
		req.ResultType = "merged_by_type"
	} else if req.ResultType == "merge" {
		// 将merge转换为merged_by_type，以兼容内部处理
		req.ResultType = "merged_by_type"
	}

	// 如果未指定数据来源类型，默认为全部
	if req.SourceType == "" {
		req.SourceType = "all"
	}

	// 参数互斥逻辑：当src=tg时忽略plugins参数，当src=plugin时忽略channels参数
	if req.SourceType == "tg" {
		req.Plugins = nil // 忽略plugins参数
	} else if req.SourceType == "plugin" {
		req.Channels = nil // 忽略channels参数
	} else if req.SourceType == "all" {
		// 对于all类型，如果plugins为空或不存在，统一设为nil
		if req.Plugins == nil || len(req.Plugins) == 0 {
			req.Plugins = nil
		}
	}

	response, err := a.service.Search(req)
	if err != nil {
		return pansou.SearchResponse{}, err
	}

	if len(response.Results) == 0 {
		response.Results = make([]pansou.SearchResult, 0)
	}

	if len(response.MergedByType) == 0 {
		response.MergedByType = make(pansou.MergedLinks)
	}

	return response, nil
}
