FROM golang:latest

RUN mkdir /app 


ARG port  


COPY . /app  
WORKDIR /app  
 
RUN export GO111MODULE=on  
RUN go mod tidy  

EXPOSE ${port}
CMD go run main.go