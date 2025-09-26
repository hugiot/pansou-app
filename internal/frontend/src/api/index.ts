import axios from 'axios';
import type { SearchResponse } from '@/types';
import { Health, Search } from "../../wailsjs/go/main/App";
import {pansou, model} from "../../wailsjs/go/models.ts";

const api = axios.create({
    baseURL: '/api',
    timeout: 10000
});

// 搜索参数接口
export interface SearchParams {
    kw: string;
    refresh?: boolean;
    res?: 'all' | 'results' | 'merge';
    src?: 'all' | 'tg' | 'plugin';
    plugins?: string;
    ext?: string;
}

// API响应包装类型
interface ApiResponse<T> {
    code: number;
    message: string;
    data: T;
}

// 健康状态接口（基于实际API返回）
export interface HealthStatus {
    status: string;
    plugins_enabled: boolean;
    plugin_count: number;
    plugins: string[];
    channels: string[];
}

// 获取API健康状态
export const getHealth = async (): Promise<HealthStatus> => {
    try {
        const healthParams = pansou.HealthRequest.createFrom({})
        const response = await Health(healthParams)
        return response as HealthStatus
    } catch (error) {
        console.error('获取健康状态失败:', error);
        // 返回模拟数据
        return getMockHealthData();
    }
};

// 模拟健康状态数据
const getMockHealthData = (): HealthStatus => {
    return {
        status: "ok",
        plugins_enabled: true,
        plugin_count: 6,
        plugins: ["pansearch", "hdr4k", "shandian", "muou", "duoduo", "labi"],
        channels: ["tgsearchers3", "SharePanBaidu", "yunpanxunlei", "tianyifc", "BaiduCloudDisk"]
    };
};

// 搜索API
export const search = async (params: SearchParams): Promise<SearchResponse> => {
    const searchParams = model.SearchRequest.createFrom({
        ...params,
        ext: { referer: "https://dm.xueximeng.com" }
    })

    // console.log('搜索参数:', searchParams);
    try {
        const response = await Search(searchParams);
        const results = response as unknown as SearchResponse;
        console.log(results)
        return results
    } catch (error) {
        console.error('API错误:', error);
        return getMockData();
    }
};

// 模拟数据（开发阶段使用）
const getMockData = (): SearchResponse => {
    return {
        total: 15,
        results: [
            {
                message_id: "12345",
                unique_id: "channel-12345",
                channel: "tgsearchers3",
                datetime: "2023-06-10T14:23:45Z",
                title: "速度与激情全集1-10",
                content: "速度与激情系列全集，1080P高清...",
                links: [
                    {
                        type: "baidu",
                        url: "https://pan.baidu.com/s/1abcdef",
                        password: "1234"
                    }
                ],
                tags: ["电影", "合集"]
            }
        ],
        merged_by_type: {
            baidu: [
                {
                    url: "https://pan.baidu.com/s/1abcdef",
                    password: "1234",
                    note: "速度与激情全集1-10",
                    datetime: "2023-06-10T14:23:45Z",
                    source: "tgsearchers3"
                },
                {
                    url: "https://pan.baidu.com/s/1ghijkl",
                    password: "5678",
                    note: "速度与激情9",
                    datetime: "2023-05-15T10:20:30Z",
                    source: "SharePanBaidu"
                }
            ],
            aliyun: [
                {
                    url: "https://www.aliyundrive.com/s/abcdef",
                    note: "速度与激情系列合集",
                    datetime: "2023-07-01T08:15:20Z",
                    source: "yunpanxunlei"
                }
            ],
            "115": [
                {
                    url: "https://115.com/s/abcdefg",
                    password: "abc123",
                    note: "速度与激情1-10全集高清资源",
                    datetime: "2023-04-22T16:45:12Z",
                    source: "pansearch插件"
                }
            ]
        }
    };
};

export default api;