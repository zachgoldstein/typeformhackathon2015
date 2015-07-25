import React from "react";

export default React.createClass({
  mixins: [ReactFireMixin],

  componentWillMount: function() {
    var ref = new Firebase("https://typeformhackathon.firebaseio.com/votedUsers");
    this.bindAsArray(ref, "votedUsers");
  },

  render: function() {

    var voters = [];

    var count = 0;
    for (var voterKey in this.state.votedUsers) {
      var voter = this.state.votedUsers[voterKey];
      console.log(voter.profile["image_192"]);
      voters.push(
        <img className="voterImg" ref={"image-"+count} src={voter.profile["image_192"]}></img>
      )
      count++;
    }

    return (
      <div>
          <h1 id="voter-header">voters</h1>
        {voters}
      </div>
    );
  },
});
