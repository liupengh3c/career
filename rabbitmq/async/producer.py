from datetime import datetime
import asyncio
import time
from rabbitmq_utils import DwRabbitMQ


async def main():
    rabbitmq = DwRabbitMQ("amqp://guest:guest@192.168.3.15:5672/addw", 'dw_jiesuan_queue', 'dw_jiesuan_exchange')
    await rabbitmq.connect()
    await rabbitmq.queue_bind_exchange()
    while True:
        current_time_str = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        await rabbitmq.publish_msg(current_time_str)
        time.sleep(1)
        pass
if __name__ == '__main__':
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print('Interrupted')