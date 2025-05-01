API REALIZATION:
MAIN ADDR: "localhost/TaskList:7878"
If we want to register: "localhost/TaskList/Register:7878"
JSON: 
{
    "username" : "your_login",
    "email" : "your_email",
    "password" : "your_password",
    "full_name" : "full_name"
}
GOT :
json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":   response.Id,
			"username":  response.Username,
			"email":     response.Email,
			"full_name": response.FullName,
		})

If we want to login and get JWT key: "localhost/TaskList/Login:7878"
JSON:
{
    "username" : "your_login",
    "password" : "your password"
}

GOT:
{
    "jwt" : "your_personal_code_jwt"
}

After authorization you can use service by jwt code, you need just add "jwt" field in json request

