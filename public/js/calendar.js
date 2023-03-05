const drawEvents = (events) => {
    const eventsDiv = document.getElementById('events-div');
    events.map((event) => {
        const eventDuration = (new Date(event.end.dateTime) - new Date(event.start.dateTime)) / 1000;
        const partWidth = Math.round(eventDuration/(24*60*60)*1000)/10;

        var dayStart = new Date();
        dayStart.setUTCHours(0,0 + dayStart.getTimezoneOffset(),0,0);
        const partBeginOffset = Math.round(
            ((new Date(event.start.dateTime) - dayStart) / 1000) / (24*60*60)*1000
        ) / 10;

        console.log(event);

        const newEventTitle = document.createElement('p');
        newEventTitle.innerText = event.summary;
        newEventTitle.style.color = 'white';
        newEventTitle.style.fontSize = '16px';

        const newEvent = document.createElement('button');
        newEvent.type = 'button';
        newEvent.style.width = `${partWidth}%`;
        newEvent.style.marginLeft = `${partBeginOffset}%`;
        newEvent.style.height = '30px';
        newEvent.style.position = 'absolute';
        newEvent.style.top = '60px';
        newEvent.style.right = 0;
        newEvent.style.left = 0;
        newEvent.dataset.name = event.summary;
        newEvent.style.padding = '0';
        newEvent.classList.add('rounded', 'shadow', 'btn', 'btn-secondary', 'btn-sm');
        newEvent.appendChild(newEventTitle);
        eventsDiv.appendChild(newEvent);
    });
}
