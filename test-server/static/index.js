var evtSource = new EventSource("/events");

evtSource.onmessage = function (e) {
    var newElement = document.createElement("li");

    newElement.innerHTML = "message: " + e.data;
    document.getElementById("events").appendChild(newElement);
}

evtSource.onerror = function (e) {
    console.log("EventSource failed.:", e);
};

$(document).ready(function () {
    $("form").submit(function (e) {
        e.preventDefault();

        var name = $(":text[name=name]").val()
        var message = $(":text[name=message]").val()

        $.post($(this).attr("action"), $(this).serialize());
    })
});
