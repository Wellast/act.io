const { google } = require('googleapis');
require('dotenv').config();

const clientId = process.env.CLIENT_ID;
const clientSecret = process.env.CLIENT_SECRET;
const redirectUrl = process.env.REDIRECT_URL;
const scopes = ['https://www.googleapis.com/auth/calendar.readonly'];
const oauth2Client = new google.auth.OAuth2(clientId, clientSecret, redirectUrl);

google.options({auth: oauth2Client});

exports.homepage = (req, res) => {
    res.render('index', {})
}

exports.listEvents = async (req, res) => {
    // console.log(req.query.token)
    if (!req.query.token) {
        return res.sendStatus(400)
    }

    oauth2Client.setCredentials(JSON.parse(decodeURIComponent(req.query.token)));
    const tzOffset = Number(req.query.tzoffset) || 0;

    const start = new Date();
    start.setUTCHours(0,0 + tzOffset,0,0);
    const end = new Date();
    end.setUTCHours(23,59 + tzOffset,59,999);

    const calendar = google.calendar({ version: 'v3', auth: oauth2Client });
    let calendarRes;
    try {
        calendarRes = await calendar.events.list({
            calendarId: 'primary',
            timeMin: start.toISOString(),
            timeMax: end.toISOString(),
            maxResults: 10,
        });
    } catch (err) {
        return res.sendStatus(401);
    }
    const events = calendarRes.data.items;
    if (!events || events.length === 0) {
        return res.sendStatus(204);
    }
    return res.send(events);
}

exports.auth = async (_, res) => {
    // 'online' (default) or 'offline' (gets refresh_token)
    // If you only need one scope you can pass it as a string
    const url = oauth2Client.generateAuthUrl({ access_type: 'offline', scope: scopes });
    res.send(url);
}

exports.oauth2callback = async (req, res) => {
    if (!req || !req.query || !req.query.code) {
        return res.redirect(301, '/');
    }
    const { tokens } = await oauth2Client.getToken(req.query.code);
    return res.redirect(301, `/?token=${encodeURIComponent(JSON.stringify(tokens))}`);
}

exports.favicon = (req, res) => res.sendFile('static/favicon.ico', { root : __dirname});
