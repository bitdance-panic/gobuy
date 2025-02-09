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

struct GetUserReq {
    1: string user_id
}


struct GetUserResp {
    1: bool success
    2: string user_id
    3: string email
    4: string username
}

struct UpdateUserReq {
    1: string user_id
    2: optional string email
    3: optional string username
}

struct UpdateUserResp {
    1: bool success
}

//软删除
struct DeleteUserReq {
    1: string user_id
}

struct DeleteUserResp {
    1: bool success
}

// 定义服务接口
service UserService {
    RegisterResp Register(1: RegisterReq req)
    LoginResp Login(1: LoginReq req)
    GetUserResp GetUser(1: GetUserReq req)
    UpdateUserResp UpdateUser(1: UpdateUserReq req)
    DeleteUserResp DeleteUser(1: DeleteUserReq req)
}