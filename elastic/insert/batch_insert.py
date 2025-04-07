from elasticsearch import AsyncElasticsearch
from elasticsearch.helpers import async_bulk
import asyncio

async def main():
    # 创建异步 Elasticsearch 客户端
    es = AsyncElasticsearch(
        hosts=["https://localhost:9200"],  # 集群地址
        # 可选参数示例:
        http_auth=("elastic", "jh_qCnIHjYz_sMp-7uT-"),  # 认证信息
        verify_certs=False,  # 关闭SSL验证
        ca_certs = '/Users/liupeng/Documents/elastic/elasticsearch-8.17.0/config/certs/http_ca.crt'
        # timeout=30  # 请求超时时间
    )

    # 生成测试数据（替换为实际数据）
    async def generate_actions():
        for i in range(1000):
            yield {
                "_index": "test-index",
                "_op_type": "index",  # 操作类型：index/create/update/delete
                "_id": f"doc_{i}",     # 可选文档ID
                "_source": {
                    "title": f"Document {i}",
                    "content": "Sample content",
                    "value": i * 10
                }
            }

    try:
        # 执行批量写入（推荐参数配置）
        success_count, errors = await async_bulk(
            es,
            generate_actions(),
            chunk_size=500,        # 每批处理文档数（根据集群性能调整）
            max_retries=3,        # 最大重试次数
            initial_backoff=1,     # 初始重试等待时间(秒)
            request_timeout=60     # 单次请求超时时间
        )
        
        print(f"成功写入 {success_count} 条文档")
        if errors:
            print(f"发现 {len(errors)} 条错误:")
            for error in errors:
                print(f"文档 {error['index']['_id']} 错误: {error['index']['error']}")

    except Exception as e:
        print(f"发生批量写入异常: {str(e)}")
    finally:
        await es.close()

if __name__ == "__main__":
    asyncio.run(main())