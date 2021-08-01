FROM python:3.9.6-slim

COPY requirements.txt /app/
RUN pip install --no-cache-dir -r /app/requirements.txt

COPY main.py /app
CMD ["python", "/app/main.py"]
