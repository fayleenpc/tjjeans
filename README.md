TJ Jeans E-Commerce ( Simple E Commerce )

It has features :

1. Sign Up user verified directly using random token ( let user click on it to be able for login in future )

2. Login unless verified with token earlier from register

3. Rate Limit per Request ( 1 second per 1 request )

4. JWT Token accross API Endpoint for securely

5. Payment Gateway ( Xendit & Midtrans able )

6. Database Migration

7. Validation


To Do :

1. Implement with Autocert https for production

2. Verify Email with real email using SMTP/Gmail-like verification code

3. Build Microservices between types.go

4. Monitor any services with Prometheus/Jaeger UI combined with Grafana

5. Use secure environtment variable d with Hashicorp Vault instead for production 

6. CI/CD, Docker, K8s, Github Workflows/Action, any-ignored file should be it

7. Accounting Details about every products has been sold so that in future there's structure how much profit

8. Frontend need to be separated-like ( for example, don't make it combined together, split instead only index.html -> using templ+htmx )

9. A lack of admin to manage even CRUD so it's self seed application.

10. Use cache liftbridge/nats

11. Use callback for every payment so that 'status' order state 'pending' changed lively 'paid', we need production mode api from xendit/Midtrans for this service



