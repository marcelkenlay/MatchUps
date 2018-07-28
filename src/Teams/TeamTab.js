import React, {Component} from 'react';
import {Tab, Tabs} from "react-bootstrap";
import {connect} from "react-redux";
import {fetchTeams} from "../actions/teamsActions";
import "./Stylesheets/TeamCardsGrid.css";
import SimpleMap from "./MiniMap";
import TeamInvitationSearch from "./TeamInvitationSearch"
import {Button} from 'react-bootstrap';


var axios = require('axios');

class TeamTab extends Component {

    state = {
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

    send_invitations() {
        var invitees = this.state.invitees.map(inv => inv.UserID)
        var teamInvInfo = {
            TeamID: this.props.team.id,
            Invitees: invitees
        }
        axios.post("/sendInvitations", teamInvInfo)
        this.setState({
            invitees: []
        })
    }

    render() {
        return (
            <div class="TeamCardsGrid">
                <div class="teamname">
                    <h1>{this.props.team.name}</h1>
                </div>
                <div class="map">
                    <SimpleMap
                        team={this.props.team}/>
                </div>
                <div class="AddPlayer">
                    <h3>Invite more players</h3>
                    <TeamInvitationSearch addInvitee={(this.addInvitee)}
                                          removeInvitee={this.removeInvitee}
                                          invitees={this.state.invitees}/>
                    <Button onClick={() => this.send_invitations()}>Send Invitations</Button>
                </div>
            </div>
        );
    }
}

const mapStateToProps = (state, ownProps) => ({
    ...ownProps
});

export default connect(mapStateToProps)(TeamTab)
