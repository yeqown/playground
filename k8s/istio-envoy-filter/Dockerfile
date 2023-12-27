FROM python:3.8-slim-buster

WORKDIR /app

COPY . .

RUN pip install flask

EXPOSE 8080

CMD ["python", "app.py"]