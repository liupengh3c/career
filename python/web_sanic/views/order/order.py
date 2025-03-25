from sanic import Blueprint
from sanic import response

bp_order = Blueprint("order", url_prefix="/order")

@bp_order.route("/add", methods=["POST"])
async def add(request):
    return response.json({})

@bp_order.route("/list", methods=["POST"])
async def list(request):
    return response.json({})