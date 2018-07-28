import React, {Component} from 'react';
import ChatOverviewCard from './ChatOverviewCard';
import {ListGroup} from 'react-bootstrap';
import ActiveUserID from '../Profile/ActiveUserID'
import './Stylesheets/ChatList.css';

var axios = require('axios');

class ChatList extends Component {


    generateChatCard(chat) {
      var _this = this
      return (
        <div class={(chat == this.props.active_chat) ? "activeChatCardHolder" : "chatCardHolder"}
                onClick={() => this.props.setActiveChat(chat)}>
          <ChatOverviewCard key={`li-${chat.id}`} data={chat}/>
        </div>
      )
    }

    render() {
        return (
            <div id="chatlist">
                <h1>Chats</h1>
                {(this.props.chats.length > 0) ?
                  (<ListGroup>{
                    this.props.chats.map(chat => this.generateChatCard(chat))
                  }</ListGroup>)
                  : (<h3> No chats available, Join/Create a team first </h3>)}
            </div>
        );
    }
}

export default ChatList;
