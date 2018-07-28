import React, {Component} from 'react';

const QueryResultCard = props => {
  return <div class={"wrapper"}>
    <h3>Username: {props.username}</h3>
    Name: {props.name}
  </div>
}

export default QueryResultCard;
