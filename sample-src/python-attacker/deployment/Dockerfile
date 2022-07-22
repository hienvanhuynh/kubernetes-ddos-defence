FROM python:3.6
RUN apt update && yes | apt install hping3
COPY src /home/server
COPY ihulk.py /home/server/ihulk
CMD python3 /home/server/main.py
WORKDIR /home/server