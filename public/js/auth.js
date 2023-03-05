let creds = undefined;

const retrieveAuthUrl = () => {
    const xhttp = new XMLHttpRequest();
    xhttp.onload = function () {
        console.log(this.responseText);
        window.open(this.responseText, "_self").focus();
    }
    xhttp.open("GET", "/auth", true);
    xhttp.send();
}
document.getElementById("auth-button").addEventListener("click", retrieveAuthUrl);

const retrieveCalendarEvents = (creds) => {
    const xhttp = new XMLHttpRequest();
    xhttp.onload = function () {
        const events = JSON.parse(this.responseText);
        drawEvents(events);
        events.map((event, i) => {
            const start = event.start.dateTime || event.start.date;
            document.getElementById('events').innerHTML += `${start} - ${event.summary}</br>`;
        });
    }
    const dt = new Date();
    xhttp.open("GET", `/listEvents?tokens=${creds}&tzoffset=${dt.getTimezoneOffset()}`, true);
    xhttp.send();
}

const alreadyAuthorized = (tokens) => {
    const authButton = document.getElementById('auth-button');
    authButton.disabled = true;
    authButton.innerHTML = 'Authorized';
    console.log(tokens);
}

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
tokens = urlParams.get('tokens');
if (tokens) {
    localStorage.setItem('tokens', tokens);
    window.location.href = '/';
}

tokens = localStorage.getItem('tokens');
if (tokens) {
    alreadyAuthorized(tokens);
    retrieveCalendarEvents(tokens);
}