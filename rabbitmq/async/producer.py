#!/usr/bin/env python
import aio_pika
import time
import asyncio
from datetime import datetime

# connection = pika.BlockingConnection(
#     pika.ConnectionParameters(host='localhost',port=5672))
async def main() -> None:
    print('Connected to RabbitMQ')
    connection = await aio_pika.connect_robust(
        "amqp://guest:guest@192.168.3.15:5672/"
    )
    channel = await connection.channel()
    exchange = await channel.declare_exchange(name='exchange_fan', type='fanout')
    # durable=True 队列持久化
    q1 = await channel.declare_queue(name='dw_clip_insert', durable=True)
    await q1.bind(exchange='exchange_fan')

    q2 = await channel.declare_queue(name='dw_tag_insert', durable=True)
    await q2.bind(exchange='exchange_fan')
    while True:
        current_time_str = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        # 消息持久化+队列持久化，才能保证消息在mq重启后仍然存在
        message = aio_pika.Message(current_time_str.encode(),delivery_mode=aio_pika.DeliveryMode.PERSISTENT)
        await exchange.publish(message, routing_key='')
        print("Message published to fanout exchange!")
        time.sleep(1)
    print(" [x] Sent 'Hello World!'")
    connection.close()

if __name__ == '__main__':
    a = {"a": 1}
    clip_obj = {"name": '{"action": "login", "timestamp": 1622505600, "ip": "192.168.1.1"}', "age": 30, "city": "New York"}
    clip_obj.update(a)
    new_clip = clip_obj.items()
    
    print(clip_obj)
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print('Interrupted')