FROM node:16-buster
# set working direction
WORKDIR /kkdgui-server
# add `/app/node_modules/.bin` to $PATH
ENV PATH /kkdgui-sever/node_modules/.bin:$PATH
# install application dependencies
COPY package.json ./
COPY package-lock.json ./

RUN apt-get install curl
RUN npm i

# add app
COPY . ./

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin

# start app
CMD ["node", "app.js"]
