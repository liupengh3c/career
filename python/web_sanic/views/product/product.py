from sanic import Blueprint
from sanic import response

bp_product = Blueprint("product", url_prefix="/product")

@bp_product.route("/add", methods=["POST"])
async def add(request):
    return response.json({})

@bp_product.route("/list", methods=["POST"])
async def list(request):
    return response.json({})