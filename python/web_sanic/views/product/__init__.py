from sanic import Blueprint
from .static import bp_static
from .product import bp_product

product = Blueprint.group(bp_static, bp_product, url_prefix="/product")