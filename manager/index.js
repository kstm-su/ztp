const express = require('express');
const bodyParser = require('body-parser');
const cors = require('cors');

const app = express();

app.use(express.static('public'));
app.use(bodyParser.json());
app.use(cors());
app.use(require('./db'));
app.use('/api', require('./controllers'));

app.listen(process.env.PORT);
