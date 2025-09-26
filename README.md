# Pansou App

使用 **Wails** 开发，跨平台的 **Pansou** 桌面应用

## 快速开始

1. 通过 Releases 下载最新版本
2. 安装

## 开发者


### 构建

```bash
# build all
wails build -tags webkit2_41 -platform darwin,windows,linux -upx -nsis

# build for windows/amd64
wails build -tags webkit2_41 -platform windows/amd64 -upx -nsis
```

## 更新记录

* 2025-09-26
  * 发布 v0.1.0
  * pansou 版本：main@789cba8（新增插件bixin）
  * pansou-web 版本：main@b708ab9（update）

## 感谢

* [pansou](https://github.com/fish2018/pansou)
* [pansou-web](https://github.com/fish2018/pansou-web)
* [Wails](https://wails.io/)


