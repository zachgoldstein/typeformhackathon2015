var bootstrapcss = require('./css/bootstrap.min.css');
var css = require('./css/style.css');

import React from "react";
import ShowCounts from "./showCounts.jsx"
import ShowVoterImages from "./ShowVoterImages.jsx"
import ShowVoteFeed from "./ShowVoteFeed.jsx"

React.render(
  <ShowCounts />,
  document.getElementById('showCounts')
);

React.render(
  <ShowVoterImages />,
  document.getElementById('showVoters')
);

React.render(
  <ShowVoteFeed />,
  document.getElementById('showFeed')
);
