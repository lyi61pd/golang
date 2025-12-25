package main

import (
	"fmt"
	"log"

	// 导入生成的 proto 代码
	pb "protobuf-demo/proto"

	"google.golang.org/protobuf/proto"
)

func main() {
	fmt.Println("=== Protobuf 学习示例 ===\n")

	// 1. 创建 protobuf 消息对象
	user := &pb.User{
		Id:    1001,
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
		Tags:  []string{"golang", "protobuf", "developer"},
		Address: &pb.Address{
			City:    "北京",
			Street:  "中关村大街1号",
			ZipCode: "100000",
		},
		Gender: pb.Gender_MALE,
	}

	fmt.Println("1. 创建的用户对象:")
	fmt.Printf("   ID: %d\n", user.Id)
	fmt.Printf("   姓名: %s\n", user.Name)
	fmt.Printf("   邮箱: %s\n", user.Email)
	fmt.Printf("   年龄: %d\n", user.Age)
	fmt.Printf("   标签: %v\n", user.Tags)
	fmt.Printf("   地址: %s %s\n", user.Address.City, user.Address.Street)
	fmt.Printf("   性别: %s\n", user.Gender.String())

	// 2. 序列化：将对象转换为二进制数据
	data, err := proto.Marshal(user)
	if err != nil {
		log.Fatalf("序列化失败: %v", err)
	}
	fmt.Printf("\n2. 序列化后的二进制数据长度: %d 字节\n", len(data))
	fmt.Printf("   二进制数据(hex): %x\n", data)

	// 3. 反序列化：将二进制数据还原为对象
	newUser := &pb.User{}
	err = proto.Unmarshal(data, newUser)
	if err != nil {
		log.Fatalf("反序列化失败: %v", err)
	}
	fmt.Println("\n3. 反序列化后的对象:")
	fmt.Printf("   ID: %d\n", newUser.Id)
	fmt.Printf("   姓名: %s\n", newUser.Name)
	fmt.Printf("   邮箱: %s\n", newUser.Email)

	// 4. 创建用户列表
	userList := &pb.UserList{
		Users: []*pb.User{
			user,
			{
				Id:     1002,
				Name:   "李四",
				Email:  "lisi@example.com",
				Age:    30,
				Gender: pb.Gender_FEMALE,
			},
		},
	}

	fmt.Printf("\n4. 用户列表包含 %d 个用户:\n", len(userList.Users))
	for i, u := range userList.Users {
		fmt.Printf("   用户%d: %s (%s)\n", i+1, u.Name, u.Email)
	}

	// 5. 模拟请求响应
	request := &pb.GetUserRequest{
		UserId: 1001,
	}
	fmt.Printf("\n5. 请求用户ID: %d\n", request.UserId)

	response := &pb.GetUserResponse{
		User:    user,
		Message: "查询成功",
	}
	fmt.Printf("   响应消息: %s\n", response.Message)
	fmt.Printf("   返回用户: %s\n", response.User.Name)

	// 6. 比较两个消息是否相等
	isEqual := proto.Equal(user, newUser)
	fmt.Printf("\n6. 原始对象与反序列化对象是否相等: %v\n", isEqual)

	fmt.Println("\n=== 示例结束 ===")
}
