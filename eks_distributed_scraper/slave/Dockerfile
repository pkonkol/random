FROM python:3.9.6-slim

ENV ZMQ_MASTER_HOSTNAME=${ZMQ_MASTER_HOSTNAME:-localhost}

COPY requirements.txt /app/
RUN pip install --no-cache-dir -r /app/requirements.txt

COPY main.py /app/
CMD ["python", "/app/main.py"]
