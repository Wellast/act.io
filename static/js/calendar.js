const retrieveCalendarEvents = (token) =>
    fetch(`/listEvents?token=${token}&tzoffset=${(new Date()).getTimezoneOffset()}`)
        .then((res) => res.json())

const drawEvents = (events) => {
    document.getElementById('events-div').hidden = false;
    document.getElementById('events-list-text').hidden = false;

    const eventsDiv = document.getElementById('events-spreasheet');
    events.map((event) => {
        console.log(event);

        const start = new Date(event.start.dateTime || event.start.date);
        document.getElementById('events').innerHTML += `<a href='${event.htmlLink}'>${start.getHours()}:${start.getMinutes()} - ${event.summary}</a></br>`;

        const dayStart = new Date();
        dayStart.setUTCHours(0,0 + dayStart.getTimezoneOffset(),0,0);
        event.start.dateTime = new Date(event.start.dateTime);
        event.start.dateTime.setDate(dayStart.getDate());
        event.start.dateTime.setMonth(dayStart.getMonth());
        event.start.dateTime.setYear(dayStart.getFullYear());
        event.end.dateTime = new Date(event.end.dateTime);
        event.end.dateTime.setDate(dayStart.getDate());
        event.end.dateTime.setMonth(dayStart.getMonth());
        event.end.dateTime.setYear(dayStart.getFullYear());

        const eventDuration = (event.end.dateTime - event.start.dateTime) / 1000;
        const partWidth = Math.round(eventDuration/(24*60*60)*1000)/10;

        const partBeginOffset = Math.round(
            ((event.start.dateTime - dayStart) / 1000) / (24*60*60)*1000
        ) / 10;

        // const newEventTitle = document.createElement('p');
        // newEventTitle.innerText = `${event.summary}`;
        // newEventTitle.style.color = 'white';
        // newEventTitle.style.fontSize = '16px';

        const newEvent = document.createElement('button');
        newEvent.onclick = () => window.open(event.htmlLink, "_blank").focus();
        newEvent.type = 'button';
        newEvent.style.width = `${partWidth}%`;
        newEvent.style.marginLeft = `${partBeginOffset}%`;
        newEvent.style.height = '30px';
        newEvent.style.position = 'absolute';
        newEvent.style.top = '60px';
        newEvent.style.right = 0;
        newEvent.style.left = 0;
        newEvent.style.textOverflow = 'ellipsis';
        newEvent.style.overflow = 'hidden';
        newEvent.style.whiteSpace = 'nowrap';
        newEvent.style.padding = '0';
        newEvent.style.fontSize = '16px';
        newEvent.style.textOverflow = 'clip';
        newEvent.innerText = event.summary;
        newEvent.classList.add('rounded', 'shadow', 'btn', 'btn-secondary', 'btn-sm');
        // newEvent.appendChild(newEventTitle);
        eventsDiv.appendChild(newEvent);
    });
}
