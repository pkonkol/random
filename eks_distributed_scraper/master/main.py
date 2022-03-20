import asyncio
import time

from aio_pika import Message, connect
from collections import deque
proxies = deque(range(0, 10))

print("Started master")

async def push_tasks(connection):
    async with connection:
        tasks_channel = await connection.channel()
        tasks_todo_q = await tasks_channel.declare_queue('tasks_todo')
        while True:
            print("In push_tasks")
            await tasks_channel.default_exchange.publish(
                Message(bytes(f"Task at {time.ctime()}", 'ascii')),
                routing_key=tasks_todo_q.name
            )
            await asyncio.sleep(5)
        

async def recv_result(connection):
    async with connection:
        tasks_channel = await connection.channel()
        tasks_done_q = await tasks_channel.declare_queue('tasks_done')
        while True:
            print("In recv_result")
            if resp := await tasks_done_q.get(fail=False):
                await resp.ack()
                print(f"got response {resp.body}")
            await asyncio.sleep(1)


async def return_proxy(connection):
    async with connection:
        proxy_channel = await connection.channel()
        proxy_request_q = await proxy_channel.declare_queue('proxy_request')
        proxy_response_q = await proxy_channel.declare_queue('proxy_response')
        while True:
            print("In return_proxy")
            if resp := await proxy_request_q.get(fail=False):
                await resp.ack()
                print(f"got proxy request {resp.body}")
                proxy = proxies[0]
                proxies.rotate(-1)
                await proxy_channel.default_exchange.publish(
                    Message(bytes(str(proxy), 'ascii')),
                    routing_key=proxy_response_q.name
                )
            await asyncio.sleep(1)


async def parse_categories():
    sitemap_url = "olx.pl/sitemap"
    while True:
        print("Main loop parsing")
        await asyncio.sleep(10)

def create_tasks(c):
    tasks = []
    tasks.append(asyncio.create_task(push_tasks(c)))
    tasks.append(asyncio.create_task(recv_result(c)))
    tasks.append(asyncio.create_task(return_proxy(c)))
    # tasks.append(asyncio.create_task(parse_categories()))
    return tasks

async def main():
    connection = await connect('amqp://localhost')
    await asyncio.gather(*create_tasks(connection))

if __name__ == "__main__":
    print('in __main__')
    asyncio.run(main())
