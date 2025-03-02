import torch
import onnx
from transformers import GPT2LMHeadModel, GPT2Tokenizer

# 加载 GPT-2
model_name = "/Users/liupeng/Downloads/gpts"
tokenizer = GPT2Tokenizer.from_pretrained(model_name)
model = GPT2LMHeadModel.from_pretrained(model_name)
model.config.use_cache = False  # 禁用 past_key_values

# 示例输入
text = "Hello, how are you?"
inputs = tokenizer(text, return_tensors="pt")

# 包装模型，避免传递 `use_cache`
class GPT2ONNXWrapper(torch.nn.Module):
    def __init__(self, model):
        super().__init__()
        self.model = model

    def forward(self, input_ids, attention_mask):
        return self.model(input_ids=input_ids, attention_mask=attention_mask).logits

onnx_model = GPT2ONNXWrapper(model)

# 导出 ONNX
onnx_path = "gpt2.onnx"
torch.onnx.export(
    onnx_model,
    (inputs["input_ids"], inputs["attention_mask"]),
    onnx_path,
    input_names=["input_ids", "attention_mask"],
    output_names=["logits"],
    dynamic_axes={
        "input_ids": {0: "batch_size", 1: "seq_length"},
        "attention_mask": {0: "batch_size", 1: "seq_length"},
        "logits": {0: "batch_size", 1: "seq_length", 2: "vocab_size"},
    }
)

print(f"✅ GPT-2 ONNX 模型导出完成: {onnx_path}")
