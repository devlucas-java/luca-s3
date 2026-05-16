
// protoc --go_out=. --go-grpc_out=. proto/video.proto

olha funcionou o endpoint de trancoder e o de 

    - delete /obs ele sempre retorna true e nem um erro, tem que colocar validacao se exixte tal coisa
    
    - get job / ele retorna ok mas ele tem que usar o id do minio e nao dele tem qeu ajsutar
    
    - trancoder / ele faz o trancoder mas ele demora tem qeu colcar para runtime e talvex um stream em tempo real 

    - get hls manifest / funciona mas o link nao funcionou no test, e tem que ajustar a data na idicacao do grpc
