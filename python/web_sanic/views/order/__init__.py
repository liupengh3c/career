from sanic import Blueprint
from .static import bp_static
from .order import bp_order

order = Blueprint.group(bp_static, bp_order, url_prefix="/order")