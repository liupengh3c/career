import torch
from PIL import Image
from transformers import CLIPProcessor, CLIPModel
import numpy as np
import os
from elasticsearch import Elasticsearch, helpers 

# 加载预训练的CLIP模型和处理器
model = CLIPModel.from_pretrained("/Users/liupeng/Documents/career/clip-vit-large-patch14")
processor = CLIPProcessor.from_pretrained("/Users/liupeng/Documents/career/clip-vit-large-patch14")

# 加载图像并进行预处理
# folder = "/Users/liupeng/Documents/career/cats_and_dogs_v2/train/cats"
image_file = "/Users/liupeng/Documents/career/cats_and_dogs_v2/train/dogs/dog.11002.jpg"
image = Image.open(image_file)
input_data = np.array(image)
print("input shape:",input_data.shape)
# inputs = processor(images=image, return_tensors="pt")
inputs = processor(images=input_data, return_tensors="pt")
# 提取图像特征
with torch.no_grad():
    image_features = model.get_image_features(**inputs)

print("shape:",image_features.shape)

# 对图像特征进行 L2 归一化
# 使用 .norm() 计算 L2 范数并进行归一化
image_features_normalized = image_features / image_features.norm(p=2, dim=-1, keepdim=True)
numpy_array = image_features_normalized.numpy()
# 打印归一化后的特征和特征的模长（应该为 1）
print("归一化后的图像特征:", numpy_array[0])
# print("归一化后的模长:", image_features_normalized.norm(p=2, dim=-1))  # 应该接近 1

        

