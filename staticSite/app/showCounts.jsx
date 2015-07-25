import React from "react";

export default React.createClass({
  mixins: [ReactFireMixin],

  componentWillMount: function() {
      var ref = new Firebase("https://typeformhackathon.firebaseio.com/answers");
    this.bindAsArray(ref, "answers");
  },

  render: function() {

    var answers = [];

    var choiceCounts = {};

    for (var answerKey in this.state.answers) {
      var answer = this.state.answers[answerKey];

      var choice = answer["answers"][0]["data"]["value"]["label"];
      choice = choice.replace('?', "");
      if (!choiceCounts[choice]) {
        choiceCounts[choice] = 1;
      } else {
        choiceCounts[choice] += 1;
      }
    }

    for (var choice in choiceCounts) {
      var count = choiceCounts[choice];
      if (answers.length == 0) {
        var style = {
          "text-align":"left"
        }
        answers.push(
          <div className="col-md-6" style={style}>
            <h3 key={choice}>Would {choice} - {count}</h3>
          </div>
        );
      } else {
        var style = {
          "text-align":"right"
        }
        answers.push(
          <div className="col-md-6" style={style}>
            <h3 key={choice}>{count} - Would {choice}</h3>
          </div>
        );
      }
    }

    var overallStyle = {
      "font-family":"04b_19regular"
    }

    return (
      <div style={overallStyle}>
        <h1>TYPEFORM HACKATHON 2015</h1>
        <div className="row" id="count-row">
          {answers}
        </div>
      </div>
    );
  },
});
