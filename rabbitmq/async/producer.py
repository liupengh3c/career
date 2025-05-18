from datetime import datetime
import asyncio
import time
from rabbitmq_utils import DwRabbitMQ
import json


async def main():
    rabbitmq = DwRabbitMQ("amqp://guest:guest@127.0.0.1:5672/", 'dw_queue', 'dw_exchange')
    await rabbitmq.connect()
    await rabbitmq.queue_bind_exchange()
   
    while True:
        # current_time_str = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        current_time_str = json.dumps({
            "message": f"Hello RabbitMQ at {datetime.now()}",
            "timestamp": int(time.time()),
        })
        await rabbitmq.publish_msg(current_time_str)
        # time.sleep(1)
        pass
if __name__ == '__main__':
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print('Interrupted')