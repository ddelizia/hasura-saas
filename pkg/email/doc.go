/*
This microservice will take care of sendinng email

* can be sent to multiple addresses
* has a subject
* has template
* has data that will be used from the template

The email srcice can have multiple implemetations such as Sendgrid, Mailchimp, AWS SES.
Some of the services above offers Template buiding but templating system can be actually implementes separetelly.
*/
package email