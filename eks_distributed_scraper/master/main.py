import zmq
import time

context = zmq.Context()
socket = context.socket(zmq.PUSH)
socket.bind("tcp://*:5556")

while True:
    print("I'm alive")
    socket.send_string(f"I'm alive {time.ctime()}")
    time.sleep(20)
