import React, {Component} from 'react';
import Popover from "react-bootstrap/es/Popover";
import './Stylesheets/FixtureMap.css'

var moment = require('moment');

class FixtureMapPopover extends Component {

   pretty_date(date) {
    var day = moment(date)
    if (day.isSame(moment(), 'week')) {
      return (<div> This {day.format("dddd")} <br/> {day.format("h a")} </div>)
    } else if (day.subtract(1,'week').isSame(moment(), 'week')) {
      return (<div> Next {day.format("dddd")} <br/> {day.format("h a")} </div>)
    } else {
      return (<div> {day.format("Do ddd MMM")} <br/> {day.format("h a")} </div>)
    }
  }

  render() {
    return (<div class={(this.props.fixture.AdID===this.props.selectedFixture) ? "selected-popover" : "fixture-popover"}>
              {this.pretty_date(this.props.fixture.StartTime)}
            </div>)
  }
}

export default FixtureMapPopover
