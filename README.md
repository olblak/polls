# Polls

Polls is a simple application that collect people interested to participate to a poll by listening on a specific endpoint

## API

**/participants(GET)[poll]**  

Return every participants for a specific poll  

`curl "http://localhost:8082/participants?poll=12312321"`

**/participants(POST)[poll]** 

Generate a list of every participants for a specific poll, by default consider every jenkins account created before the first of september 2019  

`curl -X POST 'http://localhost:8080/participants?poll=12312321&'`

**/participate(GET)**  

Request participation for a specific poll by specifying the right token associated to the right email address  

`curl "http://localhost:8082/participate?poll='12312321'&email='me@olblak.com'&token=d3d3311f641372c0f777cb79ed7fea01"`
