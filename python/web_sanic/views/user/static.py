from sanic import Blueprint
from sanic import response
from sanic.request import Request
import time
bp_static = Blueprint('user_static', url_prefix='/static')

@bp_static.middleware
async def print_on_request(request:Request):
    # print("I am a spy")
    request.ctx.start_time = time.time()

@bp_static.middleware("response")
async def halt_response(request:Request, response):
    # 计算总耗时（单位：秒）
    duration = time.time() - request.ctx.start_time
    
    # 获取API路径和HTTP方法
    api_path = request.url.split('?')[0]  # 去除查询参数
    method = request.method
    
    # 输出日志（可根据需求替换为监控上报）
    print(f"[{method}] {api_path} - 耗时: {duration:.4f}s")


@bp_static.route('/info',methods=['POST'])
async def info(request):
    time.sleep(3)
    print("I am the user static route handler")
    return response.json({'status': True})