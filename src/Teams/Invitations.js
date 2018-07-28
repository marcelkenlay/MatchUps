import React, {Component} from 'react';
import {ListGroup, ListGroupItem, Button} from 'react-bootstrap'
import {declineInvitation, acceptInvitation, fetchInvitations} from "../actions/invitationsActions";
import {connect} from "react-redux";

import "./Stylesheets/invitation.css"

class Invitations extends Component {

    componentWillMount() {
        this.props.fetchInvitations();
    }

    render() {
        return (
            <div>
                <h3 id="invitations-header">Invitations from other teams:</h3> 
                <ListGroup>
                    {this.props.invitations.map((invitation) =>
                        <ListGroupItem class="invitation" header={invitation.name}>
                            <Button bsStyle="success"
                                    onClick={this.props.acceptInvitation.bind(this, invitation)}> Accept </Button>
                            <Button bsStyle="danger"
                                    onClick={this.props.declineInvitation.bind(this, invitation)}> Decline </Button>
                        </ListGroupItem>
                    )}
                </ListGroup>
            </div>
        );
    }
}

const mapStateToProps = state => ({
    invitations: state.invitations
});

export default connect(mapStateToProps, {acceptInvitation, declineInvitation, fetchInvitations})(Invitations)
