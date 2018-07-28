var ActiveUserID = (function() {
  var id = -1;

  var getID = function() {
    return id;
  };

  var setID = function(new_id) {
    id = new_id;
  };

  return {
    getID: getID,
    setID: setID
  }

})();

export default ActiveUserID;
