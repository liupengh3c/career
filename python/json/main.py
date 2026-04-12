import json


def flatten_json_func(data, parent_key=""):
    """flatten json (非递归高性能版本)"""

    def detect_type(value):
        """检测单个值的类型"""
        if isinstance(value, bool):
            return "tinyint"
        elif isinstance(value, int):
            return "double"
        elif isinstance(value, float):
            return "double"
        elif isinstance(value, str):
            if value in ["NaN", "nan", "inf", "-inf"]:
                return "double"
            else:
                return "varchar(65533)"
        elif isinstance(value, dict):
            return "text"
        elif isinstance(value, list):
            return "list"
        elif value is None:
            return "null"
        else:
            return "jsonb"

    def detect_array_type(arr):
        """检测数组内部元素类型"""
        if not arr:
            return "array<unknown>"  # 空数组情况
        element_types = {detect_type(el) for el in arr}
        if len(element_types) == 1:
            return f"array<{element_types.pop()}>"
        else:
            return f"array<mixed:{','.join(sorted(element_types))}>"

    columns = {}
    body = {}
    parquet_body = {}

    # 显式栈，避免递归 & 频繁 dict 合并
    stack = [(data, parent_key)]

    while stack:
        current_value, current_key = stack.pop()

        if isinstance(current_value, dict):
            for k, v in current_value.items():
                new_key = f"{current_key}.{k}" if current_key else k
                stack.append((v, new_key))

        elif isinstance(current_value, list):
            if len(current_value) > 0:
                true_parent_key = current_key.replace(".", "__")
                parent_key_size = current_key + "_alen"
                true_parent_key_size = f"{true_parent_key}_alen"

                columns[current_key] = detect_array_type(current_value)

                if columns[current_key] in ["array<text>", "array<list>", "array<jsonb>"]:
                    # 复杂数组整体序列化到 parquet_body
                    parquet_body[true_parent_key] = json.dumps(current_value)
                else:
                    columns[parent_key_size] = "bigint"
                    body[true_parent_key_size] = len(current_value)
                    body[true_parent_key] = current_value

        else:
            true_parent_key = current_key.replace(".", "__")
            columns[current_key] = detect_type(current_value)

            v = current_value
            if isinstance(v, bool):
                v = 1 if v else 0

            if isinstance(v, str) and v in ["NaN", "nan"]:
                body[true_parent_key] = None
            else:
                body[true_parent_key] = v

    return columns, body, parquet_body
    
if __name__ == "__main__":
    data = {
        "name": "John",
        "age": 30,
        "address": {
            "city": "New York",
            "state": "NY"
        },
        "items": [
            {
                "name": "Item 1",
                "price": 100
            },
            {
                "name": "Item 2",
                "price": 200
            }
        ]
    }
    columns, body, parquet_body = flatten_json_func(data)
    print(columns)
    print(body)
    print(parquet_body)