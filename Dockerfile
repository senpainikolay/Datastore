FROM golang:latest

RUN mkdir /app 


ARG config config
ARG port1  
ARG port2

COPY . /app  

COPY . /app  
# Replacing the configurations folder files with needed configurations 
COPY ${config} /app/config
WORKDIR /app 

RUN export GO111MODULE=on  
RUN go mod tidy 
EXPOSE ${port1} ${port2}
CMD go run ./main