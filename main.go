package main

import (
        "fmt"
        "net/http"
)

func main() {
        fmt.Println("Starting hello-world server...")
        http.HandleFunc("/", helloServer)
        if err := http.ListenAndServe(":8080", nil); err != nil {
                panic(err)
        }
}

func helloServer(w http.ResponseWriter, r *http.Request) {
        message := `
<!DOCTYPE html>
<html>
<head>
    <title>Welcome to Okteto</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            color: #333;
        }
        h1 {
            color: #0066cc;
        }
        .container {
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 20px;
            background-color: #f9f9f9;
        }
        .highlight {
            font-weight: bold;
            color: #0066cc;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to Okteto!</h1>
        <p>This application is running in a <span class="highlight">Kubernetes cluster</span> powered by Okteto.</p>
        <p>Okteto is a developer platform that makes it easy to develop, deploy, and test Cloud Native applications on Kubernetes.</p>
        <p>With Okteto, you can:</p>
        <ul>
            <li>Develop directly in Kubernetes</li>
            <li>Automatically deploy your applications</li>
            <li>Preview your changes in real-time</li>
            <li>Collaborate with your team seamlessly</li>
            <li>Test your applications in a production-like environment</li>
        </ul>
        <p>Happy cloud native development!</p>
    </div>
</body>
</html>
`
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, message)
}
