import {FETCH_INVITATIONS, DECLINE_INVITATION, ACCEPT_INVITATION, ADD_TEAM} from "./types";
import {serverDomain} from "../ipconfig"
import UserProfile from "../Profile/UserProfile";

export const fetchInvitations = () => dispatch => {
    fetch(serverDomain + "/getInvitations/" + UserProfile.getName())
        .then(res => res.json())
        .then(invitations => dispatch({
            type: FETCH_INVITATIONS,
            payload: invitations
        }));
};

export const acceptInvitation = invitation => dispatch => {
    dispatch({
        type: ACCEPT_INVITATION,
        payload: invitation
    }).then(
        fetch(serverDomain + "/addUserToTeam/" + UserProfile.getName() + "/" + invitation.name,
            {method: "POST"})
            .then(res => res.json())
            .then(updatedteam => dispatch({
                type: ADD_TEAM,
                payload: updatedteam
            }))
    );
};

export const declineInvitation = invitation => dispatch => {
    fetch(serverDomain + "/deleteInvitation/" + UserProfile.getName() + "/" + invitation.name,
        {method: "DELETE"})
        .then(
            dispatch({
                type: DECLINE_INVITATION,
                payload: invitation
            }));
};