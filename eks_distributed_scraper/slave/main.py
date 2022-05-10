import random
import os
import time
import pika

print("Started slave")

HOSTNAME = os.environ.get("HOSTNAME")
PROXY = None

connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
channel = connection.channel()
channel.queue_declare(queue='tasks_todo', passive=True)
channel.queue_declare(queue='tasks_done', passive=True)
proxy_channel = connection.channel()
proxy_channel.queue_declare(queue='proxy-request')


def get_proxy():
    print("sending proxy request")
    proxy_channel.basic_publish(exchange='',
                                routing_key='proxy_request',
                                body=f"proxy_request from {HOSTNAME}")
    print("sent proxy request")

    method, header, body = channel.basic_get('proxy_response')
    if method:
        channel.basic_ack(method.delivery_tag)
    else:
        print('didnt retrieve proxy')
        return 0
    return body


while True:
    print("Main loop")
    method, header, body = channel.basic_get('tasks_todo')
    if method:
        print(f"Received task: {body}")
        channel.basic_ack(method.delivery_tag)
    else:
        print('didnt retrieve task')
        time.sleep(1)
        continue
    time.sleep(random.random()*10)

    if (not PROXY):
        PROXY = get_proxy()
        print(f"got proxy {PROXY}")
    time.sleep(random.random()*5)

    channel.basic_publish(exchange='',
                          routing_key='tasks_done',
                          body=f"task result for {body}")
    print(f"returned result for {body}")

    if (random.random() < 0.1):
        PROXY = None
        print("destroying proxy")
