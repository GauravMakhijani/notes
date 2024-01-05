# Replace the following URL with your actual endpoint URL
API_URL="http://localhost:8080/api/auth/ping"

ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImViMGE1ZjZjLTI0YmQtNGRiOC1iZTI2LTRmNTk5NTliY2Y5MyIsInVzZXJuYW1lIjoidGVzdDQifQ.X0wCcxOgDXZX5xEdjenqgOQulg00Qfa0J2Y8Jc03rU4"
# Number of requests to send
REQUESTS=100

# Loop to send requests
for ((i=1; i<=$REQUESTS; i++)); do
  echo "Sending request $i with token"
  curl -X GET $API_URL \
    -H "Authorization: Bearer $ACCESS_TOKEN"
  echo ""
done