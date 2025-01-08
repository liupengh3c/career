import torch
from PIL import Image
from transformers import CLIPProcessor, CLIPModel
import numpy as np
import os
import sys


def main(folder):
    # 加载预训练的CLIP模型和处理器
    model = CLIPModel.from_pretrained("/Users/liupeng/Documents/career/clip-vit-large-patch14")
    processor = CLIPProcessor.from_pretrained("/Users/liupeng/Documents/career/clip-vit-large-patch14")

    # 加载图像并进行预处理
    # folder = "/Users/liupeng/Documents/career/cats_and_dogs_v2/train/cats"
    # folder = "/Users/liupeng/Documents/career/cats_and_dogs_v2/train/dogs"

    for root, dirs, files in os.walk(folder):
        index_id = 1000
        for file in files:
            index_id += 1
            # print(os.path.join(root, file))
            image = Image.open(os.path.join(root, file))  
            inputs = processor(images=image, return_tensors="pt")
            # 提取图像特征
            with torch.no_grad():
                image_features = model.get_image_features(**inputs)

            # print("shape:",image_features.shape)

            # 对图像特征进行 L2 归一化
            # 使用 .norm() 计算 L2 范数并进行归一化
            image_features_normalized = image_features / image_features.norm(p=2, dim=-1, keepdim=True)
            numpy_array = image_features_normalized.numpy()
            numpy_to_list = numpy_array.tolist()
            # 打印归一化后的特征和特征的模长（应该为 1）
            print(numpy_to_list)
            # print("归一化后的图像特征:", numpy_to_list)
            # print("归一化后的模长:", image_features_normalized.norm(p=2, dim=-1))  # 应该接近 1
        
if __name__ == "__main__":
    args = sys.argv[1:]  # 获取所有传递给脚本的参数
    image_path = args[0]
    main(image_path)
