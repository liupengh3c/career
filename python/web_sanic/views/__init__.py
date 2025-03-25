from sanic import Blueprint
from .order import order
from .product import product
from .user import user
from sanic import Sanic, response
from sanic.request import Request
import time

bp_views = Blueprint.group(order, product, user,url_prefix="/api")
