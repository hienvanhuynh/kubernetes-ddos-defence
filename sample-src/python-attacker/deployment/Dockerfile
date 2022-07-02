FROM python:3.6
RUN apt update && yes | apt install hping3
COPY src /home/server
CMD python3 /home/server/main.py