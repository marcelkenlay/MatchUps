import React, {Component} from 'react';
import Header from './Header';
import Main from './Main';
import ActiveUserID from './Profile/ActiveUserID'
import './Stylesheets/App.css';

class App extends Component {
    render() {
      if (ActiveUserID.getID() == -1) {
        this.props.history.push('/')
      }
      return (
        <div class="wrapper">
            <div class="header"><Header/></div>
            <div class="main"><Main/></div>
        </div>
      );
    }
}
export default App;
