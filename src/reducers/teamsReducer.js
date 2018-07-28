import {ADD_TEAM, FETCH_PIN, FETCH_TEAMS} from "../actions/types";

const initialState = [];

export default function (state = initialState, action) {
    switch (action.type) {
        case FETCH_TEAMS:
            return action.payload;
        case FETCH_PIN:
            const teamIndex = state.indexOf(action.payload.team);
            const newState = state.slice();
            const playerIndex = newState[teamIndex].players.indexOf(action.payload.player);
            newState[teamIndex].players[playerIndex] = {
                ...newState[teamIndex].players[playerIndex],
                latlng: action.payload.latlng
            };
            return newState;
        case ADD_TEAM:
            return [...state, action.payload];
        default:
            return state;
    }
}