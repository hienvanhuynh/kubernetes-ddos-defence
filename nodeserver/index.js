const express = require('express');
const app = express();
const cors = require('cors')


app.use(express.json());
app.use(express.urlencoded({
    extended: true
  }));
app.use(cors())
numAccess=0
app.use('*', (req, res) => {
  numAccess=numAccess+1
  console.log(":",numAccess, "\n")
  res.send('Welcome to our service!!!')
})

const PORT = 5050
app.listen(PORT, () => console.log(`Server is listening on port ${PORT}`))
