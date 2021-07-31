FROM python:3.9.6-slim

COPY main.py /app/
COPY requirements.txt /app/

RUN pip install --no-cache-dir -r /app/requirements.txt

CMD ["python", "/app/main.py"]
