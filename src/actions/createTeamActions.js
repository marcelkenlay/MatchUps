import {FETCH_USERNAME_QUERY_RESULTS, UPDATE_USERNAME_QUERY} from "./types";
import ActiveUserID from "../Profile/ActiveUserID";
import * as axios from "moment";

export const getQueryResults = query => dispatch => {
    var req = "/getUsernameMatches?pattern=" + query;
    this.serverRequest = axios.get(req)
        .then(result => {
            if (result.data != null) {
                result.data.filter((res) => res.UserID !== ActiveUserID.getID())
                dispatch({
                    type: FETCH_USERNAME_QUERY_RESULTS,
                    payload: result.data
                });
            }
        });
};