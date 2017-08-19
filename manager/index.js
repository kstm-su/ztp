const express = require('express');
const bodyParser = require('body-parser');

const app = express();

app.use(bodyParser.json());
app.use(require('./db'));
app.use(require('./controllers'));

app.listen(5001);
