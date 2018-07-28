import React, {Component} from 'react';
import ReactDOM from 'react-dom';
import ChatList from './ChatList';
import OpenChat from './OpenChat';
import ChatRightPanel from './ChatRightPanel';
import Pusher from 'pusher-js';
import ActiveUserID from '../Profile/ActiveUserID'

import './Stylesheets/Chat.css';

var axios = require('axios');

class Chat extends Component {
  state = {
    active_chat: null,
    chats: [],
    pusher: new Pusher('112f8743d26528cd9b7e', {
        cluster: 'eu',
        encrypted: true
      }),
    channel: null
  }

  constructor(props) {
    super(props);
    this.getChats = this.getChats.bind(this);
    this.getChatName = this.getChatName.bind(this);
    this.initialiseActiveChat = this.initialiseActiveChat.bind(this);
  }

  initialiseActiveChat() {
    if (this.state.chats.length != 0) {
      this.setActiveChat(this.state.chats[this.state.chats.length - 1])
    }

  }
  getChats() {
    console.log("Getting chats");
    var _this = this;
    var query = "getChats?userID=" + ActiveUserID.getID()
    this.serverRequest =
        axios
            .get(query)
            .then(function (result) {
              console.log(result.data)
              console.log(_this.state.chats)
              if (result.data != null) {
                _this.setState({
                    chats: result.data
                }, _this.initialiseActiveChat())
              }
            })
  }

  componentDidMount() {
    this.getChats()
  }

  getChatName(chat) {
    var name = (chat.FixtureID !== -1) ?
                        ("_fixture" + chat.FixtureID)
                        : ("_team" + chat.UserTeamID)
    console.log(name)
    return name
  }

  setActiveChat(chat) {
    console.log("Setting active_chat");
    if (this.state.active_chat != null && chat != this.state.active_chat) {
      console.log("Unsubscribing from channel")
      this.state.pusher.unsubscribe(this.getChatName(this.state.active_chat))
    }
    this.setState({
      channel: this.state.pusher.subscribe(this.getChatName(chat)),
      active_chat: chat
    })
  }

  render() {
      var _this = this
      return (
          <div id="chatwrapper">
          <div id="left-block" class="chatpagepanel">
            <ChatList active_chat={this.state.active_chat}
                      chats={this.state.chats}
                      setActiveChat={(chat) => _this.setActiveChat(chat)}/>
          </div>

          <div id="central-block" class="chatpagepanel">
            {
              (this.state.active_chat != null) ?
              (<OpenChat active_chat={this.state.active_chat}
              is_fixture={this.state.active_chat.FixtureID != -1}
              channel={this.state.channel}/>)
              : (<h3>Select a chat in the left panel to view messages.</h3>)
            }
          </div>

          <div id="right-block" class="chatpagepanel">
            {
              (this.state.active_chat != null) ?
              <ChatRightPanel active_chat={this.state.active_chat}
                              get_chats={() => _this.getChats()}
                              channel={this.state.channel}/>
              : ""
            }
          </div>

          </div>
      )
    }
}


export default Chat;
