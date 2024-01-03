import sys
import redis

def calculate_square(n):
    return n**2

def deliver_integer(r, n):
    r.publish("square", n) 

if __name__ == "__main__":
    r = redis.Redis(
    host='127.0.0.1',
    port=6379,
    decode_responses=True
    )
  
    deliver_integer(r, calculate_square(int(sys.argv[1])))