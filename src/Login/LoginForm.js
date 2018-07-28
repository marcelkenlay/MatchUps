import React, { Component } from 'react';
import {FormGroup, FormControl, HelpBlock, ControlLabel, Col, Button, Form} from 'react-bootstrap';
import UserProfile from '../Profile/UserProfile';
import ActiveUserID from '../Profile/ActiveUserID';
import {Link} from 'react-router-dom';

import './Stylesheets/Login.css';
import '../Stylesheets/master.css';


var axios = require('axios');

class LoginForm extends React.Component {

  state = {
    username: "",
    pwd: ""
  }

  usernameChange(e) {
      const value = e.target.value;
      this.setState({
          username: value
      })
  }

  pwdChange(e) {
      const value = e.target.value;
      this.setState({
          pwd: value
      })
  }

  login() {
    var _this = this
    var body = {
      Username: this.state.username,
      Password: this.state.pwd
    }
    this.serverRequest =
      axios
        .post('checkLogin', body).then(function(result) {
          // console.log(result.data)
          if (result.data.Error == "none") {
            ActiveUserID.setID(result.data.UserID)
            _this.proceedToProfile()
          } else {
            alert(result.data.Error)
          }
        })
  }

  proceedToProfile() {
    UserProfile.setName(this.state.username)
    this.props.history.push('/profile')
  }

  usernameGroup() {
    return <FormGroup controlId="formHorizontal">

      <Col>
        <FormControl type="text" placeholder="Username" value={this.state.username}
          onChange={e => this.usernameChange(e)}/>
      </Col>
    </FormGroup>
  }

  passwordGroup() {
    return (<FormGroup controlId="formHorizontalPassword">

      <Col>
        <FormControl type="password" placeholder="Password" value={this.state.pwd}
          onChange={e => this.pwdChange(e)}/>
      </Col>
    </FormGroup>);
  }

  submitGroup() {
    return (<FormGroup>
        <Col>
          <Button type="submit">Sign in</Button>
        </Col>
    </FormGroup>)
  }

  render() {
    return (

      <div id='pagewrapper'>
        <div id='splashwrapper'>

          <div id="titlewrapper">
            <h1>MatchUps</h1>
          </div>

          <h3> Enter your details below to log in. </h3>
          <Form horizontal onSubmit={(e) => {this.login(); e.preventDefault();}}>
            {this.usernameGroup()}
            {this.passwordGroup()}
            {this.submitGroup()}
          </Form>
          <Link id="newuserlink" to="/newUser">or create a new account here</Link>
        </div>
      </div>




      );
  }
}

export default LoginForm;
