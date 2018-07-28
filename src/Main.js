import React from 'react';
import Chat from './Chat/Chat';
import Profile from './Profile';
import Leaderboards from './Leaderboards/Leaderboards';
import Matchmaking from './Matchmaking/Matchmaking';
import Teams from './Teams/Teams';
import CreateTeam from './Teams/CreateTeam';
import {Switch, Route} from 'react-router-dom'

const Main = () => (
    <Switch>
        <Route exact path='/' component={Profile}/>
        <Route path='/chat' component={Chat}/>
        <Route path='/profile' component={Profile}/>
        <Route path='/leaderboards' component={Leaderboards}/>
        <Route path='/matchmaking' component={Matchmaking}/>
        <Route path='/createTeam' component={CreateTeam}/>
        <Route path='/Teams' component={Teams}/>
    </Switch>
)

export default Main;
