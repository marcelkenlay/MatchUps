import React, {Component} from 'react';
import FixtureRightPanel from './FixtureRightPanel';
import TeamRightPanel from './TeamRightPanel';
import {Nav, NavItem} from 'react-bootstrap';

class ChatRightPanel extends Component {

  handleSelect(key) {
    this.setState({
      active_key: key
    })
  }


  render() {
    var _this = this

    return (
      <div id="rightpanelwrapper">
        {
          (this.props.active_chat.FixtureID !== -1) ?
          (<FixtureRightPanel fixture_id={this.props.active_chat.FixtureID}/>) :
          (<TeamRightPanel get_chats={() => _this.props.get_chats()}
                           team_id={this.props.active_chat.UserTeamID}
                           channel={this.props.channel}/>)
        }
      </div>
    );
  }
}

export default ChatRightPanel;
