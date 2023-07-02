FROM python:3.10-slim

WORKDIR /code

COPY . .

RUN pip install -r requirements.txt

CMD ["python", "-u", "./kafka_to_minio.py"]