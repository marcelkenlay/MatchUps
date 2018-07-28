import {FETCH_PIN} from "./types";


export const fetchPins = (player,team) => dispatch => {
    fetch("https://maps.googleapis.com/maps/api/geocode/json?&address=" +
        encodeURI(player.location) + "&key=AIzaSyC8Ux3avYGKFFPFl3EEmHVOqqRF4sfBJdk")
        .then(res => res.json())
        .then(geoinfo => {
            dispatch({
                type: FETCH_PIN,
                payload: {
                    team: team,
                    player: player,
                    latlng: geoinfo.results[0].geometry.location
                }
            });
        });
};