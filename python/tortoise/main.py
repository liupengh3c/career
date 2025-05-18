import asyncio
import datetime
from tortoise import Tortoise, fields, run_async,Model

class BaseModel(Model):
    _hidden_fields = ["id", "sex"]
    _field_map = {}
    _filed_order = []

    def to_dict(self):
        result = {}
        fields_order = self._meta.fields or self._field_order
        for field in fields_order:
            if field in self._hidden_fields:
                continue
            value = getattr(self, field)
            if isinstance(value, datetime.datetime):
                value = value.strftime("%Y-%m-%d %H:%M:%S")
            s_key = self._field_map.get(field, field)
            result[s_key] = value
        return result
    
    class Meta:
        abstract = True

# 定义user模型
class User(BaseModel):
    id = fields.IntField(pk=True)
    name = fields.CharField(max_length=50)
    age = fields.IntField()
    sex = fields.CharField(max_length=10)
    grade = fields.IntField()

    class Meta:
        table = "users"  # 数据库中表名

    def __str__(self):
        return f"User(id={self.id}, name='{self.name}', age={self.age})"


# 初始化数据库连接
async def init():
    await Tortoise.init(
        db_url="mysql://sanic:sanic123@xx.xx.xx.189:3306/career",
        modules={"models": ["__main__"]},  # 当前模块下定义模型
    )
    # await Tortoise.generate_schemas()


# 增删改查操作示例
async def run():
    await init()

    # 🔹 新增
    user = await User.create(name="Alice", age=30, sex = '女',grade= 4)
    # print("Inserted:", "name:",user.name, "age:",user.age, "sex:",user.sex, "grade:",user.grade)

    # 🔹 查询
    # user_fetched = await User.get(id=user.id)
    user_fetched = await User.filter(id=user.id).first()
    print("to_dict:", user_fetched.to_dict())
    # print("Fetched:", "name:",user.name, "age:",user.age, "sex:",user.sex, "grade:",user.grade)

    # 🔹 更新
    user_fetched = await User.filter(id=10).first()
    user_fetched.age = 100
    await user_fetched.save()
    new_user = await User.filter(id=10).first()
    # print("Updated:", "name:",new_user.name, "age:",new_user.age, "sex:",new_user.sex, "grade:",new_user.grade)

    cnt = await User.filter(id=10).update(age=3100)
    print("cnt:", cnt)
    new_doc = await User.filter(id=10).first()
    # print("Updated 2:", "name:",new_doc.name, "age:",new_doc.age, "sex:",new_doc.sex, "grade:",new_doc.grade)
    # 🔹 删除
    dele = await User.filter(id=9).delete()
    # print("Deleted cnt", dele)

    rec = await User.filter(id=9).first()
    if rec is not None:
        print("Not deleted")
    else:
        print("Deleted successfully")
    await Tortoise.close_connections()


if __name__ == "__main__":
    run_async(run())
