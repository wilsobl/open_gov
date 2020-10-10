
const AUTH0_CLIENT_ID = "9NH1MWLWcM54FRBn0Xvo2dOEFaJKG3gr";
const AUTH0_DOMAIN = "dev-vb3a8shg.us.auth0.com";
const AUTH0_CALLBACK_URL = location.href;
const AUTH0_API_AUDIENCE = "https://gordieh.opengov/";

class App extends React.Component {
  parseHash() {
    this.auth0 = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID
    });
    this.auth0.parseHash(window.location.hash, (err, authResult) => {
      if (err) {
        return console.log(err);
      }
      if (
        authResult !== null &&
        authResult.accessToken !== null &&
        authResult.idToken !== null
      ) {
        localStorage.setItem("access_token", authResult.accessToken);
        localStorage.setItem("id_token", authResult.idToken);
        localStorage.setItem(
          "profile",
          JSON.stringify(authResult.idTokenPayload)
        );
        window.location = window.location.href.substr(
          0,
          window.location.href.indexOf("#")
        );
      }
    });
  }

  setup() {
    $.ajaxSetup({
      beforeSend: (r) => {
        if (localStorage.getItem("access_token")) {
          r.setRequestHeader(
            "Authorization",
            "Bearer " + localStorage.getItem("access_token")
          );
        }
      }
    });
  }

  setState() {
    let idToken = localStorage.getItem("id_token");
    if (idToken) {
      this.loggedIn = true;
    } else {
      this.loggedIn = false;
    }
  }

  componentWillMount() {
    this.setup();
    this.parseHash();
    this.setState();
  }

  render() {
    if (this.loggedIn) {
      return <LoginHome />;
    }
    return <Home />;
  }
}

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.authenticate = this.authenticate.bind(this);
  }
  authenticate() {
    this.WebAuth = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID,
      scope: "openid profile",
      audience: AUTH0_API_AUDIENCE,
      responseType: "token id_token",
      redirectUri: AUTH0_CALLBACK_URL
    });
    this.WebAuth.authorize();
  }

  render() {
    return (
      <div className="container">
        <div className="row">
          <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
            <h1>Open-Gov</h1>
            <p>An open-source app for engaging citizens with government</p>
            <p>Sign in to get started </p>
            <a
              onClick={this.authenticate}
              className="btn btn-primary btn-lg btn-login btn-block"
            >
              Sign In
            </a>
          </div>
        </div>
      </div>
    );
  }
}

class LoginHome extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
        reps: []
    };

    this.serverRequest = this.serverRequest.bind(this);
    this.logout = this.logout.bind(this);
  }

  logout() {
    localStorage.removeItem("id_token");
    localStorage.removeItem("access_token");
    localStorage.removeItem("profile");
    location.reload();
  }

  serverRequest() {
    $.get("http://localhost:3000/api/localreps", res => {
      this.setState({
        reps: res
      });
    });
  }

  componentDidMount() {
    this.serverRequest();
  }

  render() {
    const userList = this.state.reps.users_rep_list;
    return (
      <div className="container">
        <br />
        <span className="pull-right">
          <a onClick={this.logout}>Log out</a>
        </span>
        <h2>Open-Gov</h2>
    <p>Hey user</p>
        <div className="row">
          <div className="container">
            {userList && (userList.map(function(localrep, i) {
              return <RepName key={i} localrep={localrep} />;
            }))}
          </div>
        </div>
      </div>
    );
  }
}


class RepName extends React.Component {
  constructor(props) {
    super(props);
  }
  render() {
    const localRep = this.props.localrep
    console.log("repName: ", localRep)
    return (
      <div className="col-xs-4">
        <div className="panel panel-default">
          <div className="panel-heading">
            {localRep.name}{" "}
            <span className="pull-right"></span>
          </div>
          <div className="panel-body joke-hld">Office: {localRep.office}</div>
          <div className="panel-body joke-hld">Location: {localRep.location}</div>
        </div>
      </div>
    );
  }
}
ReactDOM.render(<App />, document.getElementById("app"));