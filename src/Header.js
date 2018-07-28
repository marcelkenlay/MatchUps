import React, {Component} from 'react';
import './Stylesheets/Header.css';
import {Link} from 'react-router-dom';

class Header extends Component {
    render() {
        return (
            <div id='menupanel'>
                <h1>MatchUps</h1>
                <div id="menuwrapper">
                    <table>
                      <tr>
                        <td></td>
                        <td></td>
                        <td class="clickable"><Link to="/profile">Profile</Link></td>
                        <td class="clickable"><Link to="/teams">Teams</Link></td>
                        <td class="clickable"><Link to="/chat">Chat</Link></td>
                        <td></td>
                        <td class="clickable"><Link to="/">Log Out</Link></td>
                      </tr>
                    </table>
                </div>
            </div>
        );
    }
}

export default Header;
