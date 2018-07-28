import React, {Component} from 'react';
import UserProfile from '../Profile/UserProfile'
import ActiveUserID from '../Profile/ActiveUserID'
import TeamInvitationSearch from './TeamInvitationSearch'
import {FormControl, ListGroup, ListGroupItem, Button, CustomComponent} from 'react-bootstrap';
import './Stylesheets/CreateTeam.css'

var axios = require('axios');

class CreateTeam extends Component {

  state = {
    team_name: "",
    invitees: []
  }

  constructor(props) {
    super(props)
    this.addInvitee = this.addInvitee.bind(this)
    this.removeInvitee = this.removeInvitee.bind(this)
  }

  addInvitee(inv) {
    this.setState({
      invitees: [...this.state.invitees, inv],
    })
  }


  removeInvitee(inv) {
    var _this = this
    this.setState({
      invitees: _this.state.invitees.filter(_inv => _inv.Username !== inv.Username)
    })
  }

  teamInputChange(e) {
   const value = e.target.value;
   var _this = this;
   this.setState({
     team_name: value
   })
  }

  send_invitations(team_id) {
    var invitees = this.state.invitees.map(inv => inv.UserID)
    var teamInvInfo = {
      TeamID:   team_id,
      Invitees: invitees
    }
    axios.post("/sendInvitations", teamInvInfo)
    this.setState({
      invitees: []
    })
  }

  create() {
    var _this = this
    var teamInfo = {
      TeamName:this.state.team_name,
      CaptainID:ActiveUserID.getID(),
    }
    axios.post("/createTeam", teamInfo)
      .then(function(response){
        if (response.data != -1) {
          _this.send_invitations(response.data)
          _this.props.history.push("/chat")
        } else {
          alert("Team name not unique")
        }
      });
  }

  render() {
    var _this = this
    return (
        <div class="CreateTeam">
          <div id="teamNameSection">
            <h2>Team Name</h2>
            <input
              id="TeamNameBox"
              type="text"
              value={this.state.team_name}
              placeholder="Enter team name"
              onChange={e => this.teamInputChange(e)}
            />
          </div>
          <TeamInvitationSearch addInvitee={(_this.addInvitee)}
                                removeInvitee={_this.removeInvitee}
                                invitees={this.state.invitees}/>
          <Button onClick={() => this.create()}>Create Team and Send Invites</Button>
        </div>
    );
  }


}

export default CreateTeam
