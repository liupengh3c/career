from sanic import Blueprint
from .static import bp_static
from .user import bp_user

user = Blueprint.group(bp_static, bp_user, url_prefix="/user")