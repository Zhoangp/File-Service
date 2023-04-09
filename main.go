package main
import (
	"fmt"
	"github.com/Zhoangp/File-Service/config"
	"github.com/Zhoangp/File-Service/internal/delivery/https"
	"github.com/Zhoangp/File-Service/pb"
	"github.com/Zhoangp/File-Service/pkg/upload"
	"google.golang.org/grpc"
	"log"
	"net"
)
func main(){
	cf, err := config.LoadConfig("config/config-local.yml")
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", ":" + cf.App.Port)
	fmt.Println("Auth Svc on", cf.App.Port)
	s3 := upload.NewS3Provider(cf)
	hdl := https.NewUploadHandler(s3)
	grpcServer := grpc.NewServer()
	pb.RegisterFileServiceServer(grpcServer, hdl)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}


}
