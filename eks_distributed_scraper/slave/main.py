import os
import random
import time
import zmq
print("Started slave")

ZMQ_MASTER_HOSTNAME = os.environ.get("ZMQ_MASTER_HOSTNAME")
HOSTNAME = os.environ.get("HOSTNAME")

context = zmq.Context()
socket = context.socket(zmq.PULL)
master_addr=f"tcp://{ZMQ_MASTER_HOSTNAME}:5556"
print(master_addr)
socket.connect(master_addr)

while True:
    # print("I'm alive")
    print(f"Pulled package: {socket.recv_string()} on {HOSTNAME}")
    time.sleep(random.random()*10)
