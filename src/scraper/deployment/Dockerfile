FROM golang:alpine
RUN apk add --update curl && apk add --update bash
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
COPY src /home/app
WORKDIR /home/app
RUN go build -o scraper . 
CMD ["./scraper"]