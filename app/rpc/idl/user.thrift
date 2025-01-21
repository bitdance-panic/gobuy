namespace go user

// 定义请求和响应的结构体
struct RegisterReq {
    1: string email
    2: string password
    3: string confirm_password
}

struct RegisterResp {
    1: i32 user_id
}

struct LoginReq {
    1: string email
    2: string password
}

struct LoginResp {
    1: bool success
    2: i32 user_id
}

// 定义服务接口
service UserService {
    RegisterResp Register(1: RegisterReq req)
    LoginResp Login(1: LoginReq req)
}