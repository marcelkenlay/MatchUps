import React, {Component} from 'react';
import FixtureCard from './FixtureCard';
import './Stylesheets/FixtureList.css';
import FixtureMap from './FixtureMap';

var axios = require('axios');

class FixtureList extends Component {
  state = {
      fixtures: [],
      selectedFixture: -1
  }

  constructor(props) {
    super(props)
    this.removeFixture = this.removeFixture.bind(this)
  }

  componentDidMount() {
    this.getPromotedFixtures()
  }

  getPromotedFixtures() {
    var _this = this;
    var req = "getPromotedFixtures?team=" + this.props.team_id
    this.serverRequest = axios.get(req)
        .then(function (result) {
            _this.setState({fixtures: result.data});
        })
  }

  removeFixture(id) {
    this.setState({
      fixtures: this.state.fixtures.filter(fixture => fixture.AdID != id)
    })
    this.props.get_chats()
  }

  componentDidUpdate(prevProps) {
    if (prevProps.team_id !== this.props.team_id) {
     this.getPromotedFixtures()
    }
  }

  setSelectedFixture(fixture) {
    this.setState({
      selectedFixture: fixture.AdID,
      fixtures: [fixture,
                  ...(this.state.fixtures.filter((fix) => fix.AdID !== fixture.AdID))]
    })
  }

  render() {
    var _this = this
      return (
        <div class="Fixtures">
          <div class="FixtureMap">
            <FixtureMap fixtures={this.state.fixtures}
                        setSelectedFixture={id => _this.setSelectedFixture(id)}
                        selectedFixture={this.state.selectedFixture}/>
          </div>
          <div class="FixtureList">
            <h3>Vote on fixtures from other teams below</h3>
            {
              this.state.fixtures.map(fixture =>
                (
                  <div>
                    <FixtureCard key={`li-${fixture.AdID}`}
                          is_selected={fixture.AdID === this.state.selectedFixture}
                          remove_fixture={this.removeFixture}
                          data={fixture}
                          team_id={this.props.team_id}
                          channel={this.props.channel}/>
                  </div>)
                      )
                }
          </div>
        </div>
      );
  }
}

export default FixtureList;
