from sanic import Sanic
# from sanic.response import text, json
from sanic import response
from tortoise.contrib.sanic import register_tortoise
from models import user
# from models.user import User
# import models

app = Sanic("MySanicApp")
register_tortoise(
    app, db_url="mysql://sanic:sanic123@60.205.177.189:3306/career", modules={"models": ["models.user"]}, generate_schemas=True
)

@app.post("/user_list")
async def user_list(request):
    resp_users = []

    users = await user.User.all()
    for u in users:
        one = {
            'id': u.id,
            'name': u.name,
            'age': u.age,
            'sex': u.sex,
            'grade': u.grade,
        }
        resp_users.append(one)
    return response.json(resp_users)

@app.post("/create_user")
async def create_user(request):
    # print(request.json)
    one = {
        'name': request.json['name'],
        'age': request.json['age'],
        'sex': request.json['sex'],
        'grade': request.json['grade'],
    }
    await user.User.create(**one)
    return response.json({'status': True})

if __name__ == "__main__":
    app.run()
