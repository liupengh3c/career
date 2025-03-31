#!/usr/bin/env python
import asyncio
import aio_pika
import time

class DwRabbitMQ:
    def __init__(self, url, exchange, queue_name):
         self.async_queue = asyncio.Queue()
         self.url = url
         self.exchange = exchange
         self.queue_name = queue_name
         self.connection = None
         self.channel = None
         self.queue = None
    async def connect(self):
        try:
            self.connection = await aio_pika.connect_robust(
                self.url
            )
            self.channel = await self.connection.channel()
            self.exchange = await self.channel.declare_exchange(name=self.exchange, type='fanout')
            self.queue = await self.channel.declare_queue(name=self.queue_name, durable=True)
        except Exception as e:
            print(e)
    async def queue_bind_exchange(self):
        try:
            await self.queue.bind(exchange=self.exchange)
        except Exception as e:
            print(e)
    async def set_callback(self, callback):
        try:
            await self.queue.consume(callback)
        except Exception as e:
            print(e)
    async def callback(self, message: aio_pika.IncomingMessage):
        try:
            # print(f"Consumer 1 received: {message.body.decode()}")
            self.async_queue.put_nowait((message,message.body.decode()))
        except Exception:
            print(f" [!] Rejected: {message.body.decode()}")

    async def is_connected(self):
        return self.connection.is_closed
    def get_msg(self):
        return self.async_queue.get_nowait()

async def main():
    buffer = []
    batch_size = 10
    last_time = time.time()
    url = "amqp://guest:guest@192.168.3.15:5672/"
    dwmq = DwRabbitMQ(url, "exchange_fan", "dw_clip_insert")
    await dwmq.connect()
    await dwmq.queue_bind_exchange()
    await dwmq.set_callback(dwmq.callback)
    print(" [*] Waiting for messages. To exit press CTRL+C")
    while True:
        if await dwmq.is_connected():
            try:
                print("reconnecting...")
                dwmq = DwRabbitMQ(url, "exchange_fan", "dw_clip_insert")
                await dwmq.connect()
                await dwmq.set_callback(dwmq.callback)
            except Exception as e:
                print(e,",retry to connect")
                await asyncio.sleep(1)
                continue
        try:
            msg = dwmq.get_msg()
        except asyncio.QueueEmpty:
            await asyncio.sleep(1)
            msg = None
        if msg is not None:
            buffer.append(msg)
        if len(buffer) < batch_size and time.time() - last_time < 60:
            continue
        last_time = time.time()
        print("batch size : {},real size : {}, begin to process".format(batch_size,len(buffer)))
        # 根据实际需求，批量
        for msg in buffer:
            try:
                await msg[0].ack()
            except Exception as e:
                print(e)
            print(msg[1])
        buffer = []

if __name__ == '__main__':
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print('Interrupted')