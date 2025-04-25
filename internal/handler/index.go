package handler

import (
	"net/http"
	"text/template"
)

type PageData struct {
	Token string
}

// TODO: this is temporary
// ServePage renders the page with the API token
func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Example API token, in a real case, retrieve this dynamically
	query := r.URL.Query()
	token := query.Get("token")

	// Prepare the data for the template
	data := PageData{Token: token}

	// Parse the template
	tmpl, err := template.New("page").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Gazzette auth token</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				background-color: #f4f4f9;
				margin: 0;
			}
			.container {
				padding: 20px;
				background-color: #fff;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
				border-radius: 8px;
				text-align: center;
				max-width: 400px;
				width: 100%;
			}
			button {
				padding: 10px 15px;
				font-size: 1rem;
				background-color: #4CAF50;
				color: white;
				border: none;
				border-radius: 5px;
				cursor: pointer;
				width: 100%;
			}
			button:hover {
				background-color: #45a049;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Gazzette auth token</h2>
			<button onclick="copyToken()">Copy token</button>
		</div>

		<script>
			function copyToken() {
				// Create a temporary input element to copy the token
				const tempInput = document.createElement('input');
				tempInput.value = "{{.Token}}";
				document.body.appendChild(tempInput);
				tempInput.select();
				document.execCommand('copy');
				document.body.removeChild(tempInput);
				alert('Copied to clipboard!');
			}
		</script>
	</body>
	</html>
	`)
	// Check for parsing errors
	if err != nil {
		http.Error(w, "Could not render page", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
