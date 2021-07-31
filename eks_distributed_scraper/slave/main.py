import time
import zmq

context = zmq.Context()
socket = context.socket(zmq.PULL)
socket.connect("tcp://localhost:5556")

while True:
    print("I'm alive slave")
    print(socket.recv_string())
    time.sleep(5)
