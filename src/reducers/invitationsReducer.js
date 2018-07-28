import {ACCEPT_INVITATION, DECLINE_INVITATION, FETCH_INVITATIONS} from "../actions/types";

const initialState = [];

export default function (state = initialState, action) {
    switch (action.type) {
        case ACCEPT_INVITATION:
            return state.filter(teams => teams.name !== action.payload.name);
        case DECLINE_INVITATION:
            return state.filter(team => team.name !== action.payload.name);
        case FETCH_INVITATIONS:
            return action.payload;
        default:
            return state;
    }
}