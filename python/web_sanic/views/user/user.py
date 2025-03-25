from sanic import Blueprint
from models import user
from sanic import response

bp_user = Blueprint("user", url_prefix="/user")

@bp_user.route("/add",methods=["POST"])
async def add(request):
    one = {
        'name': request.json['name'],
        'age': request.json['age'],
        'sex': request.json['sex'],
        'grade': request.json['grade'],
    }
    await user.User.create(**one)
    return response.json({'status': True})

@bp_user.route("/list", methods=["POST"])
async def list(request):
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