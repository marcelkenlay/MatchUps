import React, { Component } from 'react';
import AvailabiltyTable from './Profile/AvailabiltyTable';
import axios from 'axios';
import UserProfile from './Profile/UserProfile';
import PreviousFixtureCard from './Profile/PreviousFixtureCard';
import UpcomingFixtureCard from './Profile/UpcomingFixtureCard';
import WDLShow from './Profile/WDLShow';
import StandaloneSearchBox from "react-google-maps/lib/components/places/StandaloneSearchBox";

import './Stylesheets/master.css';
import './Stylesheets/profile.css';
import './Stylesheets/Searchbox.css'

const refs = {}

class Profile extends Component {

  constructor(props){
    super(props)
    this.state = {
      username: "",
      location: "",
      position:{
        lat: 0.0,
        lng: 0.0,
      },
      places: [],
      fixtures: [],
      upcoming: [],
      isEditing: false,
      wdl: [0, 0, 0]
    };

    this.onPlacesChanged = this.onPlacesChanged.bind(this)
  }


  // Searchbox init
  onSearchBoxMounted(ref) {
      refs.searchBox = ref;
  }

  // On searchbox edited
  onPlacesChanged(){
      const places = refs.searchBox.getPlaces();
      this.setState({ places,});
      var lastvisited = places[places.length - 1];
      var newlat = lastvisited.geometry.location.lat()
      var newlng = lastvisited.geometry.location.lng()
      // console.log(lastvisited.formatted_address)
      var object = {
          position: {
              lat: newlat,
              lng: newlng,
          },
          location: lastvisited.formatted_address
      };
      this.setState(object);
      // console.log(this.state)
      axios.get(
        '/updateuserloc?username=' + UserProfile.getName()
        + '&lat=' + newlat
        + '&lng=' + newlng).then(function(response) {
          if (response.data == "fail\n") {
            alert("Failed to find location, please try again.")
          } else {
            var tick = document.getElementById('searchtick');
            tick.innerHTML = "✓";
            setTimeout(function() {
              tick.innerHTML = "";
            }, 2500);
          }
        });
  }

  loadUserInformation() {
    // Query DB
    var _this = this;
    var username = UserProfile.getName();
    axios.get('/getuserinfo?username='+username)
         .then(function(response) {
           var venue_latlng = response.data.LocLat + "," + response.data.LocLng;
           var request_url = "http://maps.googleapis.com/maps/api/geocode/json?latlng="
                             + venue_latlng

           axios.get(request_url)
                 .then(function(result) {
                   _this.setState({location: result.data.results[0].formatted_address});
                   _this.scrollFixtures();
                 })
         });
  }

  loadUserFixtures() {
    var _this = this;
    var username = UserProfile.getName();
    // axios.get("prevfix.json")
    axios.get('/getuserfixtures?username=' + username)
         .then(function(response) {
           _this.setState({
             fixtures: response.data
           });
           _this.calcWDL();
           _this.scrollFixtures();
         });
  }

  loadUpcoming() {
    var _this = this;
    var username = UserProfile.getName();
    // axios.get('/upfix.json')
    axios.get('/getuserupcoming?username=' + username)
        .then(function (response) {
            _this.setState({
                upcoming: response.data
            });
            _this.scrollFixtures();
        });
  }

  componentDidMount() {
    this.loadUserInformation();
    this.loadUserFixtures();
    this.loadUpcoming();
  }

  calcWDL() {
    var w = 0;
    var d = 0;
    var l = 0;

    for (var i = 0; i < this.state.fixtures.length; i++) {
      var fixture = this.state.fixtures[i];

      if (fixture.ScoreHome == fixture.ScoreAway) {
        d++;
        continue;
      }

      if (fixture.IsHome) {
        if (fixture.ScoreHome > fixture.ScoreAway) {
          w++;
        } else {
          l++;
        }
      } else {
        if (fixture.ScoreAway > fixture.ScoreHome) {
          w++;
        } else {
          l++;
        }
      }
    }
    var newwdl = [w, d, l]

    this.setState({
      wdl: newwdl
    });

  }

  showEditBox(e) {
    e.preventDefault();

    var _this = this;
    _this.setState({
      isEditing: true
    })
  }

  scrollFixtures() {
    var box = document.getElementById("carousel");
    box.scrollLeft = box.scrollWidth;
  }

  render() {
    return (
      <div id='contentpanel'>
        <div id='contentcontainer'>
          <h1 id='username' class='centertext'>{UserProfile.getName()}</h1>

          <WDLShow data={this.state.wdl} />

          <hr />

          <h3 class="centertext subheading">← scroll left for older | scroll right for newer →</h3>
          <div id="carousel">
            <table id="carousel-table">
              <tr id="carousel-row">
              {
                this.state.fixtures.length == 0 ? <b>(No previous matches)</b> :
                this.state.fixtures.map(item => (
                  <td class="carousel-cell">
                    <div class="cell-wrap">
                      <PreviousFixtureCard data={item} />
                    </div>
                  </td>
                  )
                )
              }
              {
                this.state.upcoming.length == 0 ? <b>(No upcoming matches)</b> :
                this.state.upcoming.map(item => (
                  <td class="carousel-cell">
                    <div class="cell-wrap">
                      <UpcomingFixtureCard data={item} />
                    </div>
                  </td>
                  )
                )
              }
              </tr>
            </table>
          </div>

          <hr />

          <h3 class='centertext subheading'>Your Details</h3>


          <h3 class='centertext notopmargin'>Location: <span class='thintext'>{this.state.location} <a id='locChangeLink' onClick={e => this.showEditBox(e)}>(change)</a></span></h3>

          <div id="changelocbox">
            {this.state.isEditing ?
              <div><span id='searchtick'></span><StandaloneSearchBox
                ref={this.onSearchBoxMounted}
                bounds={this.bounds}
                onPlacesChanged={this.onPlacesChanged}
              >
              <input
                type='text'
                placeholder="Search for your location"
                id = "searchBox"
                />
              </StandaloneSearchBox></div>
              : ""}
          </div>
          <div class="AvTable">
            <AvailabiltyTable />
          </div>



        </div>
      </div>
    );
  }
}


export default Profile;
