from sanic import Blueprint
from sanic import response
bp_static = Blueprint('order_static', url_prefix='/static')

@bp_static.route('/info',methods=['POST'])
async def info(request):
    return response.json({'status': True})