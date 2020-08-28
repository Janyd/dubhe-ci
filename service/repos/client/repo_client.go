package main

import "path/filepath"

const (
// Address 连接地址
//Address string = "127.0.0.1:9090"
)

func main() {
	//// 连接服务器
	//conn, err := grpc.Dial(Address, grpc.WithInsecure())
	//if err != nil {
	//	logrus.Fatalf("net.Connect err: %v", err)
	//}
	//defer conn.Close()
	//
	//// 建立gRPC连接
	//grpcClient := pb.NewRepositoryClient(conn)
	//req := pb.Repo{
	//	Name:         "zefun",
	//	Description:  "sdfsgdfhdhd",
	//	Url:          "git@code.aliyun.com:maywant6.0/zefun.git",
	//	CredentialId: 1,
	//	Timeout:      1,
	//}
	//
	//res, err := grpcClient.Create(context.Background(), &req)
	//if err != nil {
	//	logrus.Fatalf("Call Route err: %v", err)
	//}
	//// 打印返回值
	//logrus.Println(res)
	_path := filepath.Join("/a/b/", "preject", "sdfs/demo.txt")
	println(_path)
}
