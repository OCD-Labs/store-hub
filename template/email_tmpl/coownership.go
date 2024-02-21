package email_tmpl

import "fmt"

func CoOwnershipAccessTmpl(inviteeName, inviterName, storeName, accessLevelDesc, msg, url string) string {
	cnt := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Invitation To Manage A Store</title>
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
						margin: 20px auto;
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
						<p>%s has invited you to join %s on StoreHub with %s privileges. %s</p>
						<a href="%s" class="button">Click here</a> to accept the invitation and start managing %[3]s.
				</div>
				<div class="footer">
					If you did not expect this invitation or believe it's an error, please ignore this email or contact our support.
				</div>
		</div>
		</body>
		</html>
	`, inviteeName, inviterName, storeName, accessLevelDesc, msg, url)
	return cnt
}