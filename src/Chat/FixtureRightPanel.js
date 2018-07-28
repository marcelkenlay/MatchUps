import React, {Component} from 'react';
import axios from 'axios';
import UserProfile from '../Profile/UserProfile';
import StandaloneSearchBox from "react-google-maps/lib/components/places/StandaloneSearchBox";

import "./Stylesheets/FixtureRightPanel.css"
import "../Stylesheets/master.css"

var moment = require('moment');

class FixtureRightPanel extends Component {

  constructor(props){
    super(props)
    this.state = {
      active_key: "1",
      playedBefore: true,
      currentTeamID: -1,
      upcoming: {
        homeTeam: "",
        awayTeam: "",
        sport: "",
        lat: 0.0,
        lng: 0.0,
        location: "(Location not found)",
        date: ""      },
      previous: {
        homeTeam: "",
        awayTeam: "",
        homeScore: -1,
        awayScore: -1,
        sport: ""
      },
      submittedScore: {
        scoreHome: -1,
        scoreAway: -1,
        submitterID: -1
      },
      inputScoreHome: "",
      inputScoreAway: ""
    };

  }

  queryUpcomingGame() {
    // Query DB for upcoming game
    var _this = this;
    axios.get('/getupcominggame?fixtureID=' + this.props.fixture_id)
         .then(function(response) {
           if (response.data == "fail\n") {
             return;
           }

           var prettydate = _this.formatTimeText(response.data.Date);

           var venue_latlng = response.data.LocLat + "," + response.data.LocLng;
           var request_url = "http://maps.googleapis.com/maps/api/geocode/json?latlng="
                             + venue_latlng


           var newupcoming = {
             homeTeam: response.data.HomeTeam,
             awayTeam: response.data.AwayTeam,
             sport:    response.data.Sport,
             date:     prettydate,
             location: ""
           }

           _this.setState({upcoming: newupcoming});

           axios.get(request_url)
                 .then(function(result) {

                   newupcoming = {
                     homeTeam: response.data.HomeTeam,
                     awayTeam: response.data.AwayTeam,
                     sport:    response.data.Sport,
                     date:     prettydate,
                     location: result.data.results[0].formatted_address
                   }

                   _this.setState({upcoming: newupcoming});
                 })
         });
  }

  queryPreviousGame() {
    // Query DB for upcoming game
    var _this = this;
    var currentTeamID = 21;
    var oppTeamID = 22;
    axios.get('/getprevgame?curr=' + currentTeamID +
               "&opp=" + oppTeamID)
         .then(function(response) {
           if (response.data == "fail\n") {
             return;
           }

           var prev = {
             homeTeam: response.data.HomeTeam,
             awayTeam: response.data.AwayTeam,
             homeScore: response.data.HomeScore,
             awayScore: response.data.AwayScore,
             sport: response.data.Sport
           }

           _this.setState({previous: prev});
         });
  }

  querySubmittedScore() {
    // Query DB for a submittedScore game
    var _this = this;
    axios.get('/getsubmittedscore?fixtureID=' + this.props.fixture_id)
         .then(function(response) {
           if (response.data == "played\n") {

             var newsubscore = {
               scoreHome: -1,
               scoreAway: -1,
               submitterID: 0
             }

             _this.setState({
               submittedScore: newsubscore
             })

             return;
           }

           var newsubscore = {
             scoreHome: response.data.HomeScore,
             scoreAway: response.data.AwayScore,
             submitterID: response.data.SubmitterID
           }

           _this.setState({
             submittedScore: newsubscore
           })
         });
  }

  componentDidMount() {
    this.queryUpcomingGame();
    this.querySubmittedScore();
    this.queryPreviousGame();
  }


  pretty_date(date) {
    var day = moment(date)
    return day.format("dddd, MMMM Do YYYY, H:mm")
  }

  formatTimeText(date){
    return this.pretty_date(date);
  }


  handleSelect(key) {
    this.setState({
      active_key: key
    })
  }

  location_link() {
    return "http://maps.google.com/maps?q=" + this.state.upcoming.lat + "," + this.state.upcoming.loc
  }

  getUpcomingFixtureText() {
    // if (this.state.upcoming.homeTeam == "") {
    //   return ""
    // }

    return (
      <div>
        <p>
          The upcoming <span class="thintext bigtext">{this.state.upcoming.sport}</span> match between <span class="thintext bigtext">{this.state.upcoming.homeTeam}</span> and <span class="thintext bigtext">{this.state.upcoming.awayTeam}</span> is at <a id="rpvenue" href={this.location_link()} target="_blank">{this.state.upcoming.location}</a> at <span class="thintext bigtext">{this.state.upcoming.date}</span>.
        </p>

        <hr />

        {this.getPendingResultText()}

      </div>
    )
  }

  getPreviousFixtureText() {
    if (this.state.previous.homeTeam == "") {
      return ""
    } else {

       if (this.state.previous.scoreHome > this.state.previous.awayScore) {
         return (
           <div>
             <p>The last time you met, <span class="thintext bigtext">{this.state.previous.homeTeam}</span> won by a score of <span class="thintext bigtext">{this.state.previous.homeScore}</span> to <span class="thintext bigtext">{this.state.previous.awayScore}</span></p>
           </div>
         )
       } else if (this.state.previous.homeScore < this.state.previous.awayScore) {
         return (
           <div>
             <p>The last time you met, <span class="thintext bigtext">{this.state.previous.awayTeam}</span> won by a score of <span class="thintext bigtext">{this.state.previous.awayScore}</span> to <span class="thintext bigtext">{this.state.previous.homeScore}</span></p>
           </div>
         )
       } else {
         return (
           <div>
             <p>The last time you met, it was a <span class="thintext bigtext">draw</span> with each team scoring <span class="thintext bigtext">{this.state.previous.awayScore}</span></p>
           </div>
         )
       }

    }
  }

  submitScore(e) {
    e.preventDefault();
    var _this = this;

    // TODO: add current team to query
    axios.get('/submitscore?fixtureID=' + this.props.fixture_id +
              "&scH=" + this.state.inputScoreHome +
              "&scA=" + this.state.inputScoreAway)
         .then(function(response) {
            if (response.data == "fail\n") {
              // There is already a submitted score so update state to show
              // the current user
              _this.querySubmittedScore();
            }
         });

    this.setSubmitterZero();
  }

  setSubmitterZero() {
    var newsubscore = {
     scoreHome: -1,
     scoreAway: -1,
     submitterID: 0
    }

    this.setState({
     submittedScore: newsubscore
    })
  }

  acceptScore(e) {
    e.preventDefault();
    var _this = this;
    axios.get('/acceptsubmittedscore?fixture_id=' + this.props.fixture_id);
    this.setSubmitterZero();
  }

  rejectScore(e) {
    e.preventDefault();
    var _this = this;
    axios.get('/rejectsubmittedscore?fixtures_id=' + this.props.fixture_id);

    var defaultSubScore = {
      scoreHome: -1,
      scoreAway: -1,
      submitterID: -1
    }

    _this.setState({
      submittedScore: defaultSubScore
    })
  }

  inputScoreHomeChange(e) {
    const value = e.target.value;
    this.setState({
        inputScoreHome: value
    })
  }

  inputScoreAwayChange(e) {
    const value = e.target.value;
    this.setState({
        inputScoreAway: value
    })
  }

  getPendingResultText() {
    if (this.state.submittedScore.submitterID == 0) {
      return (
        <h3>No more scores to confirm.</h3>
      )
    }

    if (this.state.submittedScore.submitterID == -1) {
      return (
        <div>
          <p>Already played? Enter the score below!</p>

          <table id="recordScoreTable">
            <tr>
              <td>
                Home ({this.state.upcoming.homeTeam})
              </td>
              <td>
                versus
              </td>
              <td>
                Away ({this.state.upcoming.awayTeam})
              </td>
            </tr>
            <tr>
              <td>
                <input value={this.state.inputScoreHome} onChange={e => this.inputScoreHomeChange(e)} type='text' id="scoreHomeInput" class="scoreInput" />
              </td>
              <td>
                <h1> - </h1>
              </td>
              <td>
                <input value={this.state.inputScoreAway} onChange={e => this.inputScoreAwayChange(e)} type='text' id="scoreAwayInput" class="scoreInput" />
              </td>
            </tr>
          </table>

          <h3 class="centertext" onClick={e => this.submitScore(e)}>Submit</h3>
          <p>NOTE: the other team must confirm the score before it is saved</p>
        </div>
      )
    } else {

      if (this.state.submittedScore.submitterID == this.state.currentTeamID) {
        return (
          <h3>Waiting for the other team to confirm the score.</h3>
        )
      } else {
        return (
          <div>
            <h4>The other team submitted a score of<br /></h4>
            <h3>{this.state.upcoming.homeTeam}: {this.state.submittedScore.scoreHome} - {this.state.upcoming.awayTeam}: {this.state.submittedScore.scoreAway}</h3>
            <h4 class='centertext'>Is this right?</h4>
            <table id="scorefeedback">
              <tr>
                <td id="yesop" onClick={e => this.acceptScore(e)}>
                  Yes
                </td>
                <td></td>
                <td id="noop" onClick={e => this.rejectScore(e)}>
                  No
                </td>
              </tr>
            </table>
          </div>
        )
      }

    }
  }

  render() {
    var _this = this
    return (
      <div id="rightPanelWrapper">

        {this.getUpcomingFixtureText()}

        <hr />

        {this.getPreviousFixtureText()}

      </div>
    )
  }
}

export default FixtureRightPanel;
