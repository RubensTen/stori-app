# Usa la imagen base de Go para construir tu código
FROM golang:1.21 as build
# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app
# Copia tu código Go al directorio de trabajo
COPY function/ .

# Go modules will be installed into a directory inside the image.
#RUN go mod download

# Compilar el binario de la función Lambda
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o /build/bootstrap main.go

# Usa la imagen Lambda de AWS como base
FROM public.ecr.aws/lambda/provided:al2023

# Copia el binario compilado desde el paso anterior
COPY --from=build /build/bootstrap ./bootstrap
# Especificar el comando de entrada para la función Lambda
ENTRYPOINT [ "./bootstrap" ]

