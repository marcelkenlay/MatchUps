import React, {Component} from 'react';
import './Stylesheets/MessageCard.css';
import UserProfile from '../Profile/UserProfile';
import ActiveUserID from '../Profile/ActiveUserID';

class FixtureCard extends React.Component {

  get_message_type() {
    if (this.props.sender_id == ActiveUserID.getID()) {
      return "sent"
    } else if (!this.props.is_fixture) {
      return "received_team"
    } else if (this.props.team_members.includes(this.props.sender_id)) {
      return "received_fixture_team"
    } else {
      return "received_fixture_opp"
    }
  }

  build_message_with_sender() {
    var message_class = this.get_message_type()
    return <div class={message_class} >
      <div class="sender"> {(message_class == "sent") ? "You" : this.props.sender_name} </div>
      <div class="message">
        <div class="message_text">{this.props.message}</div>
      </div>
    </div>
  }

  build_message_without_sender() {
    var message_class = this.get_message_type()
    return <div class={message_class} >
      <div class="message">
        <div class="message_text">{this.props.message}</div>
      </div>
    </div>
  }


  render() {
    return <div> {this.props.is_subsequent ? this.build_message_with_sender() :
              this.build_message_without_sender()}</div>
  }
}

export default FixtureCard;
