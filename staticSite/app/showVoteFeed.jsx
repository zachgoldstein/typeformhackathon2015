import React from "react";

var ref = new Firebase("https://typeformhackathon.firebaseio.com/votedUsers");
// Attach an asynchronous callback to read the data at our posts reference

export default React.createClass({

  getInitialState: function(){
    return {
      userVotes:[],
      finishedInitialLoad: false,
      currentText: "",
      currentScrambleText: "",
      currentTextIndex: 0
    }
  },

  setFinishedLoad: function() {
    var state = this.state;
    state.finishedInitialLoad = true;
    this.setState(state);
  },

  addPerson: function(name) {
    var state = this.state;
    state.currentText = name + " has voted";
    state.userVotes.push(name);
    //animate last
    console.log("added person! ", name);
    console.log("current feed! ", state.userVotes);
    this.setState(state);
    this.startScrambler();

  },

  componentDidMount: function() {
    setTimeout(this.setFinishedLoad, 2000);

    var self = this;
    ref.on("child_added", function(snapshot, prevChildKey) {
      var person = snapshot.val()
      if(self.state.finishedInitialLoad) {
        self.addPerson(person.name);
      }
    }, function (errorObject) {
      console.log("The read failed: " + errorObject.code);
    });

    this.props.scrambleLength = 500;
    this.props.scrambleFrequency = 10;

    var currentState = this.state;

  },

  startScrambler: function () {
    var currentState = this.state;
    currentState.currentRun = 1;

    currentState.textScrambler = setInterval(this.scrambleProgress, this.props.scrambleFrequency);

    currentState.scramblerTimeout = setTimeout(this.stopScrambler, this.props.scrambleLength);

  },

  scrambleProgress: function () {
      var currentState = this.state;
      currentState.currentText = this.state.currentText;

      var numTotalRuns = this.props.scrambleLength/this.props.scrambleFrequency;

      currentState.progress = currentState.currentRun/ numTotalRuns;

      var textToScramble = this.state.currentText;
      var scrambledText = "";

      for (var i = 0, len = textToScramble.length; i < len; i++) {
        var char = textToScramble[i];
        var textPosition = i/textToScramble.length;

        if (this.state.progress < textPosition) {
          var letter = this.findScrambledCharacter();
          scrambledText += letter;
        } else {
          scrambledText += textToScramble[i];
        }
      }

      currentState.currentRun += 1;
      currentState.scrambledText = scrambledText;

      this.setState(currentState);

    },

    findScrambledCharacter: function () {
      var charArray = ['A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z'];
      return charArray[Math.floor(Math.random() * charArray.length)];
    },

    stopScrambler: function () {
      clearTimeout(this.state.textScrambler);
      clearInterval(this.state.scramblerTimeout);

      var currentState = this.state;
      currentState.scrambledText = this.state.currentText;
      this.setState(currentState);
    },

  render: function(){
    var self = this;

    var feed = []
    this.state.userVotes.forEach(function(vote, index){
      if (index < self.state.userVotes.length - 1) {
        feed.push(
          <div className="row-fluid">
            <h3>{vote} has voted</h3>
          </div>
        );
      }
    });
    feed.push(
      <div className="row-fluid">
        <h3>{self.state.scrambledText}</h3>
      </div>
    );

    return (
      <div>
        <h2>Latest Data</h2>
        {feed}
      </div>
    );
  }
});
