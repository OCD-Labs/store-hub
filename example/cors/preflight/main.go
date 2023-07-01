package main

import (
	"flag"
	"log"
	"net/http"
)

const html = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
</head>
<body>
	<h1>Preflight CORS</h1>
	<div id="output"><div>
	<script>
		document.addEventListener('DOMContentLoaded', function() {
			fetch("https://storehub-zjp3.onrender.com/api/v1/users", {
				method: "POST",
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					account_id: 'v2.storehub.testnet',
					first_name: 'string',
					last_name: 'string',
					password: 'stringst',
					email: 'user-2@gmail.com',
					profile_image_url: 'string'
				})
			}).then(
				function(response) {
					response.text().then(function(text) {
						document.getElementById("output").innerHTML = text;
					});
				},
				function(err) {
					document.getElementById("output").innerHTML = err;
				}
			);
		});
	</script>
</body>
`

func main() {
	addr := flag.String("addr", ":3000", "Server address")
	flag.Parse()

	log.Printf("starting server on %s", *addr)

	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))

	log.Fatal(err)
}