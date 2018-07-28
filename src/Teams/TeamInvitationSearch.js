import React, {Component} from 'react';
import QueryResultCard from './QueryResultCard'
import UserProfile from '../Profile/UserProfile'
import ActiveUserID from '../Profile/ActiveUserID'
import './Stylesheets/CreateTeam.css'
import {FormControl, ListGroup, ListGroupItem, Button, CustomComponent} from 'react-bootstrap';

var axios = require('axios');

class TeamInvitationSearch extends Component {

  state = {
    query: "",
    query_results: [],
    team_name: ""
  }

  getQueryResults() {
    var _this = this;
    if (this.state.query.length === 0) {
      _this.setState({
        query_results: []
      })
      return
    }
    var req = "/getUsernameMatches?pattern=" + this.state.query
    this.serverRequest = axios.get(req)
        .then(function (result) {
          if (result.data != null) {
            var relevant_results = result.data.filter((res) => res.UserID !== ActiveUserID.getID())
            _this.setState({query_results: relevant_results});
          } else {
            _this.setState({query_results: []});
          }
        })
  }

  inputChange(e) {
   const value = e.target.value;
   var _this = this;
   this.setState({
     query: value
   }, () => {_this.getQueryResults()})
  }

  addInvitee(inv){
    var _this = this
    var matches = this.props.invitees.filter(_inv => _inv.UserID === inv.UserID)
    if (matches.length > 0) {
      return
    }
    this.props.addInvitee(inv)
    this.setState({
      query: "",
      query_results: []
    })
  }

  removeInvitee(inv){
    this.props.removeInvitee(inv)
  }

  inviteeListItem(inv) {
    var _this = this
    const username = inv.Username
    return (<div class="inviteeComponent">
        <div class="inviteeUsername">{"ID: " + username}</div>
        <div class="inviteeName">{"Name: " + inv.FullName}</div>
        <Button onClick={() => this.removeInvitee(inv)}>Remove</Button>
      </div>)
  }

  queryListItem(res) {
    return (<ListGroupItem onClick={() => this.addInvitee(res)}
      header={"ID: " + res.Username}>{"Name: " + res.FullName}</ListGroupItem>
    )
  }

  render() {
    var _this = this
    return (
      <div id="searchSection">
        <h3 class="nomargintop">User Search</h3>
        <input
          id="SearchBox"
          placeholder="Search by username"
          ref={input => this.search = input}
          onChange={e => this.inputChange(e)}
          value={this.state.query}
        />
        <ListGroup>
          {this.state.query_results
            .map(res => _this.queryListItem(res))}
        </ListGroup>
        <div id="invitationSection">
          <h3>Invitations will be sent to the following users:</h3>
          {
            (this.props.invitees.length === 0) ? "No users selected." :
            (<ListGroup>
              {
                this.props.invitees.map(inv => _this.inviteeListItem(inv))
              }
             </ListGroup>
           )
         }

        </div>
      </div>
    )
  }


}

export default TeamInvitationSearch
