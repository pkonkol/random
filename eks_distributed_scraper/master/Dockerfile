FROM python:3.11.0a5-slim-bullseye

COPY requirements.txt /app/
RUN pip install --no-cache-dir -r /app/requirements.txt

COPY main.py /app
CMD ["python", "/app/main.py"]
