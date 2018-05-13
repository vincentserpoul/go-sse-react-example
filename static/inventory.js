if (typeof EventSource === "undefined") {
  document.getElementById("result").innerHTML =
    "Sorry, your browser does not support server-sent events...";
} else {
  const source = new EventSource("http://localhost:3000/sse");
  source.onmessage = event => {
    console.log(event);
  };

  const Welcome = props => <h1>Hello, {props.name}</h1>;

  class App extends React.Component {
    render() {
      return <Welcome name="Vincent" />;
    }
  }

  ReactDOM.render(<App />, document.getElementById("root"));
}
