require('dotenv').config();
const express = require('express');
const bodyParser = require('body-parser');
const routes = require('./routes');
const app = express();
const port = process.env.PORT || 8000;

app.set('view engine', 'pug');
app.get('/favicon.ico', routes.favicon)
app.use(bodyParser.urlencoded({ extended: false }));
app.use('/static', express.static('static'));
app.get('/', routes.homepage);
app.get('/auth', routes.auth);
app.get('/listEvents', routes.listEvents);
app.get('/oauth2callback', routes.oauth2callback);

app.listen(port, console.log(`Server started on port http://localhost:${port}`));
