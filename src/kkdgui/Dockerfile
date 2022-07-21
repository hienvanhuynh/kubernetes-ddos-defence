FROM node:16-alpine
# set working direction
WORKDIR /kkdgui
# add `/app/node_modules/.bin` to $PATH
ENV PATH /kkdgui/node_modules/.bin:$PATH
# install application dependencies
COPY package.json ./
COPY package-lock.json ./
RUN npm i
# add app
COPY . ./
# start app
CMD ["npm", "start"]