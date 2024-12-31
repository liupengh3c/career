import torch
from PIL import Image
from transformers import CLIPProcessor, CLIPModel
import numpy as np
import os
from elasticsearch import Elasticsearch, helpers

# Elasticsearch服务器地址和端口
host = 'https://localhost:9200'
# 用户名和密码
username = 'elastic'
password = 'xpE4DQGWE9bCkoj7WXYE'
 
# 创建Elasticsearch客户端实例，并提供用户名和密码
es = Elasticsearch(hosts=[host], http_auth=(username, password), verify_certs=False,ca_certs="/Users/liupeng/Documents/study/elasticsearch-8.17.0/config/certs/http_ca.crt")
# 检查连接是否成功
if not es.ping():
    print("无法连接到Elasticsearch")
    exit()
else:
    print("成功连接到Elasticsearch")
# 现在你可以使用es变量来与Elasticsearch进行交互了

# 加载预训练的CLIP模型和处理器
model = CLIPModel.from_pretrained("/Users/liupeng/Documents/career/clip-vit-large-patch14")
processor = CLIPProcessor.from_pretrained("/Users/liupeng/Documents/career/clip-vit-large-patch14")

# 加载图像并进行预处理
# folder = "/Users/liupeng/Documents/career/cats_and_dogs_v2/train/cats"
folder = "/Users/liupeng/Documents/career/cats_and_dogs_v2/train/dogs"


for root, dirs, files in os.walk(folder):
    index_id = 1000
    for file in files:
        index_id += 1
        print(os.path.join(root, file))
        image = Image.open(os.path.join(root, file))  
        inputs = processor(images=image, return_tensors="pt")
        # 提取图像特征
        with torch.no_grad():
            image_features = model.get_image_features(**inputs)

        print("shape:",image_features.shape)

        # 对图像特征进行 L2 归一化
        # 使用 .norm() 计算 L2 范数并进行归一化
        image_features_normalized = image_features / image_features.norm(p=2, dim=-1, keepdim=True)
        numpy_array = image_features_normalized.numpy()
        # 打印归一化后的特征和特征的模长（应该为 1）
        # print("归一化后的图像特征:", numpy_array[0])
        # print("归一化后的模长:", image_features_normalized.norm(p=2, dim=-1))  # 应该接近 1
        documents = [
            {"name": "cat_"+str(index_id), "IFV": numpy_array[0].tolist(),"path":file},
        ]
        helpers.bulk(es, [
            {
                "_index": "vector_search_202412",
                "_id": index_id,
                "_source": doc
            }
            for doc in documents
        ])
        

