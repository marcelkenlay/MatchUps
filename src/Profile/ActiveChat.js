var ActiveChat = (function() {
  var active_chat = "none";

  var getChat = function() {
    return active_chat;
  };

  var setChat = function(chat) {
    active_chat = chat;
  };

  return {
    getChat: getChat,
    setChat: setChat
  }

})();

export default ActiveChat;
