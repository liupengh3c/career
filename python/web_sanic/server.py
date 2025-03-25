from sanic import Sanic
# from sanic.response import text, json
from sanic import response
from tortoise.contrib.sanic import register_tortoise
from models import user
from views import bp_views

app = Sanic("MySanicApp")
register_tortoise(
    app, db_url="mysql://sanic:sanic123@60.205.177.189:3306/career", modules={"models": ["models.user"]}, generate_schemas=True
)
# 注册蓝图
app.blueprint(bp_views)
if __name__ == "__main__":
    app.run()
