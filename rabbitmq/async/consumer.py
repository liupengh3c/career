from datetime import datetime
import time
import asyncio
from rabbitmq_utils import DwRabbitMQ


async def main():
    batch_size = 60
    max_delay = 60
    lists = []
    dw_rabbitmq = DwRabbitMQ("amqp://guest:guest@192.168.3.15:5672/addw", 'dw_jiesuan_queue', 'dw_jiesuan_exchange')
    await dw_rabbitmq.connect()
    await dw_rabbitmq.queue_bind_exchange()
    await dw_rabbitmq.set_callback(dw_rabbitmq.callback)
    last_time = time.time()
    while True:
        try:
            msg = dw_rabbitmq.get_msg()
        except asyncio.QueueEmpty:
            await asyncio.sleep(1)
            msg = None
        if msg:
            lists.append(msg)
        if len(lists) < batch_size and time.time() - last_time < max_delay:
            continue
        # 批量处理消息
        for i in range(len(lists)):
            print("[x] Received %r" % lists[i][1])
            try:
                await lists[i][0].ack()
            except Exception as e:
                print(e)
        last_time = time.time()
        lists = []
            
if __name__ == '__main__':
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print('Interrupted')