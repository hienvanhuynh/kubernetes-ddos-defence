import subprocess
import time
import redis


r = redis.Redis()

def get_flow():
    flow = subprocess.run(['hubble', 'observe', '--since=5s', '-o', 'json'], stdout=subprocess.PIPE).stdout
    # print(flow)
    r.rpush('flow', flow)
    print('Save new flow into Redis')

while True:
    get_flow()
    time.sleep(5)