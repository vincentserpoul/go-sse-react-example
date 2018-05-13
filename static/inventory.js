if (typeof EventSource === "undefined") {
  document.getElementById("result").innerHTML =
    "Sorry, your browser does not support server-sent events...";
} else {
  const source = new EventSource("http://localhost:3000/sse");
  source.onmessage = event => {
    console.log(event);
  };

  const name = "Josh Perez";
  const element = <h1>Hello, {name}</h1>;

  ReactDOM.render(element, document.getElementById("root"));
}
