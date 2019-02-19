import React, {Component} from 'react';
import Header from './Header';
import Main from './Main';
import './Stylesheets/App.css';

import Cookies from "universal-cookie";

const cookies = new Cookies();

const sessionIdKey = 'UserSessionId';

class App extends Component {
    render() {
        if (cookies.get(sessionIdKey) === null) {
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
