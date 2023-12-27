from flask import Flask, request

app = Flask(__name__)

@app.route("/")
def index():
    print(request.headers)

    dict = {}
    for key, value in request.headers:
        dict[key] = value

    return dict

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
