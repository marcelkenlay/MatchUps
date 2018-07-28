var UserProfile = (function() {
  var full_name = "none";

  var getName = function() {
    return full_name;
  };

  var setName = function(name) {
    full_name = name;
  };

  return {
    getName: getName,
    setName: setName
  }

})();

export default UserProfile;
