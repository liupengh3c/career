import asyncio
import aio_pika

class DwRabbitMQ:
    """DwRabbitMQ."""
    def __init__(self, url, queue_name, exchange_name, exchange_type='direct'):
         self.async_queue = asyncio.Queue()
         self.url = url
         self.exchange_name = exchange_name
         self.exchange_type = exchange_type
         self.queue_name = queue_name
         self.connection = None
         self.channel = None
         self.queue = None
         self.exchange = None

    async def connect(self):
        """connect、declare exchange and declare queue"""
        try:
            self.connection = await aio_pika.connect_robust(
                url=self.url
            )
            self.channel = await self.connection.channel()
            await self.channel.set_qos(prefetch_count=1000)
            self.exchange = await self.channel.declare_exchange(name=self.exchange_name, type=self.exchange_type, durable=True)
            # arguments = {
            #     "x-dead-letter-exchange": 'dlk:' + self.exchange_name,
            #     "x-dead-letter-routing-key": 'dlk:' + self.queue_name,
            # }
            # 队列声明为永久队列，持久化
            self.queue = await self.channel.declare_queue(name=self.queue_name, arguments=None, durable=True)
            # await self.channel.set_qos(prefetch_count=1000)
        except Exception as e:
            print("occur exception:",e)

    async def queue_bind_exchange(self):
        """queue_bind_exchange."""
        try:
            await self.queue.bind(exchange=self.exchange)
        except Exception as e:
            print("queue_bind_exchange error!")

    async def reconnect_and_bind(self):
        """reconnect and bind."""
        if self.connection.is_closed:
            await self.connect()
            await self.queue_bind_exchange()
            print('reconnect success')

    async def set_callback(self, callback):
        """set_callback."""
        try:
            await self.queue.consume(callback)
        except Exception as e:
            print("set_callback error!")

    async def publish_msg(self, content):
        """publish_msg."""
        message = aio_pika.Message(content.encode(),delivery_mode=aio_pika.DeliveryMode.PERSISTENT)
        try:
            # 设置为持久化消息
            await self.exchange.publish(message, routing_key=self.queue_name)
            print(" [x] Sent %r" % content)
        except Exception as e:
            self.reconnect_and_bind()
            await self.exchange.publish(message, routing_key=self.queue_name)
            print("reconnect......." ,e)

    async def callback(self, message: aio_pika.IncomingMessage):
        """callback."""
        try:
            # print(f"Consumer 1 received: {message.body.decode()}")
            await self.async_queue.put((message, message.body.decode()))
            # self.async_queue.put_nowait((message, message.body.decode()))
        except Exception:
            print("callback error!")


    def get_msg(self):
        """get_msg."""
        return self.async_queue.get_nowait()

    async def close(self):
        """close."""
        if not self.connection.is_closed:
            self.connection.close()

