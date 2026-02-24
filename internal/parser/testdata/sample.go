package main

type UserService struct {
	Field string `diagram:"type=service,name=UserService,connectsTo=UserDatabase;MessageQueue"`
}

type OrderService struct {
	Field string `diagram:"type=service,name=OrderService,connectsTo=OrderDatabase;PaymentGateway"`
}

type PaymentGateway struct {
	Field string `diagram:"type=external,name=PaymentGateway,description=Stripe payments"`
}

type UserDatabase struct {
	Field string `diagram:"type=database,name=UserDatabase"`
}

type OrderDatabase struct {
	Field string `diagram:"type=database,name=OrderDatabase"`
}

type MessageQueue struct {
	Field string `diagram:"type=queue,name=MessageQueue"`
}

type Cache struct {
	Field string `diagram:"type=cache,name=Cache,description=Redis cache"`
}

type APIGateway struct {
	Field string `diagram:"type=gateway,name=APIGateway,connectsTo=AuthService;UserService;OrderService"`
}

type AuthService struct {
	Field string `diagram:"type=service,name=AuthService,description=Authentication service"`
}

func main() {}

type AuthService struct{}

func main() {}
