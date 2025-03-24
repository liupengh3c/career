# from tortoise.models import Model
# from tortoise import fields
from tortoise import Model, fields
class User(Model):
    id = fields.IntField(pk=True)
    name = fields.CharField(max_length=255)
    age = fields.IntField()
    sex = fields.TextField()
    grade = fields.CharField(max_length=255)
    class Meta:
        table = "user"