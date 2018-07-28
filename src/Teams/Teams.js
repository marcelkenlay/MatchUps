import React, {Component} from 'react';
import "./Stylesheets/MainGrid.css"
import "../Stylesheets/master.css"
import Invitations from "./Invitations"
import TeamsCard from "./TeamsCard"
import {Button} from 'react-bootstrap'
import {Link} from 'react-router-dom';

export default class Teams extends Component {

    render() {
        return (
            <div id="teamwrapper">

                <table id="teamtable">
                  <tr>

                    <td id="teamleft">
                      <div id="teaminfo">
                          <TeamsCard/>
                      </div>
                    </td>
                    
                    <td id="teamright">
                      <div id="teamopts">

                        <div class="createTeam">
                          <Link class="centertext" to="/createTeam">
                            <h3>Create a Team</h3>
                          </Link>
                        </div>

                        <div class="invitation">
                          <Invitations/>
                        </div>

                      </div>
                    </td>

                  </tr>
                </table>






            </div>
          );
    }
}
