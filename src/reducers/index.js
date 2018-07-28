import {combineReducers} from "redux";
import teamsReducer from "./teamsReducer";
import invitationsReducer from "./invitationsReducer";

export default combineReducers({
    teams: teamsReducer,
    invitations: invitationsReducer,
})