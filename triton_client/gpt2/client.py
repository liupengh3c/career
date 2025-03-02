import tritonclient.http as httpclient
import numpy as np
from transformers import GPT2Tokenizer

# 连接 Triton Server
triton_client = httpclient.InferenceServerClient(url="192.168.3.40:8000")

# 加载分词器
tokenizer = GPT2Tokenizer.from_pretrained("/models/triton_client/gpt2/gpts")

# 交互式对话
conversation_history = []

while True:
    user_input = input("You: ")
    if user_input.lower() in ["exit", "quit"]:
        print("Chat ended.")
        break

    # 记录对话历史
    conversation_history.append(user_input)
    context = " ".join(conversation_history)

    # Tokenize 输入
    tokens = tokenizer(context, return_tensors="np")
    input_ids = tokens["input_ids"].astype(np.int64)
    attention_mask = tokens["attention_mask"].astype(np.int64)

    # 创建请求
    inputs = [
        httpclient.InferInput("input_ids", input_ids.shape, "INT64"),
        httpclient.InferInput("attention_mask", attention_mask.shape, "INT64"),
    ]
    inputs[0].set_data_from_numpy(input_ids)
    inputs[1].set_data_from_numpy(attention_mask)

    outputs = [httpclient.InferRequestedOutput("logits")]

    # 发送请求
    response = triton_client.infer(model_name="gpt2", inputs=inputs, outputs=outputs)

    # 解析结果
    logits = response.as_numpy("logits")
    next_token_id = np.argmax(logits[0, -1, :])
    reply = tokenizer.decode([next_token_id])
    print(conversation_history)

    print(f"GPT-2: ",reply)

    # 记录 AI 回复
    conversation_history.append(reply)
