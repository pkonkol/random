import asyncio
import zmq
import zmq.asyncio
import time

from collections import deque
proxies = deque(["1", "2", "3", "4"])

print("Started master")

context = zmq.asyncio.Context()

async def push_tasks():
    task_sock = context.socket(zmq.PUSH)
    task_sock.bind("tcp://*:5556")
    while True:
        print("In push_tasks")
        await task_sock.send_string(f"I'm alive {time.ctime()}")
        await asyncio.sleep(13)

async def recv_result():
    result_sock = context.socket(zmq.PULL)
    result_sock.bind("tcp://*:5558")
    while True:
        print("In recv_result")
        msg = await result_sock.recv_string()
        print(f"Recived response: {msg}")
        await asyncio.sleep(1)

async def return_proxy():
    sock = context.socket(zmq.REP)
    sock.bind("tcp://*:5557")
    while True:
        print("In return_proxy")
        msg = await sock.recv_string()
        print("Received proxy request")
        await sock.send_string(f"proxy return from master: {proxies[0]}")
        proxies.rotate(-1)


async def parse_categories():
    sitemap_url = "olx.pl/sitemap"
    while True:
        print("Main loop parsing")
        await asyncio.sleep(10)

def create_tasks():
    tasks = []
    tasks.append(asyncio.create_task(push_tasks()))
    tasks.append(asyncio.create_task(recv_result()))
    tasks.append(asyncio.create_task(return_proxy()))
    tasks.append(asyncio.create_task(parse_categories()))
    return tasks

async def main():
    await asyncio.gather(*create_tasks())


if __name__ == "__main__":
    # loop = asyncio.get_event_loop()
    print('in __main__')
    asyncio.run(main())

