import React, {Component} from 'react';
import './Stylesheets/ProfileFixtureCard.css';

var axios = require('axios');
var moment = require('moment');

class PreviousFixtureCard extends Component {
  state = {
    Location: "Failed to find."
  }

  componentDidMount() {
    var venue_latlng = this.props.data.LocLat + "," + this.props.data.LocLng;
    var request_url = "http://maps.googleapis.com/maps/api/geocode/json?latlng="
                      + venue_latlng

    var _this = this
    _this.serverRequest = axios.get(request_url)
          .then(function(result) {
            // console.log(result);
            _this.setState({Location: result.data.results[0].formatted_address});
          })
  }

  getResult(item) {
    if (item.scoreFor == "") {
      return "unplayed";
    }

    if (item.ScoreHome == item.ScoreAway) {
      return "draw";
    }

    if (item.IsHome) {
      if (item.ScoreHome > item.ScoreAway) {
        return "win";
      }
      if (item.ScoreHome < item.ScoreAway) {
        return "lose";
      }
    } else {
      if (item.ScoreHome < item.ScoreAway) {
        return "win";
      }
      if (item.ScoreHome > item.ScoreAway) {
        return "lose";
      }
    }
  }

  pretty_date(date) {
    var day = moment(date)
    return day.format("dddd, MMMM Do YYYY, H:mm")
  }

  formatDate(date){
    return this.pretty_date(date);
  }

  getFixtureVsText(item) {
    var home = "";
    var away = "";

    if (item.IsHome) {
      home = item.ForTeam;
      away = item.Opposition;
    } else {
      home = item.Opposition;
      away = item.ForTeam;
    }

    return (
      <table class="teamvs">
        <tr>
          <td>
          <span class="teamname">{home}</span>
          </td>

          <td>
            vs
          </td>

          <td>
            <span class="teamname">{away}</span>
          </td>
        </tr>
      </table>
    )
  }

  render() {
    return (
      <div class={"resultcard " + this.getResult(this.props.data)}>
        <h3 class='centertext'>{this.formatDate(this.props.data.Date)}</h3>
        <p class='centertext'>{this.state.Location}</p>
        {this.getFixtureVsText(this.props.data)}
        <p class='centertext'>{this.props.data.Sport}</p>
        <h2 class='centertext notopmargin'>{this.props.data.ScoreHome} - {this.props.data.ScoreAway}</h2>
      </div>
    )
  }
}


export default PreviousFixtureCard;
