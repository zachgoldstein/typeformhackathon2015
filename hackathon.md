Remaining TODOS:

clean up data
run through talk a few times
commit it all to git

rip out library for golang typeform API

Remaining TODOS:

1. fix images DONE
2. DM eveybody in channel with form link after creation DONE
3. Custom forms for everybody, querystring with userID DONE
4. On form completion, store user image link with answer DONE
5. Show user images after voting DONE
6. Add reaction to form message DONE
7. Send thank you DM after voting and link to results page DONE
........

add filter for images DITCH jsut do B&W DONE
feed for who just voted DONE

add sexy typewriter effect DONE
add goofy running man everytime somebody votes DITCH

fix up golang API lib

add sexy parallax bg DITCH



stuff:

@typebot: Would you rather | Fight a horse sized duck?, 2550 | Fight a thousand duck sized horses?, 2549
https://forms.typeform.io/to/f2l0rfXiX4GBIw




post msg to typeform bot DONE
bot hits golang server to post to typeform API DONE
filling form out adds a reaction in slack
redirects to results link on finish - NOPE
DM is sent to person with results link
webhook hits a golang server DONE
golang server sticks data in firebase  DONE
static site serves it result counts DONE
site takes pic and sends to firebase DITCH
displays everybody's pics on page DITCH

SHOW user's icons after they vote


outstanding questions:

POST req to typeform API, create form with multiple options
DONE
how to redirect after form

how to hit webhook after req
DONE
how to DM person with results link after request success?
https://github.com/nlopes/slack

golang store in firebase?
https://github.com/melvinmt/firebase

take picture and deal with it in react
http://www.html5rocks.com/en/tutorials/getusermedia/intro/#toc-screenshot
???? too hard


display pic with dot matrix filter?
http://evanw.github.io/glfx.js/demo/

twilio idea:
text a number with a question. Text other numbers to ask people...

text question to phone number
text url to friends
text for who won


why is fieldId "field_id":1.151587e+06 ????

webpack-dev-server --hot --inline

http://docs.typeform.io/page/sandbox
