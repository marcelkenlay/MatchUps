import React, { Component } from 'react';
import {FormGroup, FormControl, HelpBlock, ControlLabel, Col, Button, Form} from 'react-bootstrap';
import UserProfile from '../Profile/UserProfile';
import ActiveUserID from '../Profile/ActiveUserID';
import StandaloneSearchBox from "react-google-maps/lib/components/places/StandaloneSearchBox";
import LocationPicker from 'react-location-picker';
import {Link} from 'react-router-dom';

import './Stylesheets/NewUserForm.css';
import '../Stylesheets/master.css';


const defaultPosition = {
    lat: 51.509865,
    lng: -0.118092
};

const refs = {}

var axios = require('axios');

//Faster than regex
function isAlphaNumeric(str) {
  var code, i, len;
  for (i = 0, len = str.length; i < len; i++) {
    code = str.charCodeAt(i);
    if (!(code > 47 && code < 58) &&  // numeric (0-9)
        !(code > 64 && code < 91) &&  // upper alpha (A-Z)
        !(code > 95 && code < 123)) { // underscore and lower alpha (a-z)
      return false;
    }
  }
  return true;
};

function isNameValid(str) {
  var code, i, len;
  for (i = 0, len = str.length; i < len; i++) {
    code = str.charCodeAt(i);
    if (!(code == 32)             &&   // space
        !(code > 64 && code < 91) &&   // upper alpha (A-Z)
        !(code > 96 && code < 123) ) { // lower alpha (a-z)
      return false;
    }
  }
  return true;
};

class NewUserForm extends Component {
  state = {
     username: "",
     full_name: "",
     dob: "",
     position: defaultPosition,
     pwd: "",
     pwd2: ""
  }

  constructor(props) {
      super(props);
      this.state = {
         username: "",
         full_name: "",
         dob: "",
         position: defaultPosition,
         pwd: "",
         pwd2: ""
      }

      this.handleLocationChange = this.handleLocationChange.bind(this);
      this.onPlacesChanged = this.onPlacesChanged.bind(this);
      this.onSearchBoxMounted = this.onSearchBoxMounted.bind(this);
  }

  //Map
  handleLocationChange({position, address}) {
    this.setState({position, address});
  }

  onSearchBoxMounted(ref) {
    refs.searchBox = ref;
  }

  onPlacesChanged(){
    const places = refs.searchBox.getPlaces();
    this.setState({ places,});
    var lastvisited = places[places.length - 1];

    var object = {
        position: {
            lat: lastvisited.geometry.location.lat(),
            lng: lastvisited.geometry.location.lng(),
        },
        address: lastvisited.formatted_address
    };
    this.setState(object);
  }


  //Text entry validation
  getUsernameValidationState() {
    var _error = null
    const length = this.state.username.length;
    if (length >= 30 || !isAlphaNumeric(this.state.username)) {
      return 'error'
    }
  }

  getFullNameValidationState() {
    const length = this.state.full_name.length;
    if (length > 0 && length < 30 && isNameValid(this.state.full_name)) {
      return 'success';
    } else if (length> 0 ) return 'error';
    return 'error';

  }

  getDobValidationState() {
    const length = this.state.dob.length;
    if (length === 10) return 'success';
    else if (length > 0) return 'error';
    return null;
  }

  getPasswordValidationState() {
    const length = this.state.pwd.length;
    if (length >= 8 && length < 30) return 'success';
    else if (length > 0) return 'error';
    return null;
  }

  getPasswordResubmitValidationState() {
    const length = this.state.pwd2.length;
    if (length > 0 && (this.state.pwd === this.state.pwd2)) return 'success';
    else if (length > 0) return 'error';
    return null;
  }

  // Text Box Change Handlers
  usernameChange(e) {
    const value = e.target.value;
    this.setState({
        username: value
    })
  }

  fullNameChange(e) {
    const value = e.target.value;
    this.setState({
        full_name: value
    })
  }

  dobChange(e) {
    const value = e.target.value;
    this.setState({
        dob: value
    })
  }

  pwdChange(e) {
    const value = e.target.value;
    this.setState({
        pwd: value
    })
  }

  pwd2Change(e) {
      const value = e.target.value;
      this.setState({
          pwd2: value
      })
  }


  // Form components
  usernameGroup() {
    return <FormGroup controlId="formHorizontalUsername"
            validationState={this.getUsernameValidationState()}>
      <Col componentClass={ControlLabel} sm={2}>
        Username
      </Col>
      <Col sm={10}>
        <FormControl type="text" placeholder="Username" value={this.state.username}
          onChange={e => this.usernameChange(e)}/>
            <FormControl.Feedback />
            <HelpBlock>Validation is based on username length and whether it is unique.</HelpBlock>
      </Col>
    </FormGroup>
  }

  fullNameGroup() {
    return <FormGroup controlId="formHorizontalFullName"
            validationState={this.getFullNameValidationState()}>
      <Col componentClass={ControlLabel} sm={2}>
        Full Name
      </Col>
      <Col sm={10}>
        <FormControl type="text" placeholder="Enter Full Name" value={this.state.full_name}
          onChange={e => this.fullNameChange(e)}/>
            <FormControl.Feedback />
            <HelpBlock>Validation is based on string length.</HelpBlock>
      </Col>
    </FormGroup>
  }

  dobGroup() {
    return <FormGroup controlId="formHorizontalDOB"
            validationState={this.getDobValidationState()}>
      <Col componentClass={ControlLabel} sm={2}>
        Date Of Birth
      </Col>
      <Col sm={10}>
        <FormControl type="date" value={this.state.dob}
          onChange={e => this.dobChange(e)}/>
            <FormControl.Feedback />
            <HelpBlock class='helptext'>Date must be of from DD/MM/YYYY.</HelpBlock>
      </Col>
    </FormGroup>
  }

  locationGroup() {
    return <FormGroup controlId="formHorizontalLocation">
      <Col componentClass={ControlLabel} sm={2}>
        Location
      </Col>
      <Col sm={10}>
      <div id="formSearchBox">
        <StandaloneSearchBox
          ref={this.onSearchBoxMounted}
          bounds={this.bounds}
          onPlacesChanged={this.onPlacesChanged}
        >

        <input
          id="newuserloctext"
          type="text"
          placeholder="Search for your location"
          onKeyPress={event => {
            if (event.key === 'Enter') {
              this.onPlacesChanged()
              alert("In key press")
              event.preventDefault()
            }
          }}
          style={{
            boxSizing: `border-box`,
            border: `1px solid transparent`,
            width: `99%`,
            height: `32px`,
            padding: `0 12px`,
            borderRadius: `3px`,
            boxShadow: `0 2px 6px rgba(0, 0, 0, 0.3)`,
            fontSize: `14px`,
            outline: `none`,
            textOverflow: `ellipses`,
          }}
        />
        </StandaloneSearchBox>
      </div>
      <br />
        <LocationPicker
          name='Location'
          id='input2'
          value={this.state.position}
          containerElement={ <div style={ {height: '70%'} }/> }
          mapElement={ <div style={ {height: '400px'} }/> }
          defaultPosition={this.state.position}
          onChange={this.handleLocationChange} />
        <br />
      </Col>
    </FormGroup>
  }

  passwordGroup() {
    return (<FormGroup controlId="formHorizontalPassword"
            validationState={this.getPasswordValidationState()}>
      <Col componentClass={ControlLabel} sm={2}>
        Password
      </Col>
      <Col sm={10}>
        <FormControl type="password" placeholder="Password" value={this.state.pwd}
          onChange={e => this.pwdChange(e)}/>
        <FormControl.Feedback />
        <HelpBlock>Passwords must be at least 8 characters.</HelpBlock>
      </Col>
    </FormGroup>);
  }

  resubmitPasswordGroup() {
    return (<FormGroup controlId="formHorizontalPassword2"
            validationState={this.getPasswordResubmitValidationState()}>
      <Col componentClass={ControlLabel} sm={2}>
        Re-submit Password
      </Col>
      <Col sm={10}>
        <FormControl type="password" placeholder="Password" value={this.state.pwd2}
          onChange={e => this.pwd2Change(e)}/>
          <FormControl.Feedback />
          <HelpBlock>Passwords must match</HelpBlock>
      </Col>
    </FormGroup>);
  }

  submitGroup() {
    return (<FormGroup>
        <Col smOffset={2} sm={10}>
          <Button type="submit">Create User</Button>
        </Col>
    </FormGroup>)
  }

  // Form Submission
  create() {
    if (this.state.username.length > 30 || this.state.username.length == 0) {
      alert("Resolve issue with username.")
    } else if (this.getFullNameValidationState() !== 'success') {
      alert("Resolve issue with name.")
    } else if (this.getDobValidationState() !== 'success') {
      alert("Resolve issue with date of birth.")
    } else if (this.getPasswordValidationState() !== 'success') {
      alert("Resolve issue with password.")
    } else if (this.getPasswordResubmitValidationState() !== 'success') {
      alert("Ensure passwords match.")
    } else {var _this = this
      const req = "doesMatchingUserExist?username=" + this.state.username
      this.serverRequest =
        axios
          .get(req)
          .then(function(result) {
            // console.log(result.data);
            if (!result.data) {
              _this.finishCreation()
            } else {
              alert("Username is already in use")
            }
          })
    }
  }

  finishCreation() {
    var _this = this

    var body = {
      Username: this.state.username,
      Name: this.state.full_name,
      Dob: this.state.dob,
      LocLat: this.state.position.lat,
      LocLng: this.state.position.lng,
      Pwd: this.state.pwd
    }
    this.serverRequest =
      axios
        .post("addUserInfo", body).then(function(result) {
          ActiveUserID.setID(result.data)
          _this.proceedToProfile()
        })
  }

  proceedToProfile() {
    UserProfile.setName(this.state.username)
    this.props.history.push('/profile')
  }

  render() {
    return (
      <div id="contentcontainer">
        <h2 id="newusertitle"> Enter Your Details Below </h2>
        <Form horizontal onSubmit={(e) => {this.create(); e.preventDefault();}}>
          {this.fullNameGroup()}
          {this.usernameGroup()}
          {this.dobGroup()}
          {this.locationGroup()}
          {this.passwordGroup()}
          {this.resubmitPasswordGroup()}
          {this.submitGroup()}
        </Form>
        <h3 id='backtosplash' class='centertext'><Link to="/">or login to an existing account here</Link></h3>
      </div>
      );
  }
}

export default NewUserForm;
