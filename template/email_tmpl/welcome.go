package email_tmpl

import (
	"fmt"
)

func WelcomeTmpl(firstName, url string) string {
	s := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Email Verification</title>
	<style>
			body {
					font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
					margin: 0;
					padding: 0;
					background-color: #f7f7f7;
					color: #333;
			}
			.container {
					max-width: 560px;
					margin: 40px auto;
					padding: 20px;
					background-color: #ffffff;
					border: 1px solid #dedede;
					border-radius: 5px;
					box-shadow: 0 5px 15px rgba(0,0,0,0.1);
			}
			.button {
					display: block;
					width: max-content;
					margin: 20px auto 0;
					padding: 10px 25px;
					background-color: #4CAF50;
					color: #ffffff;
					text-decoration: none;
					border-radius: 4px;
					font-weight: bold;
			}
			.content {
					text-align: center;
					line-height: 1.5;
			}
			.footer {
					margin-top: 20px;
					font-size: 12px;
					text-align: center;
					color: #999;
			}
	</style>
	</head>
	<body>
	<div class="container">
			<div class="content">
					<h2>Hello %s,</h2>
					<p>Thank you for registering with us!</p>
					<a href="%s" class="button">Verify Your Email</a>
			</div>
			<div class="footer">
					If you did not request this email, please ignore it.
			</div>
			<span class="random"></span>
	</div>
	</body>
	</html>	
 `, firstName, url)
 return s
}

