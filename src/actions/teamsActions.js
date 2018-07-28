import {serverDomain} from "../ipconfig"
import {FETCH_TEAMS} from "./types";
import ActiveUserID from "../Profile/ActiveUserID";

export const fetchTeams = () => dispatch => {
    //get team
    fetch(serverDomain + "/getTeams/" + ActiveUserID.getID())
        .then(res => res.json())
        .then(teams => {
          // console.log(teams)
            dispatch({
                type: FETCH_TEAMS,
                payload: teams
            });
            return teams;
        });
};
