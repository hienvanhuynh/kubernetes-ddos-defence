import time
import redis
import json


r = redis.Redis()

def get_flow():
    return r.lpop('flow')

def process_flow(s):
    package = json.loads(s)
    print(package["time"])


while True:
    flow = get_flow()
    packages = flow.splitlines()
    for package in packages:
        process_flow(package)
    time.sleep(5)