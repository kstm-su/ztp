const http = require('http');
const express = require('express');
const bodyParser = require('body-parser');
const socketio = require('socket.io');
const cors = require('cors');

const app = express();
const server = http.Server(app);
const io = socketio(server);

app.use(express.static('public'));
app.use(bodyParser.json());
app.use(cors());
app.use((req, res, next) => {
  req.io = io;
  next();
});

app.use('/api', require('./controllers'));

server.listen(process.env.PORT);
process.env.SOCKET && app.listen(process.env.SOCKET);
