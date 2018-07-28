import React, {Component} from 'react';
import './Stylesheets/ChatOverviewCard.css';

var axios = require('axios')

class ChatOverviewCard extends Component {

  state = {
    home_members: [],
    away_members: []
  }

  componentDidMount() {
    var _this = this
    axios.get("/getTeamNames?team_id=" + this.props.data.UserTeamID)
        .then(function(result) {
          console.log(result);
          _this.setState({
            home_members: result.data
          })
        })
    if (this.props.data.FixtureID != -1) {
      axios.get("/getTeamNames?team_id=" + this.props.data.OppID)
          .then(function(result) {
            console.log(result);
            _this.setState({
              away_members: result.data
            })
          })
    }
  }

  isUpper(character) {
    return character === character.toUpperCase() && character !== " "
  }

  abbreviateName(name) {
    var uppers = name.split('').filter(this.isUpper)
    if (uppers.length === 2 || uppers.length === 3) {
      return uppers;
    } else {
      return name.substring(0,2).toUpperCase()
    }
  }

  chatAbbreviation(chat){
    var _this = this;
    if (chat.FixtureID === -1) {
      return (<div class="chat_abbreviation"><div class="chat_abbreviation_bold">
          <strong>{_this.abbreviateName(chat.UserTeamName)}</strong></div>
      </div>)
    } else {
      return (<div class="chat_abbreviation">vs<div class="chat_abbreviation_bold">
          <strong>{_this.abbreviateName(chat.OppName)}</strong></div>
      </div>)
    }
  }

  chatName(chat) {
    if (chat.FixtureID === -1) {
      return chat.UserTeamName
    } else {
      return chat.UserTeamName + " vs " + chat.OppName
    }
  }

  chatMembers() {
    var members = ""
    var i = 1;
    if (this.state.home_members.length > 0) {
      members += this.state.home_members[0].split(' ')[0]
      for (i = 1; i < this.state.home_members.length; i++) {
          members += ", " +  this.state.home_members[i].split(' ')[0];
      }
      if (this.state.away_members.length > 0) {
        members += ", "
      }
    }
    if (this.state.away_members.length > 0) {
      members += this.state.away_members[0]
      for (i = 1; i < this.state.away_members.length; i++) {
          members += ", " +  this.state.away_members[i].split(' ')[0];
      }
    }
    return members
  }

  render() {
    return (
      <div class="ChatCard">
        {this.chatAbbreviation(this.props.data)}
        <div class="chat_name"> {this.chatName(this.props.data)} </div>
        <div class="chat_members"> {this.chatMembers()}</div>
      </div>
    )
  }
}

export default ChatOverviewCard;
