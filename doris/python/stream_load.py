from contextlib import redirect_stderr
from requests.auth import HTTPBasicAuth
import requests
import orjson
import json
import base64
from typing import List, Dict, Union
from datetime import datetime, date
import time
from _io import StringIO
import aiohttp
import asyncio

def streamload_v1(data_list: List[Dict[str, Union[str, int, float, None, datetime, date]]]):
    """有问题版本，没有解决重定向问题"""
    url = "http://xx.xx.xx.xx:8030/api/db_name/table_name/_stream_load"
    user = "xxxxx"
    password = "xxxxx"

    # 手动生成 Authorization header
    auth_str = f"{user}:{password}"
    auth_b64 = base64.b64encode(auth_str.encode("utf-8")).decode("utf-8")

    headers: dict[str, str] = {
        "label": f"streamload-id-{int(time.time())}",
        "strict_mode": "false",
        "strip_outer_array": "true",
        "Content-Type": "application/json",
        "format": "json",
        "Expect": "100-continue",
        "Authorization": f"Basic {auth_b64}",
    }

    # 转换 JSON
    body = json.dumps(data_list, default=str)

    # 超过 100MB 改为 JSON Lines
    if len(body) > 104857600:
        _ = headers.pop("strip_outer_array")
        headers["read_json_by_line"] = "true"

        buffer: StringIO = StringIO()
        for item in data_list:
            json_str = orjson.dumps(item).decode('utf-8', errors='ignore')
            _ = buffer.write(json_str + "\n")
        body = buffer.getvalue()

    session = requests.Session()
    while True:
        # allow_redirects=False 禁止自动跟随重定向
        resp = session.put(url, data=body, headers=headers, allow_redirects=False)
        # 根据http状态码判断，是否发生重定向
        if resp.status_code in (301, 302, 307, 308):
            # 获取重定向url，再次发送streamload指令
            url = resp.headers["Location"]
        # 返回结果
        if resp.status_code != 200:
            print(f"StreamLoad failed, status code: {resp.status_code}")
        # 打印doris返回结果
        print(resp.text)
        break

def streamload_v2(data_list: List[Dict[str, Union[str, int, float, None, datetime, date]]]):
    # Doris FE 节点（任意一个 FE，但最好是 Leader FE）
    url = "http://xx.xx.xx.xx:8030/api/db_name/table_name/_stream_load"
    user = "xxxxx"
    password = "xxxxx"

    # 手动生成 Authorization header
    auth_str = f"{user}:{password}"
    auth_b64 = base64.b64encode(auth_str.encode("utf-8")).decode("utf-8")

    headers: dict[str, str] = {
        "label": f"streamload-id-{int(time.time())}",
        "strict_mode": "false",
        "strip_outer_array": "true",
        "Content-Type": "application/json",
        "format": "json",
        "Expect": "100-continue",
        "Authorization": f"Basic {auth_b64}",
    }

    # 转换 JSON
    body = json.dumps(data_list, default=str)

    # 超过 100MB 改为 JSON Lines
    if len(body) > 104857600:
        _ = headers.pop("strip_outer_array")
        headers["read_json_by_line"] = "true"

        buffer: StringIO = StringIO()
        for item in data_list:
            json_str = orjson.dumps(item).decode('utf-8', errors='ignore')
            _ = buffer.write(json_str + "\n")
        body = buffer.getvalue()

    session = requests.Session()
    while True:
        # 不自动跟随重定向，手动处理
        resp = session.put(url, data=body, headers=headers, allow_redirects=False)
        if resp.status_code in (301, 302, 307, 308):
            redirect_url = resp.headers.get("Location")
            if not redirect_url:
                raise RuntimeError("Redirect response without Location header")
            print(f"Redirecting to: {redirect_url}")
            url = redirect_url
            # loop 再次 PUT，保持 headers 不变
            continue
        else:
            # 返回结果
            if resp.status_code != 200:
                print(f"StreamLoad failed, status code: {resp.status_code}")
            print(resp.text)
            break

async def streamload_aio(
    data_list: List[Dict[str, Union[str, int, float, None, datetime, date]]],
    url: str,
    user: str,
    password: str,
):
    # aiohttp 自动处理 Authorization
    auth = aiohttp.BasicAuth(user, password)

    # 通用 headers
    headers = {
        "label": f"streamload-id-{int(asyncio.get_event_loop().time()*1000)}",
        "strict_mode": "false",
        "format": "json",
        "Expect": "100-continue",
        "Content-Type": "application/json",
        "strip_outer_array": "true",
    }

    # 转 JSON
    body = json.dumps(data_list, default=str)

    # 超过 100MB 转 JSON Lines
    if len(body.encode('utf-8')) > 104857600:
        _ = headers.pop("strip_outer_array", None)
        headers["read_json_by_line"] = "true"
        buffer: StringIO = StringIO()
        for item in data_list:
            json_str = orjson.dumps(item).decode('utf-8', errors='ignore')
            _ = buffer.write(json_str + "\n")
        body = buffer.getvalue()

    async with aiohttp.ClientSession(auth=auth, headers=headers) as session:
        async with session.put(url, data=body) as resp:
            text = await resp.text()
            if resp.status != 200:
                print(f"StreamLoad failed, status code: {resp.status}")
            print(text)