const { google } = require('googleapis');
const express = require('express');
const bodyParser = require('body-parser');
require('dotenv').config();

const clientId = process.env.CLIENT_ID;
const clientSecret = process.env.CLIENT_SECRET;
const redirectUrl = process.env.REDIRECT_URL;

const app = express();
const port = 8000;
const scopes = ['https://www.googleapis.com/auth/calendar.readonly'];
const oauth2Client = new google.auth.OAuth2(clientId, clientSecret, redirectUrl);

google.options({auth: oauth2Client});



const listEvents = async (req, res) => {
    oauth2Client.setCredentials(JSON.parse(decodeURIComponent(req.query.tokens)));

    var start = new Date();
    start.setUTCHours(0,0,0,0);

    var end = new Date();
    end.setUTCHours(23,59,59,999);

    const calendar = google.calendar({version: 'v3', auth: oauth2Client});
    const calendarRes = await calendar.events.list({
      calendarId: 'primary',
      timeMin: start.toISOString(),
      // timeMax: end.toISOString(),
      maxResults: 10,
    });
    const events = calendarRes.data.items;
    if (!events || events.length === 0) {
      return res.sendStatus(204);
    }
    // events.map((event, i) => {
    //   const start = event.start.dateTime || event.start.date;
    //   console.log(`${start} - ${event.summary}`);
    // });
    res.send(events);
}

const auth = async (_, res) => {
    const url = oauth2Client.generateAuthUrl({
        // 'online' (default) or 'offline' (gets refresh_token)
        access_type: 'offline',
        // If you only need one scope you can pass it as a string
        scope: scopes,
      });
    res.send(url)
}

const oauth2callback = async (req, res) => {
    const { tokens } = await oauth2Client.getToken(req.query.code);
    res.redirect(301, `/?tokens=${encodeURIComponent(JSON.stringify(tokens))}`);
}



app.use(bodyParser.urlencoded({ extended: false }));
app.use(express.static('public'));
app.use('/', express.static('public'));
app.get('/auth', auth);
app.all('/oauth2callback', oauth2callback);
app.get('/listEvents', listEvents);
app.listen(port, console.log(`Server started on port http://localhost:${port}`));
