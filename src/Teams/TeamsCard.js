import React, {Component} from 'react';
import {Tab, Tabs} from "react-bootstrap";
import {connect} from "react-redux";
import {fetchTeams} from "../actions/teamsActions";
import "./Stylesheets/TeamCardsGrid.css";
import SimpleMap from "./MiniMap";
import TeamTab from "./TeamTab";

class TeamsCard extends Component {

    componentWillMount() {
        this.props.fetchTeams();
    }

    render() {
        return (
          <div id="teamscard">
            <h2> Team Info</h2>
            <Tabs defaultActiveKey={1}>
                {this.props.teams.map((t, index) => {
                    return (
                    <Tab eventKey={index + 1} title={t.name}>
                        <TeamTab team={t}/>
                    </Tab>)}
                )}
            </Tabs>
            </div>
        );
    }
}

const mapStateToProps = (state, ownProps) => ({
    ...ownProps,
    teams: state.teams
});

export default connect(mapStateToProps, {fetchTeams})(TeamsCard)
