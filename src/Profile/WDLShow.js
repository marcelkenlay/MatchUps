import React, { Component } from 'react';
import '../Stylesheets/master.css';
import './Stylesheets/wdlshow.css';

class WDLShow extends Component {

  render() {
    return (
      <div id="wdl">
        <table id="recordScoreTable">
          <tr>
            <td>
              Wins
            </td>
            <td>
              Draws
            </td>
            <td>
              Losses
            </td>
          </tr>
          <tr>
            <td>
              <h1>{this.props.data[0]}</h1>
            </td>
            <td>
              <h1>{this.props.data[1]}</h1>
            </td>
            <td>
              <h1>{this.props.data[2]}</h1>
            </td>
          </tr>
        </table>
      </div>
    )
  }
}


export default WDLShow;
