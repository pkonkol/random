import time
import zmq
import os
print("Started slave")

ZMQ_MASTER_HOSTNAME = os.environ.get("ZMQ_MASTER_HOSTNAME")

context = zmq.Context()
socket = context.socket(zmq.PULL)
socket.connect(f"tcp://{ZMQ_MASTER_HOSTNAME}:5556")

while True:
    print("I'm alive slave")
    print(socket.recv_string())
    time.sleep(5)
