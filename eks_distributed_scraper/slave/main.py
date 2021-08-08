import os
import random
import time
import zmq
print("Started slave")

ZMQ_MASTER_HOSTNAME = os.environ.get("ZMQ_MASTER_HOSTNAME", 'localhost')
HOSTNAME = os.environ.get("HOSTNAME")
PROXY = None

# SOCKS_PROXY = None

context = zmq.Context()
task_sock = context.socket(zmq.PULL)
task_sock.connect(f"tcp://{ZMQ_MASTER_HOSTNAME}:5556")

prox_sock = context.socket(zmq.REQ)
prox_sock.connect(f"tcp://{ZMQ_MASTER_HOSTNAME}:5557")

result_sock = context.socket(zmq.PUSH)
result_sock.connect(f"tcp://{ZMQ_MASTER_HOSTNAME}:5558")


def get_proxy():
    print("sending proxy request")
    prox_sock.send_string("REQUESTING SOCKS PROXY")
    print("sent proxy request")
    return prox_sock.recv_string()


while True:
    print(f"Pulled task: {task_sock.recv_string()} on {HOSTNAME}")
    time.sleep(random.random()*10)
    if (not PROXY):
        PROXY = get_proxy()
        print(f"received proxy: {PROXY}")
    time.sleep(random.random()*5)
    result_sock.send_string("Returning result to master")
    print("returned result")
    if (random.random() < 0.1):
        PROXY = None
        print("destroying proxy")
