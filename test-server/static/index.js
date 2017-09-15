var evtSource = new EventSource("/events");

evtSource.onmessage = function (e) {
    var newElement = document.createElement("li");

    newElement.innerHTML = "message: " + e.data;
    document.getElementById("events").appendChild(newElement);
}
