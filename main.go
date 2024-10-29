package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fatimaalcaraz/TwitterGo/awsgo"
	"github.com/fatimaalcaraz/TwitterGo/bd"
	"github.com/fatimaalcaraz/TwitterGo/handlers"
	"github.com/fatimaalcaraz/TwitterGo/models"
	"github.com/fatimaalcaraz/TwitterGo/secretmanager"
)

func main() {
	lambda.Start(EjecutoLambda)

}

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var resp *events.APIGatewayProxyResponse

	awsgo.InicializoAws()

	if !ValidoParametros() {
		resp = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno, deben incluir 'SecretName', 'BucketName', 'UrlPrefix'",
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return resp, nil
	}
	SecretModel, err := secretmanager.GetSecret(os.Getenv("secretName"))
	if err != nil {
		resp = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la lectura de Secret" + err.Error(),
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return resp, nil
	}
	path := strings.Replace(request.PathParameters["TwitterGo"], os.Getenv("UrlPrefix"), " ", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtsign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))
	// chequeo conexion a la BD o conecto la BD
	err = bd.ConectarBD(awsgo.Ctx)
	if err != nil {
		resp = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error conectando la BD" + err.Error(),
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return resp, nil
	}
	respAPI := handlers.Manejadores(awsgo.Ctx, request)
	if respAPI.CustomResp == nil {
		resp = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       respAPI.Message,
			Headers: map[string]string{
				"Content-type": "application/json",
			},
		}
		return resp, nil

	} else {
		return respAPI.CustomResp, nil
	}
}

func ValidoParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}
	return traeParametro
}
